package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"api/log"
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

	pdfData, err := convertDocxToPdf(docxData)
	if err != nil {
		logger.WithField("error", err.Error()).Error("Failed to convert DOCX to PDF")
		if errors.Is(err, errPdfConverterNotFound) {
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

var errPdfConverterNotFound = errors.New("no DOCX→PDF converter found (LibreOffice or Microsoft Word required)")

// convertDocxToPdf writes the DOCX bytes to a temp file, runs an available
// converter (LibreOffice headless preferred, Microsoft Word via PowerShell on
// Windows as a fallback), then reads the resulting PDF back. Temp files are
// cleaned up on return.
func convertDocxToPdf(docxData []byte) ([]byte, error) {
	tmpDir, err := os.MkdirTemp("", "quote_pdf_*")
	if err != nil {
		return nil, fmt.Errorf("create temp dir: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	docxPath := filepath.Join(tmpDir, fmt.Sprintf("quote_%d.docx", time.Now().UnixNano()))
	if err := os.WriteFile(docxPath, docxData, 0o600); err != nil {
		return nil, fmt.Errorf("write docx: %w", err)
	}
	pdfPath := strings.TrimSuffix(docxPath, filepath.Ext(docxPath)) + ".pdf"

	if bin, ok := findLibreOffice(); ok {
		if err := runLibreOffice(bin, docxPath, tmpDir); err != nil {
			return nil, err
		}
		return os.ReadFile(pdfPath)
	}

	if runtime.GOOS == "windows" && hasMicrosoftWord() {
		if err := runWordConvert(docxPath, pdfPath); err != nil {
			return nil, err
		}
		return os.ReadFile(pdfPath)
	}

	return nil, errPdfConverterNotFound
}

// findLibreOffice returns the path of a usable LibreOffice headless binary if
// available. It checks `soffice` first (the actual binary on all platforms),
// falls back to `libreoffice` (the wrapper present on most Linux distros),
// and finally probes common Windows install locations.
func findLibreOffice() (string, bool) {
	for _, name := range []string{"soffice", "libreoffice"} {
		if path, err := exec.LookPath(name); err == nil {
			return path, true
		}
	}
	for _, candidate := range []string{
		`C:\Program Files\LibreOffice\program\soffice.exe`,
		`C:\Program Files (x86)\LibreOffice\program\soffice.exe`,
	} {
		if _, err := os.Stat(candidate); err == nil {
			return candidate, true
		}
	}
	return "", false
}

func runLibreOffice(bin, docxPath, outDir string) error {
	// -env:UserInstallation isolates the profile so concurrent conversions
	// don't collide on a shared user profile lock.
	profileURI := "file:///" + filepath.ToSlash(filepath.Join(outDir, "lo_profile"))
	cmd := exec.Command(bin,
		"-env:UserInstallation="+profileURI,
		"--headless",
		"--convert-to", "pdf",
		"--outdir", outDir,
		docxPath,
	)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("libreoffice convert failed: %w (%s)", err, strings.TrimSpace(string(out)))
	}
	return nil
}

// hasMicrosoftWord checks the standard Office install locations on Windows.
// Word automation is reached via PowerShell COM rather than direct Word CLI,
// so we just need to confirm Word is installed.
func hasMicrosoftWord() bool {
	for _, candidate := range []string{
		`C:\Program Files\Microsoft Office\root\Office16\WINWORD.EXE`,
		`C:\Program Files (x86)\Microsoft Office\root\Office16\WINWORD.EXE`,
		`C:\Program Files\Microsoft Office\Office16\WINWORD.EXE`,
		`C:\Program Files (x86)\Microsoft Office\Office16\WINWORD.EXE`,
	} {
		if _, err := os.Stat(candidate); err == nil {
			return true
		}
	}
	return false
}

// runWordConvert opens the DOCX in Microsoft Word via PowerShell COM
// automation and saves it as PDF (wdFormatPDF = 17). PowerShell is preferred
// over a direct Go COM binding because it ships with Windows and avoids extra
// dependencies. The script is invoked with -ExecutionPolicy Bypass so it
// works even on hosts where the user's policy is Restricted.
func runWordConvert(docxPath, pdfPath string) error {
	script := fmt.Sprintf(`
$ErrorActionPreference = 'Stop'
$word = New-Object -ComObject Word.Application
$word.Visible = $false
$word.DisplayAlerts = 0
try {
    $doc = $word.Documents.Open(%q, $false, $true)
    $doc.SaveAs([ref]%q, [ref]17)
    $doc.Close($false)
} finally {
    $word.Quit()
    [System.Runtime.InteropServices.Marshal]::ReleaseComObject($word) | Out-Null
}
`, docxPath, pdfPath)

	cmd := exec.Command("powershell.exe",
		"-NoProfile",
		"-NonInteractive",
		"-ExecutionPolicy", "Bypass",
		"-Command", script,
	)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("word convert failed: %w (%s)", err, strings.TrimSpace(string(out)))
	}
	if _, err := os.Stat(pdfPath); err != nil {
		return fmt.Errorf("word convert produced no PDF (%s)", strings.TrimSpace(string(out)))
	}
	return nil
}
