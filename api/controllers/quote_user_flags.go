package controllers

import (
	"api/log"
	"api/models"
	"api/services"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetUserFlags is GET /group-pricing/dashboard/user-flags.
//
// Open by default; pass status=resolved or status=all to override.
// Callers without quote:manage_user_flags get the same list but with
// the internal note redacted so flag *state* is visible to all
// dashboard viewers while observations stay confidential.
func GetUserFlags(c *gin.Context) {
	ctx := dashboardContext(c)
	logger := log.WithContext(ctx)

	var filter models.UserFlagsFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	flags, err := services.ListUserFlags(filter)
	if err != nil {
		logger.WithField("error", err.Error()).Error("Failed to list user flags")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !services.UserHasPermission(c, "quote:manage_user_flags") {
		flags = services.StripFlagNotes(flags)
	}

	c.JSON(http.StatusOK, flags)
}

// PostUserFlag is POST /group-pricing/dashboard/user-flags.
// Opens a new flag against a user. The service layer enforces
// reason validity, self-flag prevention, note length, and uniqueness
// of open flags per (user, reason).
func PostUserFlag(c *gin.Context) {
	ctx := dashboardContext(c)
	logger := log.WithContext(ctx)

	var req models.OpenUserFlagRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	actor := c.MustGet("user").(models.AppUser)
	flag, err := services.OpenUserFlag(ctx, req, actor)
	if err != nil {
		status := http.StatusBadRequest
		switch {
		case errors.Is(err, services.ErrUserFlagAlreadyOpen):
			status = http.StatusConflict
		case errors.Is(err, services.ErrUserFlagSelfFlag):
			status = http.StatusForbidden
		}
		logger.WithField("error", err.Error()).
			WithField("target", req.UserName).
			Warn("Failed to open user flag")
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, flag)
}

// PostResolveUserFlag is POST /group-pricing/dashboard/user-flags/:id/resolve.
// Closes an existing open flag with a resolution note.
func PostResolveUserFlag(c *gin.Context) {
	ctx := dashboardContext(c)
	logger := log.WithContext(ctx)

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req models.ResolveUserFlagRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	actor := c.MustGet("user").(models.AppUser)
	flag, err := services.ResolveUserFlag(id, req.ResolutionNote, actor)
	if err != nil {
		status := http.StatusBadRequest
		switch {
		case errors.Is(err, services.ErrUserFlagNotFound):
			status = http.StatusNotFound
		case errors.Is(err, services.ErrUserFlagAlreadyResolved):
			status = http.StatusConflict
		}
		logger.WithField("error", err.Error()).
			WithField("flag_id", id).
			Warn("Failed to resolve user flag")
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, flag)
}
