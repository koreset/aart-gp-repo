package controllers

import (
	"api/models"
	"api/services"
	"bytes"
	"encoding/csv"
	"io"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GenerateMonthlySchedule generates a premium schedule for a scheme.
func GenerateMonthlySchedule(c *gin.Context) {
	var req models.GenerateScheduleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.PremiumResponse{Success: false, Message: err.Error()})
		return
	}

	user := c.MustGet("user").(models.AppUser)

	schedule, err := services.GenerateMonthlySchedule(req.SchemeID, req.Month, req.Year, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.PremiumResponse{Success: false, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.PremiumResponse{Success: true, Data: schedule})
}

// GenerateAllSchedules generates premium schedules for all in-force schemes.
func GenerateAllSchedules(c *gin.Context) {
	var req models.BulkGenerateScheduleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.PremiumResponse{Success: false, Message: err.Error()})
		return
	}

	user := c.MustGet("user").(models.AppUser)
	result, err := services.GenerateAllSchedules(req.Month, req.Year, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.PremiumResponse{Success: false, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.PremiumResponse{Success: true, Data: result})
}

// GenerateInvoice creates an invoice from a schedule.
func GenerateInvoice(c *gin.Context) {
	scheduleID, _ := strconv.Atoi(c.Param("schedule_id"))
	user := c.MustGet("user").(models.AppUser)

	var req models.GenerateInvoiceRequest
	_ = c.ShouldBindJSON(&req) // optional body — ignore bind errors

	invoice, err := services.GenerateInvoice(scheduleID, req.DueDate, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.PremiumResponse{Success: false, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.PremiumResponse{Success: true, Data: invoice})
}

// RecordPayment records a payment.
func RecordPayment(c *gin.Context) {
	var req models.RecordPaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.PremiumResponse{Success: false, Message: err.Error()})
		return
	}

	user := c.MustGet("user").(models.AppUser)
	payment, err := services.RecordPayment(req, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.PremiumResponse{Success: false, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.PremiumResponse{Success: true, Data: payment})
}

// GetArrearsStatus returns the arrears status for a scheme.
func GetArrearsStatus(c *gin.Context) {
	schemeID, _ := strconv.Atoi(c.Param("scheme_id"))

	status, err := services.GetArrearsStatus(schemeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.PremiumResponse{Success: false, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.PremiumResponse{Success: true, Data: status})
}

// GetContributionConfig returns the contribution config for a scheme.
func GetContributionConfig(c *gin.Context) {
	schemeID, _ := strconv.Atoi(c.Param("scheme_id"))

	cfg, err := services.GetContributionConfig(schemeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.PremiumResponse{Success: false, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.PremiumResponse{Success: true, Data: cfg})
}

// SaveContributionConfig upserts the contribution config for a scheme.
func SaveContributionConfig(c *gin.Context) {
	schemeID, _ := strconv.Atoi(c.Param("scheme_id"))

	var cfg models.ContributionConfig
	if err := c.ShouldBindJSON(&cfg); err != nil {
		c.JSON(http.StatusBadRequest, models.PremiumResponse{Success: false, Message: err.Error()})
		return
	}
	cfg.SchemeID = schemeID

	user := c.MustGet("user").(models.AppUser)
	saved, err := services.SaveContributionConfig(cfg, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.PremiumResponse{Success: false, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.PremiumResponse{Success: true, Data: saved})
}

// GetPremiumSchedules returns a filtered list of premium schedules.
func GetPremiumSchedules(c *gin.Context) {
	schemeID, _ := strconv.Atoi(c.Query("scheme_id"))
	month, _ := strconv.Atoi(c.Query("month"))
	year, _ := strconv.Atoi(c.Query("year"))
	status := c.Query("status")

	schedules, err := services.GetPremiumSchedules(schemeID, month, year, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.PremiumResponse{Success: false, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.PremiumResponse{Success: true, Data: schedules})
}

// GetScheduleDetail returns a schedule with member rows.
func GetScheduleDetail(c *gin.Context) {
	scheduleID, _ := strconv.Atoi(c.Param("schedule_id"))

	detail, err := services.GetScheduleDetail(scheduleID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.PremiumResponse{Success: false, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.PremiumResponse{Success: true, Data: detail})
}

// FinalizeSchedule marks a schedule as finalized.
func FinalizeSchedule(c *gin.Context) {
	scheduleID, _ := strconv.Atoi(c.Param("schedule_id"))
	user := c.MustGet("user").(models.AppUser)

	schedule, err := services.FinalizeSchedule(scheduleID, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.PremiumResponse{Success: false, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.PremiumResponse{Success: true, Data: schedule})
}

// ExportScheduleCSV streams a CSV of schedule member rows.
func ExportScheduleCSV(c *gin.Context) {
	scheduleID, _ := strconv.Atoi(c.Param("schedule_id"))

	csvBytes, err := services.ExportScheduleCSV(scheduleID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.PremiumResponse{Success: false, Message: err.Error()})
		return
	}

	c.Header("Content-Disposition", "attachment; filename=schedule_"+strconv.Itoa(scheduleID)+".csv")
	c.Data(http.StatusOK, "text/csv", csvBytes)
}

// GetInvoices returns a filtered list of invoices.
func GetInvoices(c *gin.Context) {
	schemeID, _ := strconv.Atoi(c.Query("scheme_id"))
	month, _ := strconv.Atoi(c.Query("month"))
	year, _ := strconv.Atoi(c.Query("year"))
	status := c.Query("status")

	invoices, err := services.GetInvoices(schemeID, month, year, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.PremiumResponse{Success: false, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.PremiumResponse{Success: true, Data: invoices})
}

// GetInvoiceDetail returns an invoice with line items and payment history.
func GetInvoiceDetail(c *gin.Context) {
	invoiceID, _ := strconv.Atoi(c.Param("invoice_id"))

	detail, err := services.GetInvoiceDetail(invoiceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.PremiumResponse{Success: false, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.PremiumResponse{Success: true, Data: detail})
}

// MarkInvoiceSent marks an invoice as sent.
func MarkInvoiceSent(c *gin.Context) {
	invoiceID, _ := strconv.Atoi(c.Param("invoice_id"))
	user := c.MustGet("user").(models.AppUser)

	invoice, err := services.MarkInvoiceSent(invoiceID, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.PremiumResponse{Success: false, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.PremiumResponse{Success: true, Data: invoice})
}

// CreateAdjustmentNote creates a credit or debit note on an invoice.
func CreateAdjustmentNote(c *gin.Context) {
	invoiceID, _ := strconv.Atoi(c.Param("invoice_id"))

	var req models.AdjustmentNoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.PremiumResponse{Success: false, Message: err.Error()})
		return
	}
	req.InvoiceID = invoiceID

	user := c.MustGet("user").(models.AppUser)
	adj, err := services.CreateAdjustmentNote(req, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.PremiumResponse{Success: false, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.PremiumResponse{Success: true, Data: adj})
}

// GetPayments returns a filtered list of payments.
func GetPayments(c *gin.Context) {
	schemeID, _ := strconv.Atoi(c.Query("scheme_id"))
	method := c.Query("method")
	status := c.Query("status")

	payments, err := services.GetPayments(schemeID, method, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.PremiumResponse{Success: false, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.PremiumResponse{Success: true, Data: payments})
}

// VoidPayment voids a payment record.
func VoidPayment(c *gin.Context) {
	paymentID, _ := strconv.Atoi(c.Param("payment_id"))
	user := c.MustGet("user").(models.AppUser)

	if err := services.VoidPayment(paymentID, user); err != nil {
		c.JSON(http.StatusInternalServerError, models.PremiumResponse{Success: false, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.PremiumResponse{Success: true, Message: "Payment voided"})
}

// BulkImportPayments parses an uploaded CSV and records payments.
func BulkImportPayments(c *gin.Context) {
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, models.PremiumResponse{Success: false, Message: "file is required"})
		return
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.PremiumResponse{Success: false, Message: err.Error()})
		return
	}

	reader := csv.NewReader(bytes.NewReader(data))
	// Skip header row
	if _, err := reader.Read(); err != nil {
		c.JSON(http.StatusBadRequest, models.PremiumResponse{Success: false, Message: "invalid CSV"})
		return
	}
	rows, err := reader.ReadAll()
	if err != nil {
		c.JSON(http.StatusBadRequest, models.PremiumResponse{Success: false, Message: err.Error()})
		return
	}

	user := c.MustGet("user").(models.AppUser)
	result, err := services.BulkImportPayments(rows, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.PremiumResponse{Success: false, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.PremiumResponse{Success: true, Data: result})
}

// AutoMatchPayments auto-matches unmatched payments to invoices.
func AutoMatchPayments(c *gin.Context) {
	user := c.MustGet("user").(models.AppUser)

	result, err := services.AutoMatchPayments(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.PremiumResponse{Success: false, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.PremiumResponse{Success: true, Data: result})
}

// ManualMatchPayment manually links a payment to an invoice.
func ManualMatchPayment(c *gin.Context) {
	var req models.ManualMatchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.PremiumResponse{Success: false, Message: err.Error()})
		return
	}

	user := c.MustGet("user").(models.AppUser)
	payment, err := services.ManualMatchPayment(req, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.PremiumResponse{Success: false, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.PremiumResponse{Success: true, Data: payment})
}

// GetAllArrearsAging returns aging data for all in-force schemes.
func GetAllArrearsAging(c *gin.Context) {
	status := c.Query("status")

	records, err := services.GetAllArrearsAging(status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.PremiumResponse{Success: false, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.PremiumResponse{Success: true, Data: records})
}

// SendArrearsReminder logs a reminder event for a scheme.
func SendArrearsReminder(c *gin.Context) {
	schemeID, _ := strconv.Atoi(c.Param("scheme_id"))

	var req models.SendReminderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.PremiumResponse{Success: false, Message: err.Error()})
		return
	}

	user := c.MustGet("user").(models.AppUser)
	if err := services.SendArrearsReminder(schemeID, req, user); err != nil {
		c.JSON(http.StatusInternalServerError, models.PremiumResponse{Success: false, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.PremiumResponse{Success: true, Message: "Reminder logged"})
}

// RecordPaymentPlan creates a payment plan for a scheme.
func RecordPaymentPlan(c *gin.Context) {
	schemeID, _ := strconv.Atoi(c.Param("scheme_id"))

	var req models.PaymentPlanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.PremiumResponse{Success: false, Message: err.Error()})
		return
	}

	user := c.MustGet("user").(models.AppUser)
	plan, err := services.RecordPaymentPlan(schemeID, req, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.PremiumResponse{Success: false, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.PremiumResponse{Success: true, Data: plan})
}

// SuspendSchemeCover suspends cover for a scheme.
func SuspendSchemeCover(c *gin.Context) {
	schemeID, _ := strconv.Atoi(c.Param("scheme_id"))

	var req models.SuspendCoverRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.PremiumResponse{Success: false, Message: err.Error()})
		return
	}

	user := c.MustGet("user").(models.AppUser)
	if err := services.SuspendSchemeCover(schemeID, req, user); err != nil {
		c.JSON(http.StatusInternalServerError, models.PremiumResponse{Success: false, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.PremiumResponse{Success: true, Message: "Cover suspended"})
}

// ReinstateSchemeCover reinstates cover for a scheme.
func ReinstateSchemeCover(c *gin.Context) {
	schemeID, _ := strconv.Atoi(c.Param("scheme_id"))

	var req models.ReinstateCoverRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.PremiumResponse{Success: false, Message: err.Error()})
		return
	}

	user := c.MustGet("user").(models.AppUser)
	if err := services.ReinstateSchemeCover(schemeID, req, user); err != nil {
		c.JSON(http.StatusInternalServerError, models.PremiumResponse{Success: false, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.PremiumResponse{Success: true, Message: "Cover reinstated"})
}

// GetArrearsHistory returns the arrears event history for a scheme.
func GetArrearsHistory(c *gin.Context) {
	schemeID, _ := strconv.Atoi(c.Param("scheme_id"))

	history, err := services.GetArrearsHistory(schemeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.PremiumResponse{Success: false, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.PremiumResponse{Success: true, Data: history})
}

// GetEmployerStatement returns a ledger statement for a scheme.
func GetEmployerStatement(c *gin.Context) {
	schemeID, _ := strconv.Atoi(c.Param("scheme_id"))
	from := c.Query("from")
	to := c.Query("to")

	if from == "" || to == "" {
		c.JSON(http.StatusBadRequest, models.PremiumResponse{Success: false, Message: "from and to query params are required"})
		return
	}

	statement, err := services.GetEmployerStatement(schemeID, from, to)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.PremiumResponse{Success: false, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.PremiumResponse{Success: true, Data: statement})
}

// GetBrokerCommissionStatement returns a commission statement for a broker.
func GetBrokerCommissionStatement(c *gin.Context) {
	brokerID, _ := strconv.Atoi(c.Param("broker_id"))
	from := c.Query("from")
	to := c.Query("to")

	if from == "" || to == "" {
		c.JSON(http.StatusBadRequest, models.PremiumResponse{Success: false, Message: "from and to query params are required"})
		return
	}

	statement, err := services.GetBrokerCommissionStatement(brokerID, from, to)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.PremiumResponse{Success: false, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.PremiumResponse{Success: true, Data: statement})
}

// GetPremiumDashboard returns dashboard KPIs and charts data.
func GetPremiumDashboard(c *gin.Context) {
	year, _ := strconv.Atoi(c.Query("year"))

	data, err := services.GetPremiumDashboard(year)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.PremiumResponse{Success: false, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.PremiumResponse{Success: true, Data: data})
}

// GetCollectionRate returns monthly collection rates for a year.
func GetCollectionRate(c *gin.Context) {
	year, _ := strconv.Atoi(c.Query("year"))

	data, err := services.GetCollectionRate(year)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.PremiumResponse{Success: false, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.PremiumResponse{Success: true, Data: data})
}

// ReviewSchedule submits a draft schedule for review.
func ReviewSchedule(c *gin.Context) {
	scheduleID, _ := strconv.Atoi(c.Param("schedule_id"))
	user := c.MustGet("user").(models.AppUser)

	schedule, err := services.ReviewSchedule(scheduleID, user)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.PremiumResponse{Success: false, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.PremiumResponse{Success: true, Data: schedule})
}

// ApproveSchedule approves a reviewed schedule.
func ApproveSchedule(c *gin.Context) {
	scheduleID, _ := strconv.Atoi(c.Param("schedule_id"))
	user := c.MustGet("user").(models.AppUser)

	schedule, err := services.ApproveSchedule(scheduleID, user)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.PremiumResponse{Success: false, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.PremiumResponse{Success: true, Data: schedule})
}

// ReturnScheduleToDraft returns a reviewed or approved schedule back to draft.
func ReturnScheduleToDraft(c *gin.Context) {
	scheduleID, _ := strconv.Atoi(c.Param("schedule_id"))
	user := c.MustGet("user").(models.AppUser)

	schedule, err := services.ReturnScheduleToDraft(scheduleID, user)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.PremiumResponse{Success: false, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.PremiumResponse{Success: true, Data: schedule})
}

// VoidSchedule voids a finalized or invoiced schedule.
func VoidSchedule(c *gin.Context) {
	scheduleID, _ := strconv.Atoi(c.Param("schedule_id"))

	var req models.VoidScheduleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.PremiumResponse{Success: false, Message: err.Error()})
		return
	}

	user := c.MustGet("user").(models.AppUser)
	schedule, err := services.VoidSchedule(scheduleID, req.Reason, user)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.PremiumResponse{Success: false, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.PremiumResponse{Success: true, Data: schedule})
}

// CancelSchedule cancels a draft, reviewed, or approved schedule.
func CancelSchedule(c *gin.Context) {
	scheduleID, _ := strconv.Atoi(c.Param("schedule_id"))

	var req models.VoidScheduleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.PremiumResponse{Success: false, Message: err.Error()})
		return
	}

	user := c.MustGet("user").(models.AppUser)
	schedule, err := services.CancelSchedule(scheduleID, req.Reason, user)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.PremiumResponse{Success: false, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.PremiumResponse{Success: true, Data: schedule})
}

// RegenerateSchedule drops and rebuilds all member rows for a draft schedule.
func RegenerateSchedule(c *gin.Context) {
	scheduleID, _ := strconv.Atoi(c.Param("schedule_id"))
	//month, _ := strconv.Atoi(c.Param("month"))
	//year, _ := strconv.Atoi(c.Param("year"))
	user := c.MustGet("user").(models.AppUser)

	schedule, err := services.RegenerateSchedule(scheduleID, user)
	//schedule, err := services.GenerateMonthlySchedule(scheduleID, month, year, user)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.PremiumResponse{Success: false, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.PremiumResponse{Success: true, Data: schedule})
}

// RemoveScheduleMember removes a member row from a draft schedule.
func RemoveScheduleMember(c *gin.Context) {
	scheduleID, _ := strconv.Atoi(c.Param("schedule_id"))
	rowID, _ := strconv.Atoi(c.Param("row_id"))
	user := c.MustGet("user").(models.AppUser)

	if err := services.RemoveScheduleMemberRow(scheduleID, rowID, user); err != nil {
		c.JSON(http.StatusBadRequest, models.PremiumResponse{Success: false, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.PremiumResponse{Success: true, Message: "Member removed"})
}

// UpdateScheduleMemberRow updates the rate fields on a member row in a draft schedule.
func UpdateScheduleMemberRow(c *gin.Context) {
	scheduleID, _ := strconv.Atoi(c.Param("schedule_id"))
	rowID, _ := strconv.Atoi(c.Param("row_id"))

	var req models.ScheduleMemberRateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.PremiumResponse{Success: false, Message: err.Error()})
		return
	}

	user := c.MustGet("user").(models.AppUser)
	row, err := services.UpdateScheduleMemberRow(scheduleID, rowID, req, user)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.PremiumResponse{Success: false, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.PremiumResponse{Success: true, Data: row})
}

// CreateCreditNote creates a credit note for reconciliation purposes.
func CreateCreditNote(c *gin.Context) {
	var req models.AdjustmentNoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.PremiumResponse{Success: false, Message: err.Error()})
		return
	}
	req.Type = "credit"
	// Find latest invoice for scheme if invoice_id not provided
	if req.InvoiceID == 0 {
		var inv models.Invoice
		if err := services.DB.Where("scheme_id = ? AND status IN ?", req.SchemeID, []string{"sent", "partial"}).
			Order("issue_date DESC").First(&inv).Error; err == nil {
			req.InvoiceID = inv.ID
		}
	}

	user := c.MustGet("user").(models.AppUser)
	adj, err := services.CreateAdjustmentNote(req, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.PremiumResponse{Success: false, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.PremiumResponse{Success: true, Data: adj})
}

// CreateDebitNote creates a debit note for reconciliation purposes.
func CreateDebitNote(c *gin.Context) {
	var req models.AdjustmentNoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.PremiumResponse{Success: false, Message: err.Error()})
		return
	}
	req.Type = "debit"
	if req.InvoiceID == 0 {
		var inv models.Invoice
		if err := services.DB.Where("scheme_id = ? AND status IN ?", req.SchemeID, []string{"sent", "partial"}).
			Order("issue_date DESC").First(&inv).Error; err == nil {
			req.InvoiceID = inv.ID
		}
	}

	user := c.MustGet("user").(models.AppUser)
	adj, err := services.CreateAdjustmentNote(req, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.PremiumResponse{Success: false, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.PremiumResponse{Success: true, Data: adj})
}

// GetScheduleCoverageMatrix returns a matrix of all in-force schemes × last 12 months
// showing whether premium schedules have been generated.
func GetScheduleCoverageMatrix(c *gin.Context) {
	matrix, err := services.GetScheduleCoverageMatrix()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.PremiumResponse{Success: false, Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, models.PremiumResponse{Success: true, Data: matrix})
}
