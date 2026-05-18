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

// SubmitMemberDisclosure persists a disclosure on a case and triggers a
// rules-engine evaluation. Body: services.DisclosureSubmission. A
// matching ConsentRecord of type `medical_info` must exist on the case
// (record one via POST .../consent first) — POPIA enforcement.
func SubmitMemberDisclosure(c *gin.Context) {
	caseID, err := strconv.Atoi(c.Param("case_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, "invalid case_id")
		return
	}
	var payload services.DisclosureSubmission
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	user := c.MustGet("user").(models.AppUser)
	disc, err := services.SubmitDisclosure(caseID, payload, user)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, "case not found")
			return
		}
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusCreated, disc)
}

// GetMemberDisclosure returns the most recent disclosure on a case.
func GetMemberDisclosure(c *gin.Context) {
	caseID, err := strconv.Atoi(c.Param("case_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, "invalid case_id")
		return
	}
	disc, err := services.GetLatestDisclosure(caseID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	if disc == nil {
		c.JSON(http.StatusNotFound, "no disclosure on case")
		return
	}
	c.JSON(http.StatusOK, disc)
}

// ListMemberDisclosures returns all disclosure rows on a case, newest first.
func ListMemberDisclosures(c *gin.Context) {
	caseID, err := strconv.Atoi(c.Param("case_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, "invalid case_id")
		return
	}
	rows, err := services.ListDisclosures(caseID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, rows)
}

// SubmitActivelyAtWork persists an actively-at-work attestation. The IP
// and User-Agent are captured server-side from the request so the
// signature hash binds them.
func SubmitActivelyAtWork(c *gin.Context) {
	var payload services.AttestationSubmission
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	user := c.MustGet("user").(models.AppUser)
	att, err := services.SubmitAttestation(payload, user, c.ClientIP(), c.Request.UserAgent())
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusCreated, att)
}

// ListActivelyAtWorkAttestations returns attestations for a case or quote.
// Query: ?case_id=... or ?quote_id=...
func ListActivelyAtWorkAttestations(c *gin.Context) {
	caseID, _ := strconv.Atoi(c.Query("case_id"))
	quoteID, _ := strconv.Atoi(c.Query("quote_id"))
	rows, err := services.ListAttestations(caseID, quoteID)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, rows)
}

// SubmitConsent records a POPIA consent. Required before submitting a
// disclosure of type `medical_info`.
func SubmitConsent(c *gin.Context) {
	var payload services.ConsentSubmission
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	user := c.MustGet("user").(models.AppUser)
	rec, err := services.SubmitConsent(payload, user, c.ClientIP())
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusCreated, rec)
}

// ListCaseConsents returns every consent recorded on a case.
func ListCaseConsents(c *gin.Context) {
	caseID, err := strconv.Atoi(c.Param("case_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, "invalid case_id")
		return
	}
	rows, err := services.ListConsents(caseID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, rows)
}
