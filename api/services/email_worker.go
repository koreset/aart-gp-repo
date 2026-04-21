package services

import (
	"context"
	"encoding/json"
	appLog "api/log"
	"api/models"
	"api/services/crypto"
	"api/services/email"
	"math/rand"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

// emailWorker concurrency bounds — lower than calculations because each send
// is IO-bound on an external SMTP server. 5 is plenty for Phase 1.
const maxConcurrentEmailSends = 5

var (
	emailWorkerMu       sync.Mutex
	emailWorkerRunning  bool
	emailSemaphore      chan struct{}
)

// StartEmailOutboxWorker starts the background worker that drains the
// email_outbox table. Idempotent — calling it multiple times is a no-op after
// the first start, mirroring StartCalculationJobWorker.
func StartEmailOutboxWorker() {
	emailWorkerMu.Lock()
	defer emailWorkerMu.Unlock()

	if emailWorkerRunning {
		appLog.Info("Email outbox worker is already running")
		return
	}
	emailSemaphore = make(chan struct{}, maxConcurrentEmailSends)
	emailWorkerRunning = true

	go func() {
		appLog.Info("Email outbox worker started")
		for {
			processNextEmail()
			// Randomised sleep so multiple instances don't poll in lock-step.
			time.Sleep(time.Duration(4+rand.Intn(3)) * time.Second)
		}
	}()
}

// processNextEmail claims and sends one pending outbox row if a semaphore
// slot is available.
func processNextEmail() {
	select {
	case emailSemaphore <- struct{}{}:
	default:
		return
	}

	go func() {
		defer func() { <-emailSemaphore }()

		var row models.EmailOutbox
		err := DB.
			Where("status = ? AND next_attempt_at <= ?", models.EmailOutboxPending, time.Now()).
			Order("next_attempt_at ASC").
			Limit(1).
			First(&row).Error
		if err != nil {
			return
		}

		// Atomically claim the row so concurrent workers can't double-send.
		res := DB.Model(&models.EmailOutbox{}).
			Where("id = ? AND status = ?", row.ID, models.EmailOutboxPending).
			Updates(map[string]interface{}{"status": models.EmailOutboxSending})
		if res.Error != nil || res.RowsAffected == 0 {
			return
		}

		sendEmailRow(row)
	}()
}

// sendEmailRow does the actual work: resolve settings, build message,
// send, and update the row based on outcome.
func sendEmailRow(row models.EmailOutbox) {
	logger := appLog.WithFields(map[string]interface{}{
		"outbox_id":     row.ID,
		"license_id":    row.LicenseId,
		"template_code": row.TemplateCode,
		"action":        "EmailOutboxWorker",
	})

	var settings models.EmailSettings
	if err := DB.Where("license_id = ?", row.LicenseId).First(&settings).Error; err != nil {
		markEmailFailed(row, "no SMTP settings configured for license", logger, true)
		return
	}

	password, err := crypto.Decrypt(settings.AuthPasswordEncrypted)
	if err != nil {
		markEmailFailed(row, "decrypt SMTP password: "+err.Error(), logger, true)
		return
	}

	var to, cc, bcc []string
	_ = json.Unmarshal([]byte(row.ToRecipients), &to)
	_ = json.Unmarshal([]byte(row.CcRecipients), &cc)
	_ = json.Unmarshal([]byte(row.BccRecipients), &bcc)

	attachments, err := email.ResolveAttachments(row.Attachments)
	if err != nil {
		markEmailFailed(row, "resolve attachments: "+err.Error(), logger, true)
		return
	}

	fromAddr := row.FromAddress
	if fromAddr == "" {
		fromAddr = settings.FromAddress
	}
	fromName := row.FromName
	if fromName == "" {
		fromName = settings.FromName
	}
	replyTo := row.ReplyTo
	if replyTo == "" {
		replyTo = settings.ReplyTo
	}

	msg := email.Message{
		FromAddress: fromAddr,
		FromName:    fromName,
		ReplyTo:     replyTo,
		To:          to,
		Cc:          cc,
		Bcc:         bcc,
		Subject:     row.Subject,
		Body:        row.Body,
		Attachments: attachments,
	}

	mailer := email.NewSMTPMailer(settings, password)
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	if err := mailer.Send(ctx, msg); err != nil {
		markEmailFailed(row, err.Error(), logger, false)
		return
	}

	now := time.Now()
	if err := DB.Model(&models.EmailOutbox{}).
		Where("id = ?", row.ID).
		Updates(map[string]interface{}{
			"status":     models.EmailOutboxSent,
			"sent_at":    now,
			"attempts":   row.Attempts + 1,
			"last_error": "",
		}).Error; err != nil {
		logger.WithField("error", err.Error()).Error("Failed to mark email as sent")
		return
	}
	logger.Info("Email sent")
}

// markEmailFailed either reschedules the row with exponential backoff or
// marks it permanently failed when max_attempts has been reached. If terminal
// is true the failure is treated as non-retryable regardless of attempt count.
func markEmailFailed(row models.EmailOutbox, errMsg string, logger *logrus.Entry, terminal bool) {
	attempts := row.Attempts + 1
	updates := map[string]interface{}{
		"attempts":   attempts,
		"last_error": errMsg,
	}
	if terminal || attempts >= row.MaxAttempts {
		updates["status"] = models.EmailOutboxFailed
		logger.WithField("error", errMsg).Error("Email permanently failed")
	} else {
		updates["status"] = models.EmailOutboxPending
		updates["next_attempt_at"] = time.Now().Add(emailBackoff(attempts))
		logger.WithField("error", errMsg).Warn("Email send failed, will retry")
	}
	if err := DB.Model(&models.EmailOutbox{}).Where("id = ?", row.ID).Updates(updates).Error; err != nil {
		logger.WithField("error", err.Error()).Error("Failed to update outbox row after send failure")
	}
}

// emailBackoff returns the delay before the next retry for a row at the given
// attempt number. Doubles from 1 minute up to a cap of 16 minutes.
func emailBackoff(attempt int) time.Duration {
	switch {
	case attempt <= 1:
		return 1 * time.Minute
	case attempt == 2:
		return 2 * time.Minute
	case attempt == 3:
		return 4 * time.Minute
	case attempt == 4:
		return 8 * time.Minute
	default:
		return 16 * time.Minute
	}
}

