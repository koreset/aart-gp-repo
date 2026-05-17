package services

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"strings"
	"time"
)

// ──────────────────────────────────────────────
// Payment exceptions queue (Phase 4)
// ──────────────────────────────────────────────
//
// Failed reconciliation results (status = "failed") were previously only
// visible inside one schedule's Reconciliation tab. Phase 4 aggregates them
// across all schedules so an "Exceptions" screen can drive triage workflow.
// Retry routes through the existing RetryFailedPayments.

// PaymentException is the denormalised projection returned to the UI. Joins
// reconciliation result + schedule + ACB file so the row carries everything
// finance needs to triage without follow-up calls.
type PaymentException struct {
	ID              int     `json:"id"`
	ACBFileID       int     `json:"acb_file_id"`
	ScheduleID      int     `json:"schedule_id"`
	ScheduleNumber  string  `json:"schedule_number"`
	ScheduleItemID  int     `json:"schedule_item_id"`
	ClaimID         int     `json:"claim_id"`
	ClaimNumber     string  `json:"claim_number"`
	MemberName      string  `json:"member_name"`
	BenefitName     string  `json:"benefit_name"`
	AccountNumber   string  `json:"account_number"`
	Amount          float64 `json:"amount"`
	Status          string  `json:"status"`
	FailureReason   string  `json:"failure_reason"`
	BankReference   string  `json:"bank_reference"`
	ResponseCode    string  `json:"response_code"`
	CreatedAt       string  `json:"created_at"`
	Resolved        bool    `json:"resolved"`
}

// ListPaymentExceptionsRequest controls filtering of the exceptions list.
// Status="failed" is the default — pass status="unmatched" to see lines the
// bank acknowledged but couldn't be matched back to a schedule item, or
// status="" to see everything.
type ListPaymentExceptionsRequest struct {
	Status         string
	IncludeResolved bool
	Limit          int
}

// ListPaymentExceptions returns the aggregated exception queue. "Resolved" is
// inferred from the source claim's current status — if a failed line has
// since transitioned back to "paid" (via a retry that succeeded), it's
// counted resolved and hidden unless IncludeResolved is true.
func ListPaymentExceptions(req ListPaymentExceptionsRequest) ([]PaymentException, error) {
	status := strings.ToLower(strings.TrimSpace(req.Status))
	if status == "" {
		status = "failed"
	}
	limit := req.Limit
	if limit <= 0 || limit > 500 {
		limit = 100
	}

	q := DB.Table("acb_reconciliation_results r").
		Select(`r.id AS id, r.acb_file_id AS acb_file_id,
			s.id AS schedule_id, s.schedule_number AS schedule_number,
			r.schedule_item_id AS schedule_item_id, r.claim_id AS claim_id,
			r.claim_number AS claim_number, i.member_name AS member_name,
			i.benefit_name AS benefit_name, r.account_number AS account_number,
			r.amount AS amount, r.status AS status,
			r.failure_reason AS failure_reason, r.bank_reference AS bank_reference,
			r.response_code AS response_code, r.created_at AS created_at,
			c.status AS claim_status`).
		Joins("LEFT JOIN claim_payment_schedule_items i ON i.id = r.schedule_item_id").
		Joins("LEFT JOIN claim_payment_schedules s ON s.id = i.schedule_id").
		Joins("LEFT JOIN group_scheme_claims c ON c.id = r.claim_id").
		Where("LOWER(r.status) = ?", status).
		Order("r.created_at DESC").
		Limit(limit)

	type raw struct {
		ID             int
		ACBFileID      int
		ScheduleID     int
		ScheduleNumber string
		ScheduleItemID int
		ClaimID        int
		ClaimNumber    string
		MemberName     string
		BenefitName    string
		AccountNumber  string
		Amount         float64
		Status         string
		FailureReason  string
		BankReference  string
		ResponseCode   string
		CreatedAt      string
		ClaimStatus    string
	}
	var rows []raw
	if err := q.Scan(&rows).Error; err != nil {
		return nil, err
	}

	out := make([]PaymentException, 0, len(rows))
	for _, r := range rows {
		resolved := strings.EqualFold(r.ClaimStatus, "paid")
		if resolved && !req.IncludeResolved {
			continue
		}
		out = append(out, PaymentException{
			ID:             r.ID,
			ACBFileID:      r.ACBFileID,
			ScheduleID:     r.ScheduleID,
			ScheduleNumber: r.ScheduleNumber,
			ScheduleItemID: r.ScheduleItemID,
			ClaimID:        r.ClaimID,
			ClaimNumber:    r.ClaimNumber,
			MemberName:     r.MemberName,
			BenefitName:    r.BenefitName,
			AccountNumber:  r.AccountNumber,
			Amount:         r.Amount,
			Status:         r.Status,
			FailureReason:  r.FailureReason,
			BankReference:  r.BankReference,
			ResponseCode:   r.ResponseCode,
			CreatedAt:      r.CreatedAt,
			Resolved:       resolved,
		})
	}
	return out, nil
}

// ExportPaymentExceptionsCSV writes the same filtered exception rows that
// ListPaymentExceptions would return as a CSV blob. Column set mirrors the
// table shown on the Payment Exceptions screen so the download matches what
// the user sees.
func ExportPaymentExceptionsCSV(req ListPaymentExceptionsRequest) ([]byte, string, error) {
	rows, err := ListPaymentExceptions(req)
	if err != nil {
		return nil, "", err
	}

	var buf bytes.Buffer
	w := csv.NewWriter(&buf)
	_ = w.Write([]string{
		"Claim Number", "Member", "Benefit", "Schedule", "Bank Account",
		"Amount", "Status", "Resolved", "Failure Reason", "Bank Reference",
		"Response Code", "Failed At",
	})
	for _, r := range rows {
		resolved := "no"
		if r.Resolved {
			resolved = "yes"
		}
		_ = w.Write([]string{
			r.ClaimNumber,
			r.MemberName,
			r.BenefitName,
			r.ScheduleNumber,
			r.AccountNumber,
			fmt.Sprintf("%.2f", r.Amount),
			r.Status,
			resolved,
			r.FailureReason,
			r.BankReference,
			r.ResponseCode,
			r.CreatedAt,
		})
	}
	w.Flush()
	if err := w.Error(); err != nil {
		return nil, "", err
	}

	filename := fmt.Sprintf("payment_exceptions_%s.csv", time.Now().Format("20060102_150405"))
	return buf.Bytes(), filename, nil
}

// PaymentExceptionsSummary is the KPI strip on the exceptions screen.
type PaymentExceptionsSummary struct {
	OutstandingFailed   int     `json:"outstanding_failed"`
	OutstandingUnmatched int    `json:"outstanding_unmatched"`
	OutstandingValue    float64 `json:"outstanding_value"`
}

// GetPaymentExceptionsSummary returns the totals strip shown on the
// exceptions queue. Counts only rows whose source claim isn't already paid.
func GetPaymentExceptionsSummary() (PaymentExceptionsSummary, error) {
	var s PaymentExceptionsSummary
	type r struct {
		FailedCount    int64
		UnmatchedCount int64
		OutstandingSum float64
	}
	var row r
	err := DB.Raw(`
		SELECT
		    SUM(CASE WHEN LOWER(r.status) = 'failed' THEN 1 ELSE 0 END) AS failed_count,
		    SUM(CASE WHEN LOWER(r.status) = 'unmatched' THEN 1 ELSE 0 END) AS unmatched_count,
		    COALESCE(SUM(CASE WHEN LOWER(r.status) IN ('failed','unmatched') AND LOWER(COALESCE(c.status,'')) <> 'paid' THEN r.amount ELSE 0 END), 0) AS outstanding_sum
		FROM acb_reconciliation_results r
		LEFT JOIN group_scheme_claims c ON c.id = r.claim_id
		WHERE LOWER(COALESCE(c.status,'')) <> 'paid'
	`).Scan(&row).Error
	if err != nil {
		return s, err
	}
	s.OutstandingFailed = int(row.FailedCount)
	s.OutstandingUnmatched = int(row.UnmatchedCount)
	s.OutstandingValue = row.OutstandingSum
	return s, nil
}
