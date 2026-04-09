package controllers

import (
	"api/models"
	"api/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreateEmployerSubmission handles POST /bordereaux/submissions
func CreateEmployerSubmission(c *gin.Context) {
	var req models.CreateSubmissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	user := c.MustGet("user").(models.AppUser)
	sub, err := services.CreateEmployerSubmission(req, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": sub})
}

// GetEmployerSubmissions handles GET /bordereaux/submissions
// Query params: scheme_id, month, year, status
func GetEmployerSubmissions(c *gin.Context) {
	schemeID, _ := strconv.Atoi(c.Query("scheme_id"))
	month, _ := strconv.Atoi(c.Query("month"))
	year, _ := strconv.Atoi(c.Query("year"))
	status := c.Query("status")

	subs, err := services.GetEmployerSubmissions(schemeID, month, year, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": subs})
}

// GetEmployerSubmission handles GET /bordereaux/submissions/:id
func GetEmployerSubmission(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	sub, err := services.GetEmployerSubmission(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": sub})
}

// UploadEmployerSubmission handles POST /bordereaux/submissions/:id/upload (multipart: file)
func UploadEmployerSubmission(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "file is required: " + err.Error()})
		return
	}
	file, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "failed to open file: " + err.Error()})
		return
	}
	defer file.Close()

	user := c.MustGet("user").(models.AppUser)
	sub, err := services.UploadEmployerSubmission(id, file, fileHeader, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": sub})
}

// ReviewEmployerSubmission handles POST /bordereaux/submissions/:id/review
func ReviewEmployerSubmission(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var body struct {
		Notes string `json:"notes"`
	}
	_ = c.ShouldBindJSON(&body)
	user := c.MustGet("user").(models.AppUser)
	sub, err := services.ReviewEmployerSubmission(id, body.Notes, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": sub})
}

// RaiseSubmissionQuery handles POST /bordereaux/submissions/:id/query
func RaiseSubmissionQuery(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var req models.RaiseQueryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	user := c.MustGet("user").(models.AppUser)
	sub, err := services.RaiseSubmissionQuery(id, req.QueryNotes, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": sub})
}

// AcceptEmployerSubmission handles POST /bordereaux/submissions/:id/accept
func AcceptEmployerSubmission(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var body struct {
		Notes string `json:"notes"`
	}
	_ = c.ShouldBindJSON(&body)
	user := c.MustGet("user").(models.AppUser)
	sub, err := services.AcceptEmployerSubmission(id, body.Notes, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": sub})
}

// RejectEmployerSubmission handles POST /bordereaux/submissions/:id/reject
func RejectEmployerSubmission(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var req models.RejectSubmissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	user := c.MustGet("user").(models.AppUser)
	sub, err := services.RejectEmployerSubmission(id, req.Reason, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": sub})
}

// GetSubmissionRecords handles GET /bordereaux/submissions/:id/records
func GetSubmissionRecords(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	records, err := services.GetSubmissionRecords(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": records})
}

// GenerateScheduleFromSubmission handles POST /bordereaux/submissions/:id/generate-schedule
func GenerateScheduleFromSubmission(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	user := c.MustGet("user").(models.AppUser)
	schedule, err := services.GenerateScheduleFromSubmission(id, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": schedule})
}

// ComputeSubmissionDelta handles POST /bordereaux/submissions/:id/compute-delta
func ComputeSubmissionDelta(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	user := c.MustGet("user").(models.AppUser)
	summary, err := services.ComputeSubmissionDelta(id, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": summary})
}

// GetSubmissionDeltaRecords handles GET /bordereaux/submissions/:id/delta
func GetSubmissionDeltaRecords(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	records, err := services.GetSubmissionDeltaRecords(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": records})
}

// SyncSubmissionToMemberRegister handles POST /bordereaux/submissions/:id/sync-members
// Applies delta changes (new/ceased/amendment) to the live member register.
func SyncSubmissionToMemberRegister(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	user := c.MustGet("user").(models.AppUser)
	result, err := services.SyncSubmissionToMemberRegister(id, user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": result})
}

// ComputeRegisterDiff handles GET /bordereaux/submissions/:id/register-diff
// Compares valid submission records against the live member register.
func ComputeRegisterDiff(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	result, err := services.ComputeRegisterDiff(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": result})
}

// SnapshotRegisterDiff handles POST /bordereaux/submissions/:id/snapshot-diff
// Forces a fresh live computation and overwrites any existing snapshot.
// Useful for backfilling older submissions or correcting a snapshot.
func SnapshotRegisterDiff(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	user := c.MustGet("user").(models.AppUser)
	result, err := services.RefreshRegisterDiffSnapshot(id, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": result})
}

// ApplySubmissionExits handles POST /bordereaux/submissions/:id/apply-exits
func ApplySubmissionExits(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	user := c.MustGet("user").(models.AppUser)
	result, err := services.ApplySubmissionExits(id, user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": result})
}

// ApplySubmissionAmendments handles POST /bordereaux/submissions/:id/apply-amendments
func ApplySubmissionAmendments(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	user := c.MustGet("user").(models.AppUser)
	result, err := services.ApplySubmissionAmendments(id, user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": result})
}

// GetNewJoinerDetails handles GET /bordereaux/submissions/:id/new-joiner-details
func GetNewJoinerDetails(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	details, err := services.GetNewJoinerDetails(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": details})
}

// UploadNewJoinerDetails handles POST /bordereaux/submissions/:id/upload-new-joiner-details (multipart: file)
func UploadNewJoinerDetails(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "file is required: " + err.Error()})
		return
	}
	file, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "failed to open file: " + err.Error()})
		return
	}
	defer file.Close()

	user := c.MustGet("user").(models.AppUser)
	details, err := services.UploadNewJoinerDetails(id, file, fileHeader, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": details})
}

// SyncNewJoiners handles POST /bordereaux/submissions/:id/sync-new-joiners
func SyncNewJoiners(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	user := c.MustGet("user").(models.AppUser)
	result, err := services.SyncNewJoiners(id, user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": result})
}
