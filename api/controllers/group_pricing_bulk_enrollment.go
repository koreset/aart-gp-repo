package controllers

import (
	"api/models"
	"api/services"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreateBulkEnrollmentBatch ingests a CSV upload as a single batch.
// Body: { members: [GPricingMemberDataInForce, ...], skip_duplicates: bool,
//        file_name, file_size_bytes, file_checksum }
//
// Responses:
//   - 201 with the new batch + the per-row validation report on success.
//   - 400 on bad input (missing members array, scheme without an in-force quote, etc.).
//   - 500 for unexpected errors.
func CreateBulkEnrollmentBatch(c *gin.Context) {
	schemeID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid scheme id"})
		return
	}

	raw, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	raw = normalizeSlashDates(raw)

	var payload struct {
		Members        []models.GPricingMemberDataInForce `json:"members"`
		SkipDuplicates bool                               `json:"skip_duplicates"`
		FileName       string                             `json:"file_name"`
		FileSizeBytes  int64                              `json:"file_size_bytes"`
		FileChecksum   string                             `json:"file_checksum"`
	}
	if err := json.Unmarshal(raw, &payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if len(payload.Members) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "members array is required and cannot be empty"})
		return
	}

	user := c.MustGet("user").(models.AppUser)

	batch, report, err := services.CreateBulkEnrollmentBatch(
		schemeID,
		payload.Members,
		payload.SkipDuplicates,
		services.BulkEnrollmentBatchFileMeta{
			FileName:      payload.FileName,
			FileSizeBytes: payload.FileSizeBytes,
			FileChecksum:  payload.FileChecksum,
		},
		user,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"batch":             batch,
		"validation_report": report,
	})
}

// ListBulkEnrollmentBatches returns the batches for a scheme, optionally filtered
// by status (?status=pending_approval). The frontend Pending Approvals tab calls
// this with status=pending_approval to drive its grid and badge count.
func ListBulkEnrollmentBatches(c *gin.Context) {
	schemeID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid scheme id"})
		return
	}

	status := c.Query("status")
	batches, err := services.ListBulkEnrollmentBatches(schemeID, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"batches": batches})
}

// GetBulkEnrollmentBatch returns one batch with all member rows linked to it
// and the parsed validation report.
func GetBulkEnrollmentBatch(c *gin.Context) {
	batchID, err := strconv.Atoi(c.Param("batchId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid batch id"})
		return
	}

	batch, members, report, err := services.GetBulkEnrollmentBatch(batchID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"batch":             batch,
		"members":           members,
		"validation_report": report,
	})
}

// RunExternalIDCheckOnBatch triggers the external CheckID bulk validation for
// the RSA IDs in a batch. Findings merge into the per-row ValidationStatus
// and the batch's ValidationReport; counts are recalculated.
func RunExternalIDCheckOnBatch(c *gin.Context) {
	batchID, err := strconv.Atoi(c.Param("batchId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid batch id"})
		return
	}

	user := c.MustGet("user").(models.AppUser)
	batch, report, err := services.RunExternalRSAIDCheckOnBatch(batchID, user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"batch":             batch,
		"validation_report": report,
	})
}

// ApproveBulkEnrollmentBatch flips all draft member rows in the batch to
// status='Active' and marks the batch as approved with the actor + timestamp.
// Requires BlockingCount == 0.
func ApproveBulkEnrollmentBatch(c *gin.Context) {
	batchID, err := strconv.Atoi(c.Param("batchId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid batch id"})
		return
	}

	user := c.MustGet("user").(models.AppUser)
	batch, err := services.ApproveBulkEnrollmentBatch(batchID, user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"batch": batch})
}

// RejectBulkEnrollmentBatch hard-deletes the batch's draft member rows and
// marks the batch as rejected, recording the reviewer-supplied reason for audit.
//
// Body: { reason: string }
func RejectBulkEnrollmentBatch(c *gin.Context) {
	batchID, err := strconv.Atoi(c.Param("batchId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid batch id"})
		return
	}

	var payload struct {
		Reason string `json:"reason"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := c.MustGet("user").(models.AppUser)
	batch, err := services.RejectBulkEnrollmentBatch(batchID, payload.Reason, user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"batch": batch})
}
