package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"gorm.io/gorm"

	"api/models"
)

// Notification type strings emitted by the underwriting workflow. Mirror the
// renderer's `notifications.ts` channel filter when adding new ones.
const (
	NotificationTypeUWCaseCreated      = "uw_case_created"
	NotificationTypeUWCaseDecided      = "uw_case_decided"
	NotificationTypeUWEvidenceRequired = "uw_evidence_required"
)

// uwCaseAttachmentDir is the on-disk root for case attachments.
const uwCaseAttachmentDir = "tmp/uploads/underwriting_cases"

// allowedUWCaseTransitions lists the valid `from -> to` moves. A case can
// always be assigned an underwriter without a status change (handled
// separately by AssignUnderwritingCase).
var allowedUWCaseTransitions = map[models.UWCaseStatus]map[models.UWCaseStatus]bool{
	models.UWCaseStatusPendingEvidence: {
		models.UWCaseStatusInReview:     true,
		models.UWCaseStatusDeclined:     true,
		models.UWCaseStatusPostponed:    true,
		models.UWCaseStatusAutoAccepted: true,
	},
	models.UWCaseStatusInReview: {
		models.UWCaseStatusDecided:         true,
		models.UWCaseStatusDeclined:        true,
		models.UWCaseStatusPostponed:       true,
		models.UWCaseStatusPendingEvidence: true,
		models.UWCaseStatusAutoAccepted:    true,
	},
	models.UWCaseStatusPostponed: {
		models.UWCaseStatusPendingEvidence: true,
		models.UWCaseStatusInReview:        true,
		models.UWCaseStatusAutoAccepted:    true,
	},
	// Decided, Declined and AutoAccepted are terminal — re-opening is not
	// supported here. A new case is created by re-running CreateCasesForQuote
	// after re-rate.
}

// UnderwritingCaseFilter narrows ListUnderwritingCases results.
type UnderwritingCaseFilter struct {
	QuoteID                  *int
	Status                   models.UWCaseStatus
	Tier                     *int
	AssignedUnderwriterEmail string
}

// CreateCasesForQuote creates one underwriting case for every member on the
// quote with UnderwritingTier >= 1 (short-form or full-UW). It is idempotent:
// if a case already exists for (QuoteID + member identity) the snapshot
// fields are refreshed but decisions / events / attachments are preserved.
//
// `actor` is recorded on the audit event as the system or user that triggered
// the create. Callers from a controller pass the active user's email; the
// quote-calc-completion hook passes the calc trigger's email.
func CreateCasesForQuote(quoteID int, actor string) ([]models.UnderwritingCase, error) {
	var ratings []models.MemberRatingResult
	if err := DB.Where("quote_id = ?", quoteID).Find(&ratings).Error; err != nil {
		return nil, fmt.Errorf("load ratings: %w", err)
	}

	idLookup := buildMemberIDLookup(quoteID)

	out := make([]models.UnderwritingCase, 0, len(ratings))
	for _, r := range ratings {
		memberID := idLookup[memberLookupKey(r.MemberName, r.Category)]
		c, created, err := upsertCaseForMember(quoteID, memberID, r, actor)
		if err != nil {
			return nil, err
		}
		if c == nil {
			continue
		}
		out = append(out, *c)
		if created {
			emitCaseCreatedNotification(*c, actor)
		}
	}
	return out, nil
}

// upsertCaseForMember either refreshes the snapshot on an existing case or
// creates a fresh case. Returns (case, wasNewlyCreated, error).
func upsertCaseForMember(quoteID int, memberID string, r models.MemberRatingResult, actor string) (*models.UnderwritingCase, bool, error) {
	var existing models.UnderwritingCase
	q := DB.Where("quote_id = ? AND member_name = ? AND category = ?", quoteID, r.MemberName, r.Category)
	if memberID != "" {
		q = q.Where("member_id_number = ? OR member_id_number = ''", memberID)
	}
	err := q.First(&existing).Error

	if err == nil {
		previousStatus := existing.Status
		existing.Tier = r.UnderwritingTier
		existing.FCLExcessRatio = r.FCLExcessRatio
		existing.GlaSumAssured = r.GlaSumAssured
		existing.PtdSumAssured = r.PtdSumAssured
		existing.CiSumAssured = r.CiSumAssured
		existing.SpouseGlaSumAssured = r.SpouseGlaSumAssured
		// FreeCoverLimit is stored at the quote level; the rating row doesn't
		// carry it directly, so leave existing value untouched (it was set on
		// initial create from the quote and rarely changes mid-cycle).
		if memberID != "" && existing.MemberIdNumber == "" {
			existing.MemberIdNumber = memberID
		}
		// Auto-close when the member's new tier drops below short-form — the
		// case no longer needs underwriting. Only safe to do from non-terminal
		// statuses where no human has committed a decision; if a human paused
		// (postponed) or is actively reviewing, leave the workflow intact so
		// their context is preserved.
		autoClosed := false
		if r.UnderwritingTier < UnderwritingTierShortForm &&
			(existing.Status == models.UWCaseStatusPendingEvidence ||
				existing.Status == models.UWCaseStatusInReview) {
			now := time.Now()
			existing.Status = models.UWCaseStatusAutoAccepted
			existing.DecidedAt = &now
			existing.DecidedBy = "system"
			autoClosed = true
		}
		if err := DB.Save(&existing).Error; err != nil {
			return nil, false, fmt.Errorf("refresh case: %w", err)
		}
		recordCaseEvent(existing.ID, "tier_refreshed", actor, map[string]any{
			"tier":             r.UnderwritingTier,
			"fcl_excess_ratio": r.FCLExcessRatio,
		})
		if autoClosed {
			recordCaseEvent(existing.ID, "status_changed", "system", map[string]any{
				"from": previousStatus,
				"to":   models.UWCaseStatusAutoAccepted,
				"note": "member tier dropped below short-form on re-calc; no underwriting required",
			})
		}
		return &existing, false, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, false, fmt.Errorf("lookup case: %w", err)
	}

	if r.UnderwritingTier < UnderwritingTierShortForm {
		return nil, false, nil
	}

	fcl := loadQuoteFreeCoverLimit(quoteID)
	ruleSetID, ruleSetVersion := loadActiveRuleSetSnapshot()
	c := models.UnderwritingCase{
		QuoteID:             quoteID,
		MemberIdNumber:      memberID,
		MemberName:          r.MemberName,
		Category:            r.Category,
		Tier:                r.UnderwritingTier,
		FCLExcessRatio:      r.FCLExcessRatio,
		GlaSumAssured:       r.GlaSumAssured,
		PtdSumAssured:       r.PtdSumAssured,
		CiSumAssured:        r.CiSumAssured,
		SpouseGlaSumAssured: r.SpouseGlaSumAssured,
		FreeCoverLimit:      fcl,
		Status:              models.UWCaseStatusPendingEvidence,
		RuleSetID:           ruleSetID,
		RuleSetVersion:      ruleSetVersion,
		CreatedBy:           actor,
	}
	if err := DB.Create(&c).Error; err != nil {
		return nil, false, fmt.Errorf("create case: %w", err)
	}
	recordCaseEvent(c.ID, "case_created", actor, map[string]any{
		"tier":             r.UnderwritingTier,
		"fcl_excess_ratio": r.FCLExcessRatio,
	})
	return &c, true, nil
}

// loadActiveRuleSetSnapshot returns (id, version) of the currently-active
// UWRuleSet, or (0, 0) when no active set exists. Cases created when no set
// is active will have RuleSetID/Version=0 and the engine treats them as
// "no rules evaluated yet" — replaying later picks up whatever set was
// snapshotted (still 0, so the case shows "no rules applied").
func loadActiveRuleSetSnapshot() (int, int) {
	var rs models.UWRuleSet
	if err := DB.Where("active = ?", true).Order("version DESC").First(&rs).Error; err != nil {
		return 0, 0
	}
	return rs.ID, rs.Version
}

func loadQuoteFreeCoverLimit(quoteID int) float64 {
	var q models.GroupPricingQuote
	if err := DB.Select("free_cover_limit").Where("id = ?", quoteID).First(&q).Error; err != nil {
		return 0
	}
	return q.FreeCoverLimit
}

// buildMemberIDLookup maps (MemberName + Category) -> MemberIdNumber from the
// uploaded census so cases get a stable identifier when available.
func buildMemberIDLookup(quoteID int) map[string]string {
	var members []models.GPricingMemberData
	if err := DB.Select("member_name, scheme_category, member_id_number").
		Where("quote_id = ?", quoteID).Find(&members).Error; err != nil {
		return map[string]string{}
	}
	out := make(map[string]string, len(members))
	for _, m := range members {
		out[memberLookupKey(m.MemberName, m.SchemeCategory)] = m.MemberIdNumber
	}
	return out
}

func memberLookupKey(name, category string) string {
	return name + "|" + category
}

// UnderwritingQuoteSummary is one row of the quote-grouped underwriting
// queue. Aggregates every case under a single quote into a workstream
// view so underwriters can see "quote X has 7 pending, 2 in review" at
// a glance, without scrolling a flat per-member list.
type UnderwritingQuoteSummary struct {
	QuoteID              int        `json:"quote_id"`
	QuoteName            string     `json:"quote_name"`
	SchemeName           string     `json:"scheme_name"`
	BrokerName           string     `json:"broker_name"`
	QuoteStatus          string     `json:"quote_status"`
	TotalCases           int        `json:"total_cases"`
	PendingEvidenceCount int        `json:"pending_evidence_count"`
	InReviewCount        int        `json:"in_review_count"`
	DecidedCount         int        `json:"decided_count"`
	PostponedCount       int        `json:"postponed_count"`
	DeclinedCount        int        `json:"declined_count"`
	AutoAcceptedCount    int        `json:"auto_accepted_count"`
	TopTier              int        `json:"top_tier"`
	LatestActivityAt     *time.Time `json:"latest_activity_at,omitempty"`
}

// ListUnderwritingCaseQuoteSummaries returns one row per quote that has
// at least one UnderwritingCase, with case-count breakdowns by status
// and tier. The provided UnderwritingCaseFilter (status / tier /
// assignee) is applied to the inner row set BEFORE aggregation — so
// filtering by "status=in_review" returns only quotes with at least one
// in-review case, with the counts limited to that subset.
//
// Single query, GROUP BY quote_id. Left-joined to group_pricing_quotes
// so the renderer sees the quote name, scheme name and broker without
// an N+1 lookup.
func ListUnderwritingCaseQuoteSummaries(filter UnderwritingCaseFilter) ([]UnderwritingQuoteSummary, error) {
	q := DB.Table("underwriting_cases AS uc").
		Select(`uc.quote_id AS quote_id,
            COALESCE(gpq.quote_name, '') AS quote_name,
            COALESCE(gpq.scheme_name, '') AS scheme_name,
            COALESCE(gpq.broker_name, '') AS broker_name,
            COALESCE(gpq.status, '') AS quote_status,
            COUNT(*) AS total_cases,
            SUM(CASE WHEN uc.status = ? THEN 1 ELSE 0 END) AS pending_evidence_count,
            SUM(CASE WHEN uc.status = ? THEN 1 ELSE 0 END) AS in_review_count,
            SUM(CASE WHEN uc.status = ? THEN 1 ELSE 0 END) AS decided_count,
            SUM(CASE WHEN uc.status = ? THEN 1 ELSE 0 END) AS postponed_count,
            SUM(CASE WHEN uc.status = ? THEN 1 ELSE 0 END) AS declined_count,
            SUM(CASE WHEN uc.status = ? THEN 1 ELSE 0 END) AS auto_accepted_count,
            MAX(uc.tier) AS top_tier,
            MAX(uc.updated_at) AS latest_activity_at`,
			models.UWCaseStatusPendingEvidence,
			models.UWCaseStatusInReview,
			models.UWCaseStatusDecided,
			models.UWCaseStatusPostponed,
			models.UWCaseStatusDeclined,
			models.UWCaseStatusAutoAccepted,
		).
		Joins("LEFT JOIN group_pricing_quotes gpq ON gpq.id = uc.quote_id").
		Group("uc.quote_id, gpq.quote_name, gpq.scheme_name, gpq.broker_name, gpq.status").
		Order("MAX(uc.updated_at) DESC, SUM(CASE WHEN uc.status = '" + string(models.UWCaseStatusPendingEvidence) + "' THEN 1 ELSE 0 END) DESC")

	if filter.QuoteID != nil {
		q = q.Where("uc.quote_id = ?", *filter.QuoteID)
	}
	if filter.Status != "" {
		q = q.Where("uc.status = ?", filter.Status)
	}
	if filter.Tier != nil {
		q = q.Where("uc.tier = ?", *filter.Tier)
	}
	if filter.AssignedUnderwriterEmail != "" {
		q = q.Where("uc.assigned_underwriter_email = ?", filter.AssignedUnderwriterEmail)
	}

	var summaries []UnderwritingQuoteSummary
	if err := q.Scan(&summaries).Error; err != nil {
		return nil, err
	}
	return summaries, nil
}

// ListUnderwritingCases returns cases matching the filter, ordered by tier
// descending then creation date descending so the queue surfaces the
// highest-risk newest items first.
func ListUnderwritingCases(filter UnderwritingCaseFilter) ([]models.UnderwritingCase, error) {
	q := DB.Model(&models.UnderwritingCase{})
	if filter.QuoteID != nil {
		q = q.Where("quote_id = ?", *filter.QuoteID)
	}
	if filter.Status != "" {
		q = q.Where("status = ?", filter.Status)
	}
	if filter.Tier != nil {
		q = q.Where("tier = ?", *filter.Tier)
	}
	if filter.AssignedUnderwriterEmail != "" {
		q = q.Where("assigned_underwriter_email = ?", filter.AssignedUnderwriterEmail)
	}
	var cases []models.UnderwritingCase
	if err := q.Order("tier DESC, creation_date DESC").Find(&cases).Error; err != nil {
		return nil, err
	}
	return cases, nil
}

// GetUnderwritingCase returns the case with its decisions, events, and
// attachments preloaded. Attachments have ViewerURL populated by the
// controller layer.
func GetUnderwritingCase(caseID int) (*models.UnderwritingCase, error) {
	var c models.UnderwritingCase
	err := DB.Preload("Decisions").
		Preload("Events").
		Preload("Attachments").
		Where("id = ?", caseID).
		First(&c).Error
	if err != nil {
		return nil, err
	}
	return &c, nil
}

// AssignUnderwritingCase sets the assignee. Does not change status.
func AssignUnderwritingCase(caseID int, underwriterEmail string, user models.AppUser) (*models.UnderwritingCase, error) {
	var c models.UnderwritingCase
	if err := DB.Where("id = ?", caseID).First(&c).Error; err != nil {
		return nil, err
	}
	c.AssignedUnderwriterEmail = underwriterEmail
	if err := DB.Save(&c).Error; err != nil {
		return nil, err
	}
	recordCaseEvent(c.ID, "assigned", user.UserEmail, map[string]any{
		"assignee": underwriterEmail,
	})
	if underwriterEmail != "" {
		_, _ = CreateNotification(models.CreateNotificationRequest{
			RecipientEmail: underwriterEmail,
			SenderEmail:    user.UserEmail,
			SenderName:     user.UserName,
			Type:           NotificationTypeUWCaseCreated,
			Title:          "Underwriting case assigned to you",
			Body:           fmt.Sprintf("%s (%s) — tier %d", c.MemberName, c.Category, c.Tier),
			ObjectType:     "underwriting_case",
			ObjectID:       c.ID,
		})
	}
	return &c, nil
}

// TransitionUnderwritingCase moves a case to a new status. Returns an error
// if the transition is not in allowedUWCaseTransitions. Records an event and,
// when the case has an assigned underwriter, sends a notification.
func TransitionUnderwritingCase(caseID int, newStatus models.UWCaseStatus, user models.AppUser, note string) (*models.UnderwritingCase, error) {
	var c models.UnderwritingCase
	if err := DB.Where("id = ?", caseID).First(&c).Error; err != nil {
		return nil, err
	}
	allowed, ok := allowedUWCaseTransitions[c.Status]
	if !ok || !allowed[newStatus] {
		return nil, fmt.Errorf("transition %s -> %s not allowed", c.Status, newStatus)
	}
	previous := c.Status
	c.Status = newStatus
	if newStatus == models.UWCaseStatusDecided || newStatus == models.UWCaseStatusDeclined {
		now := time.Now()
		c.DecidedAt = &now
		c.DecidedBy = user.UserEmail
	}
	if err := DB.Save(&c).Error; err != nil {
		return nil, err
	}
	recordCaseEvent(c.ID, "status_changed", user.UserEmail, map[string]any{
		"from": previous,
		"to":   newStatus,
		"note": note,
	})
	return &c, nil
}

// CreateUnderwritingDecision appends a decision for a benefit. The latest row
// per benefit is the operative outcome; we do not delete prior rows so the
// audit trail is preserved. Fires uw_case_decided when the case has an
// assignee.
func CreateUnderwritingDecision(caseID int, decision models.UnderwritingDecision, user models.AppUser) (*models.UnderwritingDecision, error) {
	var c models.UnderwritingCase
	if err := DB.Where("id = ?", caseID).First(&c).Error; err != nil {
		return nil, err
	}
	decision.CaseID = caseID
	decision.CreatedBy = user.UserEmail
	if err := DB.Create(&decision).Error; err != nil {
		return nil, err
	}
	recordCaseEvent(c.ID, "decision_added", user.UserEmail, map[string]any{
		"benefit_type":    decision.BenefitType,
		"outcome":         decision.Outcome,
		"loading_percent": decision.LoadingPercent,
		"exclusion_code":  decision.ExclusionCode,
		"cover_cap":       decision.CoverCap,
	})
	// Phase 4: trigger a re-rate so the broker sees the revised premium
	// without a fresh quote run. Failure is logged on the case event but
	// does not fail the decision write — re-rate can be retried manually
	// via the recreate-cases endpoint.
	if _, err := ApplyDecisionsAndReRate(c.QuoteID, user, fmt.Sprintf("Decision on case %d: %s %s", c.ID, decision.Outcome, decision.BenefitType), c.ID); err != nil {
		recordCaseEvent(c.ID, "rerate_failed", user.UserEmail, map[string]any{"error": err.Error()})
	}
	if c.AssignedUnderwriterEmail != "" && c.AssignedUnderwriterEmail != user.UserEmail {
		_, _ = CreateNotification(models.CreateNotificationRequest{
			RecipientEmail: c.AssignedUnderwriterEmail,
			SenderEmail:    user.UserEmail,
			SenderName:     user.UserName,
			Type:           NotificationTypeUWCaseDecided,
			Title:          fmt.Sprintf("Decision recorded on case %d", c.ID),
			Body:           fmt.Sprintf("%s on %s for %s", decision.Outcome, decision.BenefitType, c.MemberName),
			ObjectType:     "underwriting_case",
			ObjectID:       c.ID,
		})
	}
	return &decision, nil
}

// AppendCaseAttachments persists multipart-uploaded files for a case and
// records an event per file. `kind` for each file is read from the matching
// `kinds` form value at the same index (defaulting to models.UWAttachmentKindMedicalReport).
func AppendCaseAttachments(caseID int, files map[string][]*multipart.FileHeader, formValues map[string][]string, user models.AppUser) ([]models.UnderwritingCaseAttachment, error) {
	var c models.UnderwritingCase
	if err := DB.Where("id = ?", caseID).First(&c).Error; err != nil {
		return nil, err
	}
	baseDir := filepath.Join(uwCaseAttachmentDir, fmt.Sprintf("case_%d", caseID))
	if err := os.MkdirAll(baseDir, 0o755); err != nil {
		return nil, fmt.Errorf("create attachment dir: %w", err)
	}

	kinds := formValues["kind"]
	created := make([]models.UnderwritingCaseAttachment, 0)
	idx := 0
	for _, fhList := range files {
		for _, fh := range fhList {
			kind := models.UWAttachmentKindMedicalReport
			if idx < len(kinds) && kinds[idx] != "" {
				kind = kinds[idx]
			}
			idx++

			name := filepath.Base(fh.Filename)
			destPath := filepath.Join(baseDir, name)
			if err := saveMultipartFile(fh, destPath); err != nil {
				return nil, fmt.Errorf("save %s: %w", name, err)
			}
			att := models.UnderwritingCaseAttachment{
				CaseID:      caseID,
				Kind:        kind,
				FileName:    name,
				ContentType: fh.Header.Get("Content-Type"),
				SizeBytes:   fh.Size,
				StoragePath: destPath,
				UploadedBy:  user.UserEmail,
			}
			if err := DB.Create(&att).Error; err != nil {
				return nil, fmt.Errorf("persist attachment: %w", err)
			}
			recordCaseEvent(caseID, "attachment_added", user.UserEmail, map[string]any{
				"kind":     kind,
				"filename": name,
				"size":     fh.Size,
			})
			created = append(created, att)
		}
	}
	return created, nil
}

func saveMultipartFile(fh *multipart.FileHeader, destPath string) error {
	src, err := fh.Open()
	if err != nil {
		return err
	}
	defer src.Close()
	dst, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer dst.Close()
	_, err = io.Copy(dst, src)
	return err
}

func recordCaseEvent(caseID int, eventType, actor string, payload map[string]any) {
	body, _ := json.Marshal(payload)
	evt := models.UnderwritingCaseEvent{
		CaseID:    caseID,
		EventType: eventType,
		Actor:     actor,
		Payload:   string(body),
	}
	_ = DB.Create(&evt).Error
}

func emitCaseCreatedNotification(c models.UnderwritingCase, actor string) {
	if c.AssignedUnderwriterEmail == "" {
		return
	}
	_, _ = CreateNotification(models.CreateNotificationRequest{
		RecipientEmail: c.AssignedUnderwriterEmail,
		SenderEmail:    actor,
		Type:           NotificationTypeUWCaseCreated,
		Title:          "New underwriting case",
		Body:           fmt.Sprintf("%s (%s) — tier %d", c.MemberName, c.Category, c.Tier),
		ObjectType:     "underwriting_case",
		ObjectID:       c.ID,
	})
}
