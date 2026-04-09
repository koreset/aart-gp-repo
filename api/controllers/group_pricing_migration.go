package controllers

import (
	"api/models"
	"api/services"
	"mime/multipart"
	"net/http"

	"github.com/gin-gonic/gin"
)

// extractMigrationFiles extracts the CSV file headers from a multipart form.
// scheme_setup, categories, and member_data are required; beneficiaries and claims_experience are optional.
func extractMigrationFiles(c *gin.Context) (map[string]*multipart.FileHeader, error) {
	form, err := c.MultipartForm()
	if err != nil {
		return nil, err
	}

	files := make(map[string]*multipart.FileHeader)
	requiredKeys := []string{"scheme_setup", "categories", "member_data"}
	for _, key := range requiredKeys {
		fh := form.File[key]
		if len(fh) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": key + " CSV file is required"})
			return nil, nil
		}
		files[key] = fh[0]
	}

	optionalKeys := []string{"beneficiaries", "claims_experience"}
	for _, key := range optionalKeys {
		if fh := form.File[key]; len(fh) > 0 {
			files[key] = fh[0]
		}
	}

	return files, nil
}

func ValidateMigration(c *gin.Context) {
	files, err := extractMigrationFiles(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse form: " + err.Error()})
		return
	}
	if files == nil {
		return // error already sent by extractMigrationFiles
	}

	user := c.MustGet("user").(models.AppUser)

	result, err := services.ValidateMigration(files, user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

func ExecuteMigration(c *gin.Context) {
	files, err := extractMigrationFiles(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse form: " + err.Error()})
		return
	}
	if files == nil {
		return
	}

	user := c.MustGet("user").(models.AppUser)

	result, err := services.ExecuteMigration(files, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

func DownloadMigrationTemplate(c *gin.Context) {
	templateName := c.Param("template_name")

	data, filename, err := services.GenerateMigrationTemplate(templateName)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Data(http.StatusOK, "text/csv", data)
}
