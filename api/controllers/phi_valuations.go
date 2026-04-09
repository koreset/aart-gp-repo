package controllers

import (
	"api/models"
	"api/services"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func GetPhiValuationTableMetadata(c *gin.Context) {
	metadata := services.GetPhiValuationTableMetaData()
	c.JSON(http.StatusOK, metadata)
}

func GetPhiValuationTableData(c *gin.Context) {
	tableType := c.Param("table_type")

	data, err := services.GetPhiValuationTableData(tableType)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, data)
}

func GetPhiValuationTableYears(c *gin.Context) {
	tableType := c.Param("table_type")

	years, err := services.GetPhiValuationTableYears(tableType)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, years)
}

func GetPhiValuationTableVersions(c *gin.Context) {
	tableType := c.Param("table_type")
	year, err := strconv.Atoi(c.Param("year"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid year parameter"})
		return
	}

	versions, err := services.GetPhiValuationTableVersions(tableType, year)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, versions)
}

func UploadPhiValuationTables(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	user := c.MustGet("user").(models.AppUser)

	tableType := form.Value["table_type"][0]
	year, _ := strconv.Atoi(form.Value["year"][0])
	version := form.Value["version"][0]
	file := form.File["file"][0]

	err, count := services.SavePhiTable(file, tableType, year, version, user)

	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"count": count})
}

func CreatePhiShockSetting(c *gin.Context) {
	var shockSetting models.PhiShockSetting
	err := c.BindJSON(&shockSetting)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	shockSetting, err = services.SavePhiShockSetting(shockSetting)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, shockSetting)
}

func GetPhiShockSettings(c *gin.Context) {
	//year, _ := strconv.Atoi(c.Query("year"))
	//version := c.Query("version")

	shockSettings, err := services.GetPhiShockSettings()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, shockSettings)
}

func GetPhiShockBases(c *gin.Context) {
	results := services.GetPhiShockBases()
	c.JSON(http.StatusOK, results)
}

func GetPhiModelPointYears(c *gin.Context) {
	years, err := services.GetPhiModelPointYears()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, years)
}

func GetPhiModelPointYearVersions(c *gin.Context) {
	year, err := strconv.Atoi(c.Param("year"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid year parameter"})
		return
	}

	versions, err := services.GetPhiModelPointYearVersions(year)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, versions)
}

func GetPhiParameterYears(c *gin.Context) {
	years, err := services.GetPhiParameterYears()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, years)
}

func GetPhiParameterYearVersions(c *gin.Context) {
	year, err := strconv.Atoi(c.Param("year"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid year parameter"})
		return
	}
	versions, err := services.GetPhiParameterYearVersions(year)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, versions)
}

func GetPhiMortalityYears(c *gin.Context) {
	years, err := services.GetPhiMortalityYears()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, years)
}

func GetPhiMortalityYearVersions(c *gin.Context) {
	year, err := strconv.Atoi(c.Param("year"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid year parameter"})
		return
	}
	versions, err := services.GetPhiMortalityYearVersions(year)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, versions)
}

func GetPhiRecoveryRateYears(c *gin.Context) {
	years, err := services.GetPhiRecoveryRateYears()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, years)
}

func GetPhiRecoveryRateYearVersions(c *gin.Context) {
	year, err := strconv.Atoi(c.Param("year"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid year parameter"})
		return
	}
	versions, err := services.GetPhiRecoveryRateYearVersions(year)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, versions)
}

func GetPhiYieldCurveYears(c *gin.Context) {
	years, err := services.GetPhiYieldCurveYears()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, years)
}

func GetPhiYieldCurveYearVersions(c *gin.Context) {
	year, err := strconv.Atoi(c.Param("year"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid year parameter"})
		return
	}
	versions, err := services.GetPhiYieldCurveYearVersions(year)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, versions)
}

func RunPhiProjections(c *gin.Context) {
	reqbody, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Error().Err(err).Send()
	}
	var runJobs models.RunPhiJob

	fmt.Println(string(reqbody))
	err = json.Unmarshal(reqbody, &runJobs)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := c.MustGet("user").(models.AppUser)

	runJobs.UserEmail = user.UserEmail
	runJobs.UserName = user.UserName

	//var runJobs []models.RunPhiJob

	//user := c.MustGet("user").(models.AppUser)

	go func() {
		err = services.RunPhiProjection(runJobs)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}()

	c.JSON(http.StatusOK, nil)
}

func GetAllPhiRunJobs(c *gin.Context) {
	user := c.MustGet("user").(models.AppUser)

	jobs, err := services.GetPhiRunJobs(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, jobs)
}

func GetPhiRunJob(c *gin.Context) {
	jobID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid job ID parameter"})
		return
	}

	user := c.MustGet("user").(models.AppUser)

	job, err := services.GetPhiRunJob(jobID, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, job)
}

func GetPhiRunJobControl(c *gin.Context) {
	runId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid job ID parameter"})
		return
	}

	user := c.MustGet("user").(models.AppUser)

	control, err := services.GetPhiRunJobControl(runId, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, control)
}

func DeletePhiRunJob(c *gin.Context) {
	jobID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid job ID parameter"})
		return
	}

	user := c.MustGet("user").(models.AppUser)

	err = services.DeletePhiRunJob(jobID, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Job deleted successfully"})
}

func GetPhiModelPointCount(c *gin.Context) {
	counts, err := services.GetPhiModelPointCount()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, counts)
}

func GetPhiModelPointsData(c *gin.Context) {
	year, err := strconv.Atoi(c.Param("year"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid year parameter"})
		return
	}
	version := c.Param("version")
	mps, err := services.GetPhiModelPointsForYear(year, version)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, mps)
}

func GetPhiModelPointsExcel(c *gin.Context) {
	year, err := strconv.Atoi(c.Param("year"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid year parameter"})
		return
	}
	version := c.Param("version")
	excelData, err := services.GetPhiModelPointsExcel(year, version)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to generate Excel: %v", err)})
		return
	}
	c.Header("Content-Disposition", `attachment; filename="phi_model_points.xlsx"`)
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", excelData)
}

func DeletePhiModelPoints(c *gin.Context) {
	year, err := strconv.Atoi(c.Param("year"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid year parameter"})
		return
	}
	version := c.Param("version")
	if err := services.DeletePhiModelPoints(year, version); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "PHI model points deleted successfully"})
}

func DeletePhiValuationTableData(c *gin.Context) {
	tableType := c.Param("table_type")
	year, err := strconv.Atoi(c.Param("year"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid year parameter"})
		return
	}
	version := c.Param("version")

	err = services.DeletePhiValuationTableData(tableType, year, version)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Data deleted successfully"})
}

func ListPhiRunConfigs(c *gin.Context) {
	configs, err := services.ListPhiRunConfigs()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, configs)
}

func SavePhiRunConfig(c *gin.Context) {
	user := c.MustGet("user").(models.AppUser)
	var config models.PhiRunConfig
	if err := c.ShouldBindJSON(&config); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	config.CreatedBy = user.UserName
	saved, err := services.SavePhiRunConfig(config)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, saved)
}

func DeletePhiRunConfig(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id parameter"})
		return
	}
	if err := services.DeletePhiRunConfig(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Config deleted"})
}
