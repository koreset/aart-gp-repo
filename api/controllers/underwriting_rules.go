package controllers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"api/models"
	"api/services"
)

// ListUWRuleSets returns every rule set, latest first.
func ListUWRuleSets(c *gin.Context) {
	sets, err := services.ListUWRuleSets()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, sets)
}

// GetUWRuleSet returns one rule set with its rules preloaded.
func GetUWRuleSet(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("rule_set_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, "invalid rule_set_id")
		return
	}
	rs, err := services.GetUWRuleSet(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, "rule set not found")
			return
		}
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, rs)
}

// CreateUWRuleSet creates a new (empty) rule set. Body: {name, version}.
func CreateUWRuleSet(c *gin.Context) {
	var payload models.UWRuleSet
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	user := c.MustGet("user").(models.AppUser)
	payload.CreatedBy = user.UserName
	rs, err := services.CreateUWRuleSet(payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusCreated, rs)
}

// ActivateUWRuleSet marks the given rule set active and deactivates every
// other set. Only one rule set is ever active at a time.
func ActivateUWRuleSet(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("rule_set_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, "invalid rule_set_id")
		return
	}
	if _, err := services.ActivateUWRuleSet(id); err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	rs, err := services.GetUWRuleSet(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, rs)
}

// CreateUWRule appends a rule to a rule set.
func CreateUWRule(c *gin.Context) {
	ruleSetID, err := strconv.Atoi(c.Param("rule_set_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, "invalid rule_set_id")
		return
	}
	var payload models.UWRule
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	rule, err := services.CreateUWRule(ruleSetID, payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusCreated, rule)
}

// UpdateUWRule replaces mutable fields on a rule.
func UpdateUWRule(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("rule_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, "invalid rule_id")
		return
	}
	var payload models.UWRule
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	rule, err := services.UpdateUWRule(id, payload)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, "rule not found")
			return
		}
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, rule)
}

// DeleteUWRule removes a rule.
func DeleteUWRule(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("rule_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, "invalid rule_id")
		return
	}
	if err := services.DeleteUWRule(id); err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

// ExportUWRuleSetCSV streams the canonical CSV for the given rule set —
// round-trips with ImportUWRulesCSV so admins can export, edit, re-import.
func ExportUWRuleSetCSV(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("rule_set_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, "invalid rule_set_id")
		return
	}
	body, filename, err := services.ExportUWRuleSetCSV(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, "rule set not found")
			return
		}
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.Header("Content-Disposition", `attachment; filename="`+filename+`"`)
	c.Data(http.StatusOK, "text/csv", body)
}

// DuplicateUWRuleSet creates a new rule set with bumped version and a
// fresh copy of every rule. Returns the new set with rules preloaded.
func DuplicateUWRuleSet(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("rule_set_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, "invalid rule_set_id")
		return
	}
	user := c.MustGet("user").(models.AppUser)
	rs, err := services.DuplicateUWRuleSet(id, user.UserName)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, "rule set not found")
			return
		}
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusCreated, rs)
}

// DeleteUWRuleSet removes a rule set — refused when the set is active
// or any UnderwritingCase references it. Returns 409 with a message the
// renderer surfaces to the admin.
func DeleteUWRuleSet(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("rule_set_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, "invalid rule_set_id")
		return
	}
	if err := services.DeleteUWRuleSet(id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, "rule set not found")
			return
		}
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

// SeedStarterRuleSet creates "Standard underwriting v1" with a fixture
// set of common rules so first-time admins have something to iterate on.
// Idempotent — re-clicks return the existing set.
func SeedStarterRuleSet(c *gin.Context) {
	user := c.MustGet("user").(models.AppUser)
	rs, err := services.SeedStarterRuleSet(user.UserName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, rs)
}

// ImportUWRulesCSV expects multipart/form-data with a single `file` field
// containing the rules CSV. Returns the number of rules imported.
func ImportUWRulesCSV(c *gin.Context) {
	fh, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, "file required")
		return
	}
	f, err := fh.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	defer f.Close()
	user := c.MustGet("user").(models.AppUser)
	count, err := services.ImportUWRulesCSV(f, user.UserName)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"imported": count, "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"imported": count})
}

// DryRunUWRules evaluates the active (or specified) rule set against an
// EvaluationContext supplied in the request body:
//
//	{"context": {"bmi": 32, "smoker": true}, "rule_set_id": 0}
//
// rule_set_id=0 uses the currently-active set.
func DryRunUWRules(c *gin.Context) {
	var payload struct {
		Context    map[string]any `json:"context"`
		RuleSetID  int            `json:"rule_set_id"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	summary, err := services.EvaluateAgainstRuleSet(payload.RuleSetID, payload.Context)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, summary)
}

// DryRunUWRulesForCase replays the engine against the rule set snapshotted
// on the case. The context for the case is built from per-member rating
// fields. When MemberDisclosure (Phase 5) lands, this builds the union of
// disclosure + rating fields. For now, returns whatever rating-derived
// fields we can populate; rules referencing unset fields are skipped.
func DryRunUWRulesForCase(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("case_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, "invalid case_id")
		return
	}
	uw, err := services.GetUnderwritingCase(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, "case not found")
			return
		}
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	ctx := services.BuildRuleContextForCase(*uw, c.Request.URL.Query())
	summary, err := services.EvaluateAgainstRuleSet(uw.RuleSetID, ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, summary)
}
