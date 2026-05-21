package controllers

import (
	"api/models"
	"api/services"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// GetCpiIndices returns all CPI rows ordered most-recent first.
func GetCpiIndices(c *gin.Context) {
	rows, err := services.GetCpiIndices()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, rows)
}

// UpsertCpiIndex writes (or replaces) a single CPI row.
func UpsertCpiIndex(c *gin.Context) {
	var payload models.CpiIndex
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user := c.MustGet("user").(models.AppUser)
	payload.CreatedBy = user.UserName
	stored, err := services.UpsertCpiIndex(payload)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, stored)
}

// UploadCpiIndices accepts a multipart .csv or .xlsx file, parses it, and
// upserts rows on (year_index, month_index).
func UploadCpiIndices(c *gin.Context) {
	fh, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
		return
	}
	// Reject files larger than 5MB — CPI uploads are tiny by nature.
	const maxBytes = int64(5 << 20)
	if fh.Size > maxBytes {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file too large (max 5MB)"})
		return
	}

	src, err := fh.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer src.Close()

	tmpDir := filepath.Join("tmp", "uploads", "cpi_index")
	if err := os.MkdirAll(tmpDir, 0o755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	safeName := strconv.FormatInt(fh.Size, 10) + "_" + filepath.Base(fh.Filename)
	tmpPath := filepath.Join(tmpDir, safeName)
	dst, err := os.Create(tmpPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if _, err := dst.ReadFrom(src); err != nil {
		dst.Close()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	dst.Close()
	defer os.Remove(tmpPath)

	rows, err := services.ParseCpiUploadFile(tmpPath, fh.Filename)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if len(rows) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no rows found in file"})
		return
	}
	user := c.MustGet("user").(models.AppUser)
	for i := range rows {
		rows[i].CreatedBy = user.UserName
	}
	n, err := services.BulkUpsertCpiIndices(rows)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"upserted": n,
		"file":     strings.TrimSpace(fh.Filename),
	})
}
