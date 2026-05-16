package services

import (
	appLog "api/log"
	"api/models"
	"api/services/sms"
	"context"
	"errors"
	"fmt"
	"strings"
)

// ──────────────────────────────────────────────
// Claimant notifications (Phase 4)
// ──────────────────────────────────────────────
//
// Fires after a successful ACB reconciliation. For each line that the bank
// confirmed as "paid", we:
//
//   1. Generate an IT3(a) tax certificate (best-effort; only for benefits
//      flagged tax-disclosing — see services/tax_certificates.go).
//   2. Send an SMS to the claimant's contact number if available.
//   3. Enqueue an email to the claimant (template code "claim_payment_confirmed")
//      if a contact email is available on the claim.
//
// All three steps are best-effort: any one failing logs and continues so a
// single unreachable claimant doesn't block the reconciliation summary
// returned to the user.

// notifyClaimantTemplateCode is the email template the outbox renderer
// looks up. Installs override the template content via the email-templates
// admin UI; the code itself is fixed so the service can refer to it.
const notifyClaimantTemplateCode = "claim_payment_confirmed"

// NotifyClaimantsForReconciliation walks the paid lines from a reconciliation
// run and dispatches per-line notifications + tax certificates. Called from
// ProcessBankResponse in a goroutine so the HTTP response isn't blocked.
func NotifyClaimantsForReconciliation(results []models.ACBReconciliationResult, user models.AppUser) {
	for _, r := range results {
		if !strings.EqualFold(r.Status, "paid") {
			continue
		}
		notifyClaimantForResult(r, user)
	}
}

func notifyClaimantForResult(r models.ACBReconciliationResult, user models.AppUser) {
	// Look up the source claim + the schedule item snapshot so notifications
	// match what finance actually paid (not later edits to the claim).
	var item models.ClaimPaymentScheduleItem
	if err := DB.First(&item, r.ScheduleItemID).Error; err != nil {
		appLog.WithField("error", err.Error()).
			WithField("schedule_item_id", r.ScheduleItemID).
			Warn("claimant notify: schedule item not found")
		return
	}

	var claim models.GroupSchemeClaim
	if err := DB.Select("id, claim_number, claimant_name, claimant_contact_number, member_name").
		First(&claim, item.ClaimID).Error; err != nil {
		appLog.WithField("error", err.Error()).Warn("claimant notify: claim not found")
		return
	}

	if cert, err := EnsureTaxCertificateForItem(item.ScheduleID, item.ID, user); err == nil && cert.ID != 0 {
		appLog.WithField("claim_number", claim.ClaimNumber).
			WithField("cert_ref", cert.CertificateRef).
			Info("tax certificate issued for paid claim")
	} else if err != nil {
		appLog.WithField("error", err.Error()).Warn("claimant notify: tax cert generation failed")
	}

	sendSMSConfirmation(claim, item)
	enqueueEmailConfirmation(claim, item, user)
}

func sendSMSConfirmation(claim models.GroupSchemeClaim, item models.ClaimPaymentScheduleItem) {
	provider := sms.Active()
	if provider == nil {
		return
	}
	to := strings.TrimSpace(claim.ClaimantContactNumber)
	if to == "" {
		return
	}
	body := fmt.Sprintf(
		"Hi %s, your claim %s has been paid. Net amount: R %.2f. Ref: %s.",
		firstName(coalesce(claim.ClaimantName, claim.MemberName, item.BeneficiaryName)),
		claim.ClaimNumber,
		nonZero(item.NetPayable, item.GrossAmount, item.ClaimAmount),
		claim.ClaimNumber,
	)
	if _, err := provider.Send(context.Background(), sms.Message{
		To:        to,
		Body:      body,
		Reference: claim.ClaimNumber,
	}); err != nil {
		appLog.WithField("error", err.Error()).
			WithField("claim_number", claim.ClaimNumber).
			Warn("claimant notify: SMS send failed")
	}
}

func enqueueEmailConfirmation(claim models.GroupSchemeClaim, item models.ClaimPaymentScheduleItem, user models.AppUser) {
	// Try to pull a contact email from the most recent claim communication
	// row; many installs don't capture a structured email on the claim
	// itself yet, so this is best-effort.
	var email string
	type emailRow struct{ Message string }
	var rows []emailRow
	if err := DB.Raw(`
		SELECT message FROM group_scheme_claim_communications
		WHERE claim_id = ? AND LOWER(method) = 'email'
		ORDER BY created_at DESC LIMIT 5
	`, claim.ID).Scan(&rows).Error; err == nil {
		for _, r := range rows {
			if e := extractEmail(r.Message); e != "" {
				email = e
				break
			}
		}
	}
	if email == "" {
		appLog.WithField("claim_number", claim.ClaimNumber).Debug("claimant notify: no email on file, skipping")
		return
	}

	licenseId := resolveDefaultLicense()
	if licenseId == "" {
		appLog.Debug("claimant notify: no email_settings configured, skipping email")
		return
	}
	req := EnqueueEmailRequest{
		LicenseId:    licenseId,
		TemplateCode: notifyClaimantTemplateCode,
		To:           []string{email},
		Vars: map[string]interface{}{
			"claim_number":     claim.ClaimNumber,
			"beneficiary_name": coalesce(item.BeneficiaryName, claim.MemberName, claim.ClaimantName),
			"benefit_name":     item.BenefitName,
			"gross_amount":     item.GrossAmount,
			"deductions":       item.PremiumArrearsDeduction + item.PolicyLoanDeduction + item.TaxWithheld,
			"net_amount":       item.NetPayable,
		},
		RelatedObjectType: "claim_payment_schedule_item",
		RelatedObjectID:   fmt.Sprintf("%d", item.ID),
		CreatedBy:         user.UserName,
	}
	if _, err := EnqueueEmail(req); err != nil {
		// Missing-template is normal in dev when the install hasn't seeded
		// the "claim_payment_confirmed" template yet; downgrade to info.
		if isMissingTemplate(err) {
			appLog.WithField("template", notifyClaimantTemplateCode).
				Info("claimant notify: email template not configured, skipping")
			return
		}
		appLog.WithField("error", err.Error()).
			WithField("claim_number", claim.ClaimNumber).
			Warn("claimant notify: email enqueue failed")
	}
}

func isMissingTemplate(err error) bool {
	if err == nil {
		return false
	}
	msg := strings.ToLower(err.Error())
	if strings.Contains(msg, "template") && (strings.Contains(msg, "not found") || strings.Contains(msg, "no such")) {
		return true
	}
	return errors.Is(err, ErrEmailTemplateNotFound)
}

// ErrEmailTemplateNotFound is exported so the enqueue path can return a
// distinguishable error when the requested template isn't seeded yet.
// Phase 4 declares this here rather than perturbing email_enqueue.go; the
// enqueue function returns whatever error its lookup produces and this
// helper picks it up by message-matching, so installs that pre-seed the
// template won't hit this branch.
var ErrEmailTemplateNotFound = errors.New("email template not found")

func coalesce(vals ...string) string {
	for _, v := range vals {
		if strings.TrimSpace(v) != "" {
			return v
		}
	}
	return ""
}

func nonZero(vals ...float64) float64 {
	for _, v := range vals {
		if v > 0 {
			return v
		}
	}
	return 0
}

func firstName(full string) string {
	full = strings.TrimSpace(full)
	if full == "" {
		return "there"
	}
	if i := strings.IndexAny(full, " \t"); i > 0 {
		return full[:i]
	}
	return full
}

// resolveDefaultLicense returns the singleton install's license id, derived
// from the first EmailSettings row. Single-tenant installs have exactly one
// row; multi-tenant rollouts will need a different resolver but the call
// shape here stays the same.
func resolveDefaultLicense() string {
	type r struct{ LicenseId string }
	var row r
	if err := DB.Raw(`SELECT license_id FROM email_settings ORDER BY id ASC LIMIT 1`).Scan(&row).Error; err != nil {
		return ""
	}
	return row.LicenseId
}

// extractEmail pulls the first plausible email address out of a free-text
// blob (used because the schema's "email" communication channel stores the
// body without a dedicated To field).
func extractEmail(s string) string {
	for _, tok := range strings.Fields(s) {
		t := strings.Trim(tok, "<>,;:\"'()[]")
		if strings.Count(t, "@") == 1 && strings.Contains(t, ".") {
			return t
		}
	}
	return ""
}
