package services

import (
	"api/log"
	"api/models"
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
)

// Sentinel errors so controllers can map service-layer failures to the
// right HTTP status without string-matching.
var (
	ErrUserFlagNotFound       = errors.New("user flag not found")
	ErrUserFlagAlreadyOpen    = errors.New("user already has an open flag with this reason")
	ErrUserFlagAlreadyResolved = errors.New("user flag is already resolved")
	ErrUserFlagInvalidReason  = errors.New("flag_reason must be coaching or capacity")
	ErrUserFlagSelfFlag       = errors.New("you cannot flag yourself")
	ErrUserFlagNoteTooShort   = errors.New("note must be at least 10 characters")
)

const userFlagMinNoteLen = 10

// ListUserFlags returns flags matching the filter, ordered by opened_at
// descending so the most recent activity sits at the top. The default
// status is "open" — admins typically want the active list first; "all"
// or "resolved" is opt-in.
func ListUserFlags(filter models.UserFlagsFilter) ([]models.QuoteUserFlag, error) {
	var flags []models.QuoteUserFlag
	q := DB.Order("opened_at DESC")

	status := strings.ToLower(strings.TrimSpace(filter.Status))
	switch status {
	case "", "open":
		q = q.Where("resolved_at IS NULL")
	case "resolved":
		q = q.Where("resolved_at IS NOT NULL")
	case "all":
		// no constraint
	default:
		return nil, fmt.Errorf("invalid status %q (open|resolved|all)", filter.Status)
	}

	if name := strings.TrimSpace(filter.UserName); name != "" {
		q = q.Where("user_name = ?", name)
	}
	if reason := strings.TrimSpace(filter.Reason); reason != "" {
		q = q.Where("flag_reason = ?", reason)
	}

	if err := q.Find(&flags).Error; err != nil {
		return nil, err
	}
	return flags, nil
}

// OpenUserFlag creates a new open flag. Enforces: actor != target,
// reason ∈ {coaching, capacity}, note length, and uniqueness of open
// flags per (user, reason). Fires a neutral in-app notification to the
// flagged user when their email is resolvable; the manager's internal
// note is never exposed to the recipient.
func OpenUserFlag(ctx context.Context, req models.OpenUserFlagRequest, actor models.AppUser) (models.QuoteUserFlag, error) {
	var flag models.QuoteUserFlag

	reason := strings.ToLower(strings.TrimSpace(req.FlagReason))
	if reason != models.QuoteUserFlagReasonCoaching && reason != models.QuoteUserFlagReasonCapacity {
		return flag, ErrUserFlagInvalidReason
	}

	targetName := strings.TrimSpace(req.UserName)
	if targetName == "" {
		return flag, fmt.Errorf("user_name is required")
	}
	note := strings.TrimSpace(req.Note)
	if len(note) < userFlagMinNoteLen {
		return flag, ErrUserFlagNoteTooShort
	}

	// Resolve target email: prefer the client-supplied value, fall back to
	// the org_users lookup by name. Email is optional — if we can't find
	// it, the flag still opens but the notification step is skipped.
	targetEmail := strings.TrimSpace(req.UserEmail)
	if targetEmail == "" {
		var ou models.OrgUser
		if err := DB.Where("name = ?", targetName).First(&ou).Error; err == nil {
			targetEmail = ou.Email
		}
	}

	// Self-flag check, by either email or display name — name is the
	// primary identifier on the leaderboard, but the JWT carries the
	// email for the active user.
	if targetEmail != "" && strings.EqualFold(targetEmail, actor.UserEmail) {
		return flag, ErrUserFlagSelfFlag
	}
	if strings.EqualFold(targetName, actor.UserName) {
		return flag, ErrUserFlagSelfFlag
	}

	// Uniqueness of open flags per (user_name, flag_reason). Application
	// level because MySQL doesn't support partial unique indexes.
	var existing models.QuoteUserFlag
	err := DB.Where("user_name = ? AND flag_reason = ? AND resolved_at IS NULL", targetName, reason).
		First(&existing).Error
	if err == nil {
		return existing, ErrUserFlagAlreadyOpen
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return flag, err
	}

	flag = models.QuoteUserFlag{
		UserName:     targetName,
		UserEmail:    targetEmail,
		FlagReason:   reason,
		Note:         note,
		OpenedBy:     actor.UserEmail,
		OpenedByName: actor.UserName,
		OpenedAt:     time.Now(),
	}
	if err := DB.Create(&flag).Error; err != nil {
		return flag, err
	}

	// Fire-and-forget notification. The note is *not* included — the
	// flagged user only sees a neutral prompt to follow up with the
	// manager. Errors are logged, never returned: a delivery failure
	// shouldn't fail the flag write.
	if targetEmail != "" {
		notifType := "user_flagged_" + reason // user_flagged_coaching | user_flagged_capacity
		title, body := userFlagNotificationCopy(reason, actor.UserName)
		_, nerr := CreateNotification(models.CreateNotificationRequest{
			RecipientEmail: targetEmail,
			SenderEmail:    actor.UserEmail,
			SenderName:     actor.UserName,
			Type:           notifType,
			Title:          title,
			Body:           body,
			ObjectType:     "quote_user_flag",
			ObjectID:       flag.ID,
		})
		if nerr != nil {
			log.WithContext(ctx).
				WithField("error", nerr.Error()).
				WithField("recipient", targetEmail).
				Warn("Failed to send user-flag notification")
		}
	}

	return flag, nil
}

// userFlagNotificationCopy returns the neutral title and body used for
// the flagged user's notification. The manager's internal note is never
// surfaced — the recipient only sees this prompt to follow up.
func userFlagNotificationCopy(reason, managerName string) (string, string) {
	switch reason {
	case models.QuoteUserFlagReasonCoaching:
		return "Your manager would like to talk",
			fmt.Sprintf("%s has flagged your recent quote turnaround for a coaching conversation. Please reach out to them at your earliest convenience.", managerName)
	case models.QuoteUserFlagReasonCapacity:
		return "Your manager would like to review your workload",
			fmt.Sprintf("%s has flagged your current workload for review. Please reach out to them at your earliest convenience.", managerName)
	}
	return "Your manager would like to talk",
		fmt.Sprintf("%s has flagged you for review. Please reach out to them at your earliest convenience.", managerName)
}

// ResolveUserFlag closes an open flag. Errors with ErrUserFlagNotFound
// or ErrUserFlagAlreadyResolved so the controller can return 404/409.
// No notification fires on resolve — see plan §Notifications: the
// manager has already had the conversation by then.
func ResolveUserFlag(id int, resolutionNote string, actor models.AppUser) (models.QuoteUserFlag, error) {
	var flag models.QuoteUserFlag

	note := strings.TrimSpace(resolutionNote)
	if len(note) < userFlagMinNoteLen {
		return flag, ErrUserFlagNoteTooShort
	}

	if err := DB.First(&flag, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return flag, ErrUserFlagNotFound
		}
		return flag, err
	}
	if flag.ResolvedAt != nil {
		return flag, ErrUserFlagAlreadyResolved
	}

	now := time.Now()
	resolvedBy := actor.UserEmail
	resolvedByName := actor.UserName
	flag.ResolvedAt = &now
	flag.ResolvedBy = &resolvedBy
	flag.ResolvedByName = &resolvedByName
	flag.ResolutionNote = &note

	if err := DB.Save(&flag).Error; err != nil {
		return flag, err
	}
	return flag, nil
}

// OpenFlagsByUserName returns, for each user_name in the input, the list
// of currently open flags. Used by the dashboard KPI builder to attach
// flag state to leaderboard rows in one DB hit (no N+1).
func OpenFlagsByUserName(userNames []string) (map[string][]models.QuoteUserFlag, error) {
	out := make(map[string][]models.QuoteUserFlag, len(userNames))
	if len(userNames) == 0 {
		return out, nil
	}

	var rows []models.QuoteUserFlag
	if err := DB.Where("user_name IN ? AND resolved_at IS NULL", userNames).
		Order("opened_at DESC").
		Find(&rows).Error; err != nil {
		return nil, err
	}
	for _, r := range rows {
		out[r.UserName] = append(out[r.UserName], r)
	}
	return out, nil
}

// StripFlagNotes redacts the internal manager note from each flag —
// used by the public-facing leaderboard endpoint so non-managers see
// flag state (the chip) but not the confidential observations.
func StripFlagNotes(flags []models.QuoteUserFlag) []models.QuoteUserFlag {
	out := make([]models.QuoteUserFlag, len(flags))
	for i, f := range flags {
		f.Note = ""
		if f.ResolutionNote != nil {
			f.ResolutionNote = nil
		}
		out[i] = f
	}
	return out
}
