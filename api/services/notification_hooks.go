package services

import (
	"api/log"
	"api/models"
	"fmt"
)

// resolveEmailByName looks up an OrgUser's email address by their name.
// Quote fields like CreatedBy, Reviewer, and ModifiedBy store display names,
// so we need this lookup to obtain the actual email for notification delivery.
func resolveEmailByName(name string) string {
	if name == "" {
		return ""
	}
	var user models.OrgUser
	if err := DB.Where("name = ?", name).First(&user).Error; err != nil {
		log.WithField("name", name).Warn("Could not resolve email for user name")
		return ""
	}
	return user.Email
}

// NotifyQuoteSubmitted notifies the reviewer that a quote has been submitted for review.
func NotifyQuoteSubmitted(quote models.GroupPricingQuote, submitter models.AppUser) {
	recipientEmail := resolveEmailByName(quote.Reviewer)
	if recipientEmail == "" {
		return
	}
	CreateNotification(models.CreateNotificationRequest{
		RecipientEmail: recipientEmail,
		SenderEmail:    submitter.UserEmail,
		SenderName:     submitter.UserName,
		Type:           "quote_submitted",
		Title:          "Quote submitted for review",
		Body:           fmt.Sprintf("Quote '%s' has been submitted for your review", quote.QuoteName),
		ObjectType:     "quote",
		ObjectID:       quote.ID,
	})
}

// NotifyQuoteApproved notifies the quote creator that their quote was approved.
func NotifyQuoteApproved(quote models.GroupPricingQuote, approver models.AppUser) {
	recipientEmail := resolveEmailByName(quote.CreatedBy)
	if recipientEmail == "" {
		return
	}
	CreateNotification(models.CreateNotificationRequest{
		RecipientEmail: recipientEmail,
		SenderEmail:    approver.UserEmail,
		SenderName:     approver.UserName,
		Type:           "quote_approved",
		Title:          "Quote approved",
		Body:           fmt.Sprintf("Quote '%s' has been approved", quote.QuoteName),
		ObjectType:     "quote",
		ObjectID:       quote.ID,
	})
}

// NotifyQuoteRejected notifies the quote creator that their quote was rejected.
func NotifyQuoteRejected(quote models.GroupPricingQuote, rejecter models.AppUser, reason string) {
	recipientEmail := resolveEmailByName(quote.CreatedBy)
	if recipientEmail == "" {
		return
	}
	body := fmt.Sprintf("Quote '%s' has been rejected", quote.QuoteName)
	if reason != "" {
		body += fmt.Sprintf(": %s", reason)
	}
	CreateNotification(models.CreateNotificationRequest{
		RecipientEmail: recipientEmail,
		SenderEmail:    rejecter.UserEmail,
		SenderName:     rejecter.UserName,
		Type:           "quote_rejected",
		Title:          "Quote rejected",
		Body:           body,
		ObjectType:     "quote",
		ObjectID:       quote.ID,
	})
}

// NotifyQuoteAccepted notifies the quote creator that their quote was accepted.
func NotifyQuoteAccepted(quote models.GroupPricingQuote, accepter models.AppUser) {
	recipientEmail := resolveEmailByName(quote.CreatedBy)
	if recipientEmail == "" {
		return
	}
	CreateNotification(models.CreateNotificationRequest{
		RecipientEmail: recipientEmail,
		SenderEmail:    accepter.UserEmail,
		SenderName:     accepter.UserName,
		Type:           "quote_accepted",
		Title:          "Quote accepted",
		Body:           fmt.Sprintf("Quote '%s' has been accepted", quote.QuoteName),
		ObjectType:     "quote",
		ObjectID:       quote.ID,
	})
}

// NotifySchemeStatusChange notifies relevant users when a scheme status changes.
func NotifySchemeStatusChange(recipientEmail string, schemeName string, schemeID int, newStatus string, changedBy models.AppUser) {
	if recipientEmail == "" || recipientEmail == changedBy.UserEmail {
		return
	}
	CreateNotification(models.CreateNotificationRequest{
		RecipientEmail: recipientEmail,
		SenderEmail:    changedBy.UserEmail,
		SenderName:     changedBy.UserName,
		Type:           "scheme_status_change",
		Title:          "Scheme status updated",
		Body:           fmt.Sprintf("Scheme '%s' status changed to %s", schemeName, newStatus),
		ObjectType:     "scheme",
		ObjectID:       schemeID,
	})
}

// notifyIfNotSender is a helper that resolves recipientName to an email and sends a notification,
// skipping if the recipient is the sender or cannot be resolved.
func notifyIfNotSender(recipientName string, sender models.AppUser, req models.CreateNotificationRequest) {
	email := resolveEmailByName(recipientName)
	if email == "" || email == sender.UserEmail {
		return
	}
	req.RecipientEmail = email
	CreateNotification(req)
}

// ─── Premium Schedule Hooks ────────────────────────────────────────────────

func NotifyScheduleReviewed(s models.PremiumSchedule, reviewer models.AppUser) {
	notifyIfNotSender(s.GeneratedBy, reviewer, models.CreateNotificationRequest{
		SenderEmail: reviewer.UserEmail,
		SenderName:  reviewer.UserName,
		Type:        "schedule_reviewed",
		Title:       "Premium schedule reviewed",
		Body:        fmt.Sprintf("Schedule for '%s' has been reviewed", s.SchemeName),
		ObjectType:  "premium_schedule",
		ObjectID:    s.ID,
	})
}

func NotifyScheduleApproved(s models.PremiumSchedule, approver models.AppUser) {
	base := models.CreateNotificationRequest{
		SenderEmail: approver.UserEmail,
		SenderName:  approver.UserName,
		Type:        "schedule_approved",
		Title:       "Premium schedule approved",
		Body:        fmt.Sprintf("Schedule for '%s' has been approved", s.SchemeName),
		ObjectType:  "premium_schedule",
		ObjectID:    s.ID,
	}
	notifyIfNotSender(s.GeneratedBy, approver, base)
	notifyIfNotSender(s.ReviewedBy, approver, base)
}

func NotifyScheduleFinalized(s models.PremiumSchedule, finalizer models.AppUser) {
	notifyIfNotSender(s.GeneratedBy, finalizer, models.CreateNotificationRequest{
		SenderEmail: finalizer.UserEmail,
		SenderName:  finalizer.UserName,
		Type:        "schedule_finalized",
		Title:       "Premium schedule finalized",
		Body:        fmt.Sprintf("Schedule for '%s' has been finalized", s.SchemeName),
		ObjectType:  "premium_schedule",
		ObjectID:    s.ID,
	})
}

func NotifyScheduleVoided(s models.PremiumSchedule, voider models.AppUser) {
	body := fmt.Sprintf("Schedule for '%s' has been voided", s.SchemeName)
	if s.VoidReason != "" {
		body += fmt.Sprintf(": %s", s.VoidReason)
	}
	notifyIfNotSender(s.GeneratedBy, voider, models.CreateNotificationRequest{
		SenderEmail: voider.UserEmail,
		SenderName:  voider.UserName,
		Type:        "schedule_voided",
		Title:       "Premium schedule voided",
		Body:        body,
		ObjectType:  "premium_schedule",
		ObjectID:    s.ID,
	})
}

func NotifyScheduleCancelled(s models.PremiumSchedule, canceller models.AppUser) {
	notifyIfNotSender(s.GeneratedBy, canceller, models.CreateNotificationRequest{
		SenderEmail: canceller.UserEmail,
		SenderName:  canceller.UserName,
		Type:        "schedule_cancelled",
		Title:       "Premium schedule cancelled",
		Body:        fmt.Sprintf("Schedule for '%s' has been cancelled", s.SchemeName),
		ObjectType:  "premium_schedule",
		ObjectID:    s.ID,
	})
}

func NotifySchemeSuspended(scheme models.GroupScheme, suspender models.AppUser) {
	notifyIfNotSender(scheme.CreatedBy, suspender, models.CreateNotificationRequest{
		SenderEmail: suspender.UserEmail,
		SenderName:  suspender.UserName,
		Type:        "scheme_suspended",
		Title:       "Scheme suspended",
		Body:        fmt.Sprintf("Scheme '%s' has been suspended", scheme.Name),
		ObjectType:  "scheme",
		ObjectID:    scheme.ID,
	})
}

func NotifySchemeReinstated(scheme models.GroupScheme, reinstater models.AppUser) {
	notifyIfNotSender(scheme.CreatedBy, reinstater, models.CreateNotificationRequest{
		SenderEmail: reinstater.UserEmail,
		SenderName:  reinstater.UserName,
		Type:        "scheme_reinstated",
		Title:       "Scheme reinstated",
		Body:        fmt.Sprintf("Scheme '%s' has been reinstated", scheme.Name),
		ObjectType:  "scheme",
		ObjectID:    scheme.ID,
	})
}

// ─── Bordereaux Inbound Hooks ──────────────────────────────────────────────

func NotifySubmissionReviewed(sub models.EmployerSubmission, reviewer models.AppUser) {
	notifyIfNotSender(sub.SubmittedBy, reviewer, models.CreateNotificationRequest{
		SenderEmail: reviewer.UserEmail,
		SenderName:  reviewer.UserName,
		Type:        "submission_reviewed",
		Title:       "Submission under review",
		Body:        fmt.Sprintf("Submission for '%s' is now under review", sub.SchemeName),
		ObjectType:  "employer_submission",
		ObjectID:    sub.ID,
	})
}

func NotifySubmissionQueryRaised(sub models.EmployerSubmission, querier models.AppUser) {
	body := fmt.Sprintf("A query has been raised on submission for '%s'", sub.SchemeName)
	if sub.QueryNotes != "" {
		body += fmt.Sprintf(": %s", sub.QueryNotes)
	}
	notifyIfNotSender(sub.SubmittedBy, querier, models.CreateNotificationRequest{
		SenderEmail: querier.UserEmail,
		SenderName:  querier.UserName,
		Type:        "submission_query_raised",
		Title:       "Submission query raised",
		Body:        body,
		ObjectType:  "employer_submission",
		ObjectID:    sub.ID,
	})
}

func NotifySubmissionAccepted(sub models.EmployerSubmission, accepter models.AppUser) {
	notifyIfNotSender(sub.SubmittedBy, accepter, models.CreateNotificationRequest{
		SenderEmail: accepter.UserEmail,
		SenderName:  accepter.UserName,
		Type:        "submission_accepted",
		Title:       "Submission accepted",
		Body:        fmt.Sprintf("Submission for '%s' has been accepted", sub.SchemeName),
		ObjectType:  "employer_submission",
		ObjectID:    sub.ID,
	})
}

func NotifySubmissionRejected(sub models.EmployerSubmission, rejecter models.AppUser, reason string) {
	body := fmt.Sprintf("Submission for '%s' has been rejected", sub.SchemeName)
	if reason != "" {
		body += fmt.Sprintf(": %s", reason)
	}
	notifyIfNotSender(sub.SubmittedBy, rejecter, models.CreateNotificationRequest{
		SenderEmail: rejecter.UserEmail,
		SenderName:  rejecter.UserName,
		Type:        "submission_rejected",
		Title:       "Submission rejected",
		Body:        body,
		ObjectType:  "employer_submission",
		ObjectID:    sub.ID,
	})
}

// ─── Bordereaux Outbound Hooks ─────────────────────────────────────────────

func NotifyBordereauxReviewed(brd models.GeneratedBordereaux, reviewer models.AppUser) {
	notifyIfNotSender(brd.CreatedBy, reviewer, models.CreateNotificationRequest{
		SenderEmail: reviewer.UserEmail,
		SenderName:  reviewer.UserName,
		Type:        "bordereaux_reviewed",
		Title:       "Bordereaux reviewed",
		Body:        fmt.Sprintf("Bordereaux for '%s' has been reviewed", brd.SchemeName),
		ObjectType:  "bordereaux",
		ObjectID:    brd.ID,
	})
}

func NotifyBordereauxApproved(brd models.GeneratedBordereaux, approver models.AppUser) {
	base := models.CreateNotificationRequest{
		SenderEmail: approver.UserEmail,
		SenderName:  approver.UserName,
		Type:        "bordereaux_approved",
		Title:       "Bordereaux approved",
		Body:        fmt.Sprintf("Bordereaux for '%s' has been approved", brd.SchemeName),
		ObjectType:  "bordereaux",
		ObjectID:    brd.ID,
	}
	notifyIfNotSender(brd.CreatedBy, approver, base)
	notifyIfNotSender(brd.ReviewedBy, approver, base)
}

func NotifyBordereauxSubmitted(brd models.GeneratedBordereaux, submitter models.AppUser) {
	notifyIfNotSender(brd.CreatedBy, submitter, models.CreateNotificationRequest{
		SenderEmail: submitter.UserEmail,
		SenderName:  submitter.UserName,
		Type:        "bordereaux_submitted",
		Title:       "Bordereaux submitted",
		Body:        fmt.Sprintf("Bordereaux for '%s' has been submitted to insurer", brd.SchemeName),
		ObjectType:  "bordereaux",
		ObjectID:    brd.ID,
	})
}

// ─── Reinsurance Cession Hooks ─────────────────────────────────────────────

func NotifyRIBordereauxSubmitted(run models.RIBordereauxRun, submitter models.AppUser) {
	notifyIfNotSender(run.GeneratedBy, submitter, models.CreateNotificationRequest{
		SenderEmail: submitter.UserEmail,
		SenderName:  submitter.UserName,
		Type:        "ri_bordereaux_submitted",
		Title:       "RI bordereaux submitted",
		Body:        fmt.Sprintf("RI bordereaux for '%s' has been submitted", run.ReinsurerName),
		ObjectType:  "ri_bordereaux",
		ObjectID:    run.ID,
	})
}

func NotifyRIBordereauxAcknowledged(run models.RIBordereauxRun, acknowledger models.AppUser) {
	notifyIfNotSender(run.SubmittedBy, acknowledger, models.CreateNotificationRequest{
		SenderEmail: acknowledger.UserEmail,
		SenderName:  acknowledger.UserName,
		Type:        "ri_bordereaux_acknowledged",
		Title:       "RI bordereaux acknowledged",
		Body:        fmt.Sprintf("RI bordereaux for '%s' has been acknowledged", run.ReinsurerName),
		ObjectType:  "ri_bordereaux",
		ObjectID:    run.ID,
	})
}

// ─── CSM Engine Hooks ──────────────────────────────────────────────────────

func NotifyCsmRunReviewed(run models.CsmRun, reviewer models.AppUser) {
	if run.UserEmail == "" || run.UserEmail == reviewer.UserEmail {
		return
	}
	CreateNotification(models.CreateNotificationRequest{
		RecipientEmail: run.UserEmail,
		SenderEmail:    reviewer.UserEmail,
		SenderName:     reviewer.UserName,
		Type:           "csm_run_reviewed",
		Title:          "CSM run reviewed",
		Body:           fmt.Sprintf("CSM run '%s' has been reviewed", run.Name),
		ObjectType:     "csm_run",
		ObjectID:       run.ID,
	})
}

func NotifyCsmRunApproved(run models.CsmRun, approver models.AppUser) {
	// Notify the original creator
	if run.UserEmail != "" && run.UserEmail != approver.UserEmail {
		CreateNotification(models.CreateNotificationRequest{
			RecipientEmail: run.UserEmail,
			SenderEmail:    approver.UserEmail,
			SenderName:     approver.UserName,
			Type:           "csm_run_approved",
			Title:          "CSM run approved",
			Body:           fmt.Sprintf("CSM run '%s' has been approved", run.Name),
			ObjectType:     "csm_run",
			ObjectID:       run.ID,
		})
	}
	// Notify the reviewer
	notifyIfNotSender(run.ReviewedBy, approver, models.CreateNotificationRequest{
		SenderEmail: approver.UserEmail,
		SenderName:  approver.UserName,
		Type:        "csm_run_approved",
		Title:       "CSM run approved",
		Body:        fmt.Sprintf("CSM run '%s' has been approved", run.Name),
		ObjectType:  "csm_run",
		ObjectID:    run.ID,
	})
}

// ─── IFRS17 Amendment Hooks ────────────────────────────────────────────────

func NotifyAmendmentApproved(amendment models.IFRS17Amendment, approver models.AppUser) {
	notifyIfNotSender(amendment.CreatedBy, approver, models.CreateNotificationRequest{
		SenderEmail: approver.UserEmail,
		SenderName:  approver.UserName,
		Type:        "amendment_approved",
		Title:       "IFRS17 amendment approved",
		Body:        fmt.Sprintf("IFRS17 amendment #%d has been approved", amendment.ID),
		ObjectType:  "ifrs17_amendment",
		ObjectID:    amendment.ID,
	})
}

// ─── Reinsurance Settlement Hooks ──────────────────────────────────────────

func NotifySettlementUpdated(account models.TechnicalAccount, updater models.AppUser, newStatus string) {
	notifyIfNotSender(account.CreatedBy, updater, models.CreateNotificationRequest{
		SenderEmail: updater.UserEmail,
		SenderName:  updater.UserName,
		Type:        "settlement_updated",
		Title:       "Settlement account updated",
		Body:        fmt.Sprintf("Technical account '%s' (%s) status changed to %s", account.AccountNumber, account.ReinsurerName, newStatus),
		ObjectType:  "technical_account",
		ObjectID:    account.ID,
	})
}

func NotifySettlementDisputeResolved(account models.TechnicalAccount, resolver models.AppUser) {
	notifyIfNotSender(account.CreatedBy, resolver, models.CreateNotificationRequest{
		SenderEmail: resolver.UserEmail,
		SenderName:  resolver.UserName,
		Type:        "settlement_dispute_resolved",
		Title:       "Settlement dispute resolved",
		Body:        fmt.Sprintf("Dispute on technical account '%s' (%s) has been resolved", account.AccountNumber, account.ReinsurerName),
		ObjectType:  "technical_account",
		ObjectID:    account.ID,
	})
}

// ─── Claim Payment Hooks ───────────────────────────────────────────────────

// NotifyCustomTieredTableRequested notifies all system admins that a user has
// requested a custom tiered income replacement table for a scheme.
func NotifyCustomTieredTableRequested(schemeName string, schemeID int, riskRateCode string, requester models.AppUser) {
	// Find all users with system:admin permission
	var adminUsers []models.OrgUser
	DB.Joins("JOIN gp_user_roles ON org_users.gp_role_id = gp_user_roles.id").
		Joins("JOIN gp_user_role_permissions ON gp_user_roles.id = gp_user_role_permissions.gp_user_role_id").
		Joins("JOIN gp_permissions ON gp_user_role_permissions.gp_permission_id = gp_permissions.id").
		Where("gp_permissions.slug = ?", "system:admin").
		Find(&adminUsers)

	for _, admin := range adminUsers {
		if admin.Email == "" {
			continue
		}
		CreateNotification(models.CreateNotificationRequest{
			RecipientEmail: admin.Email,
			SenderEmail:    requester.UserEmail,
			SenderName:     requester.UserName,
			Type:           "custom_tir_requested",
			Title:          "Custom tiered income replacement table requested",
			Body:           fmt.Sprintf("A custom tiered income replacement table has been requested for scheme '%s' (risk rate code: %s) by %s. Please upload the custom table.", schemeName, riskRateCode, requester.UserName),
			ObjectType:     "scheme",
			ObjectID:       schemeID,
		})
	}
}

// NotifyCustomTieredTableUploaded notifies the requesting user that the custom
// tiered income replacement table has been uploaded and they can proceed.
func NotifyCustomTieredTableUploaded(schemeName string, riskRateCode string, uploader models.AppUser) {
	// Find the most recent quote that references this scheme name to identify the requester
	var quote models.GroupPricingQuote
	err := DB.Where("scheme_name = ?", schemeName).Order("id desc").First(&quote).Error
	if err != nil {
		log.WithField("scheme_name", schemeName).Warn("Could not find quote for custom tiered table notification")
		return
	}
	recipientEmail := resolveEmailByName(quote.CreatedBy)
	if recipientEmail == "" || recipientEmail == uploader.UserEmail {
		return
	}
	CreateNotification(models.CreateNotificationRequest{
		RecipientEmail: recipientEmail,
		SenderEmail:    uploader.UserEmail,
		SenderName:     uploader.UserName,
		Type:           "custom_tir_uploaded",
		Title:          "Custom tiered income replacement table uploaded",
		Body:           fmt.Sprintf("The custom tiered income replacement table for scheme '%s' (risk rate code: %s) has been uploaded. You may now proceed with the calculation.", schemeName, riskRateCode),
		ObjectType:     "scheme",
		ObjectID:       quote.SchemeID,
	})
}

func NotifyClaimPaymentSummary(summary models.ACBReconciliationSummary, user models.AppUser) {
	if user.UserEmail == "" {
		return
	}
	CreateNotification(models.CreateNotificationRequest{
		RecipientEmail: user.UserEmail,
		SenderEmail:    user.UserEmail,
		SenderName:     user.UserName,
		Type:           "claim_payment_summary",
		Title:          "Bank reconciliation complete",
		Body:           fmt.Sprintf("Reconciliation complete: %d paid, %d failed out of %d transactions", summary.Paid, summary.Failed, summary.TotalTransactions),
		ObjectType:     "claim_payment",
		ObjectID:       0,
	})
}
