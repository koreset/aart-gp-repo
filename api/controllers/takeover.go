package controllers

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"api/models"
	"api/services"
)

// UploadPriorInsurerSchedule accepts multipart/form-data with one `file`
// field (CSV) plus form values:
//
//	quote_id, insurer_name, certificate_number, effective_date, expiry_date, notes
//
// On success returns the persisted PriorInsurerSchedule with the count of
// rows imported. The matcher runs automatically afterwards so the caller
// gets an immediately-usable preview via GetPriorInsurerSchedule.
func UploadPriorInsurerSchedule(c *gin.Context) {
	form, formErr := c.MultipartForm()
	if formErr != nil || form == nil {
		c.JSON(http.StatusBadRequest, "multipart/form-data with file required")
		return
	}
	files := form.File["file"]
	if len(files) == 0 {
		c.JSON(http.StatusBadRequest, "file field required")
		return
	}
	quoteID, err := strconv.Atoi(firstFormValue(form.Value, "quote_id"))
	if err != nil || quoteID <= 0 {
		c.JSON(http.StatusBadRequest, "invalid quote_id")
		return
	}
	header := services.PriorInsurerScheduleHeader{
		InsurerName:       firstFormValue(form.Value, "insurer_name"),
		CertificateNumber: firstFormValue(form.Value, "certificate_number"),
		EffectiveDate:     parseDateForm(firstFormValue(form.Value, "effective_date")),
		ExpiryDate:        parseDateForm(firstFormValue(form.Value, "expiry_date")),
		Notes:             firstFormValue(form.Value, "notes"),
	}
	fh := files[0]
	f, err := fh.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	defer f.Close()
	user := c.MustGet("user").(models.AppUser)
	schedule, err := services.ImportPriorInsurerScheduleCSV(quoteID, header, f, fh.Filename, user)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	// Run the matcher immediately so the broker's first preview reflects
	// the actual member-by-member outcome, not just the import stats.
	summary, matchErr := services.MatchPriorMembersToCensus(schedule.ID)
	if matchErr != nil {
		// Schedule is persisted; matcher failure is non-fatal — return
		// the schedule with a warning so the broker can retry the match.
		c.JSON(http.StatusCreated, gin.H{
			"schedule":      schedule,
			"match_summary": nil,
			"match_warning": matchErr.Error(),
		})
		return
	}
	// Re-classify via the rules engine if any takeover rules exist. Empty
	// rule sets are a no-op (the default outcome stays).
	_ = services.ClassifyPriorMembersAgainstRules(schedule.ID)
	c.JSON(http.StatusCreated, gin.H{
		"schedule":      schedule,
		"match_summary": summary,
	})
}

// GetPriorInsurerScheduleForQuote returns the latest schedule (with
// members) for a quote, or 404 if none exists.
func GetPriorInsurerScheduleForQuote(c *gin.Context) {
	quoteID, err := strconv.Atoi(c.Param("quote_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, "invalid quote_id")
		return
	}
	schedule, err := services.GetPriorInsurerScheduleForQuote(quoteID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	if schedule == nil {
		c.JSON(http.StatusNotFound, "no prior insurer schedule for quote")
		return
	}
	c.JSON(http.StatusOK, schedule)
}

// RematchPriorInsurerSchedule re-runs the matcher on an existing schedule
// (useful after the census is edited or new cases are created).
func RematchPriorInsurerSchedule(c *gin.Context) {
	scheduleID, err := strconv.Atoi(c.Param("schedule_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, "invalid schedule_id")
		return
	}
	summary, err := services.MatchPriorMembersToCensus(scheduleID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, "schedule not found")
			return
		}
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	_ = services.ClassifyPriorMembersAgainstRules(scheduleID)
	c.JSON(http.StatusOK, summary)
}

// ApplyTakeoverTermsToCases pushes the prior-member outcomes onto the
// matched underwriting cases as an engine snapshot. Suggest-don't-decide
// — the underwriter still commits the human decision via the existing
// form.
func ApplyTakeoverTermsToCases(c *gin.Context) {
	scheduleID, err := strconv.Atoi(c.Param("schedule_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, "invalid schedule_id")
		return
	}
	user := c.MustGet("user").(models.AppUser)
	touched, err := services.ApplyTakeoverTermsToCases(scheduleID, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"touched": touched})
}

// ListPolicyHandoffSnapshotsForScheme returns every snapshot recorded for
// a scheme, newest first.
func ListPolicyHandoffSnapshotsForScheme(c *gin.Context) {
	schemeID, err := strconv.Atoi(c.Param("scheme_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, "invalid scheme_id")
		return
	}
	rows, err := services.ListPolicyHandoffSnapshotsForScheme(schemeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, rows)
}

// GetPolicyHandoffSnapshot returns one snapshot row by ID. The full
// canonical JSON payload is on the Payload field of the response.
func GetPolicyHandoffSnapshot(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("snapshot_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, "invalid snapshot_id")
		return
	}
	row, err := services.GetPolicyHandoffSnapshot(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, "snapshot not found")
			return
		}
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, row)
}

// DownloadPolicyHandoffSnapshot streams the canonical JSON payload as a
// downloadable file. Compliance auditors and the legal team consume this
// rather than the raw API response.
func DownloadPolicyHandoffSnapshot(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("snapshot_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, "invalid snapshot_id")
		return
	}
	row, err := services.GetPolicyHandoffSnapshot(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, "snapshot not found")
			return
		}
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	filename := "handoff-snapshot-" + strconv.Itoa(row.ID) + ".json"
	c.Header("Content-Disposition", `attachment; filename="`+filename+`"`)
	c.Data(http.StatusOK, "application/json", []byte(row.Payload))
}

// RebuildPolicyHandoffSnapshot rebuilds the snapshot for a (quote,
// scheme) pair. Idempotent — the underlying service refuses to overwrite
// an existing snapshot; manually delete the row first if a re-build is
// truly required.
func RebuildPolicyHandoffSnapshot(c *gin.Context) {
	quoteID, err := strconv.Atoi(c.Param("quote_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, "invalid quote_id")
		return
	}
	schemeID, _ := strconv.Atoi(c.Query("scheme_id"))
	user := c.MustGet("user").(models.AppUser)
	row, err := services.RecordPolicyHandoffSnapshot(nil, quoteID, schemeID, user.UserEmail, "Manual rebuild")
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, row)
}

// ───── Helpers ─────────────────────────────────────────────────────────────

func firstFormValue(values map[string][]string, key string) string {
	if vs, ok := values[key]; ok && len(vs) > 0 {
		return vs[0]
	}
	return ""
}

func parseDateForm(s string) *time.Time {
	if s == "" {
		return nil
	}
	for _, layout := range []string{"2006-01-02", "2006/01/02", time.RFC3339} {
		if t, err := time.Parse(layout, s); err == nil {
			return &t
		}
	}
	return nil
}
