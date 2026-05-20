package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

	"api/log"
	"api/services/docpdf"
	"github.com/gin-gonic/gin"
)

// GenerateGroupPricingQuotePdf generates a PDF quotation by rendering the same
// DOCX produced by GenerateGroupPricingQuoteDocx and converting it. The
// conversion uses LibreOffice if available, falling back to Microsoft Word
// (via PowerShell COM automation) on Windows hosts. Either path preserves the
// insurer's branded template formatting.
func GenerateGroupPricingQuotePdf(c *gin.Context) {
	ctx := requestContext(c)
	logger := log.WithContext(ctx)

	id := c.Param("id")
	logger.WithField("quote_id", id).Info("Processing GenerateGroupPricingQuotePdf request")

	docxFilename, docxData, err := buildQuoteDocxBytes(ctx, id)
	if err != nil {
		logger.WithField("error", err.Error()).Error("Failed to build DOCX for PDF conversion")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	pdfData, err := docpdf.ConvertDocxToPdf(docxData)
	if err != nil {
		logger.WithField("error", err.Error()).Error("Failed to convert DOCX to PDF")
		if errors.Is(err, docpdf.ErrConverterNotFound) {
			c.JSON(http.StatusNotImplemented, gin.H{
				"error": "PDF conversion is unavailable on this server. Install LibreOffice (https://www.libreoffice.org) or Microsoft Word so the .docx can be rendered to PDF.",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	pdfFilename := strings.TrimSuffix(docxFilename, filepath.Ext(docxFilename)) + ".pdf"

	logger.WithFields(map[string]interface{}{
		"quote_id": id,
		"filename": pdfFilename,
		"size":     len(pdfData),
	}).Info("Successfully generated PDF quotation")

	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", pdfFilename))
	c.Data(http.StatusOK, "application/pdf", pdfData)
}
