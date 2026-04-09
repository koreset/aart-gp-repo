package controllers

import (
	"api/models"
	"api/services"
	"api/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func CreateExpConfiguration(c *gin.Context) {
	var configuration models.ExpConfiguration
	if err := c.ShouldBindJSON(&configuration); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := services.CreateConfiguration(configuration); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Configuration created successfully"})

}

func GetExpConfigurations(c *gin.Context) {
	configurations, err := services.GetConfigurations()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, configurations)

}

func GetExpTableData(c *gin.Context) {
	var tablePayload models.TableDataPayload
	if err := c.ShouldBindJSON(&tablePayload); err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	tableData, err := services.GetTableData(tablePayload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return

	}
	c.JSON(http.StatusOK, tableData)

}

func UploadAgeBands(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	file := form.File["file"][0]
	user := c.MustGet("user").(models.AppUser)

	if err := services.UploadAgeBands(file, user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Data uploaded successfully"})

}

func GetAgeBandVersions(c *gin.Context) {
	versions, err := services.GetAgeBandVersions()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, versions)
}

func UploadExpExposureData(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	file := form.File["file"][0]
	year := form.Value["year"][0]
	version := form.Value["version"][0]
	expId := form.Value["exp_id"][0]
	user := c.MustGet("user").(models.AppUser)

	if err := services.UploadExpExposureData(file, year, version, expId, user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Data uploaded successfully"})
}

func UploadExpActualData(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	file := form.File["file"][0]
	year := form.Value["year"][0]
	version := form.Value["version"][0]
	expId := form.Value["exp_id"][0]
	user := c.MustGet("user").(models.AppUser)

	if err := services.UploadExpActualData(file, year, version, expId, user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Data uploaded successfully"})
}

func UploadExpCurrentMortality(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	file := form.File["file"][0]
	year := form.Value["year"][0]
	version := form.Value["version"][0]
	expId := form.Value["exp_id"][0]
	user := c.MustGet("user").(models.AppUser)

	if err := services.UploadExpCurrentMortality(file, year, version, expId, user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Data uploaded successfully"})
}

func UploadExpCurrentLapse(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	file := form.File["file"][0]
	year := form.Value["year"][0]
	version := form.Value["version"][0]
	expId := form.Value["exp_id"][0]
	user := c.MustGet("user").(models.AppUser)

	if err := services.UploadExpCurrentLapse(file, year, version, expId, user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Data uploaded successfully"})
}

func GetExpExposureActualYears(c *gin.Context) {
	expId := c.Param("exp_id")
	years, err := services.GetExposureActualYears(expId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, years)
}

func GetExpExposureDataVersions(c *gin.Context) {
	expId := c.Param("exp_id")
	year := c.Param("year")
	versions, err := services.GetExposureDataVersions(expId, year)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, versions)
}

func GetExpActualDataVersions(c *gin.Context) {
	expId := c.Param("exp_id")
	year := c.Param("year")
	versions, err := services.GetActualDataVersions(expId, year)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, versions)
}

func GetExpLapseDataVersions(c *gin.Context) {
	year := c.Param("year")
	versions, err := services.GetLapseDataVersions(year)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, versions)
}

func GetExpMortalityDataVersions(c *gin.Context) {
	year := c.Param("year")
	versions, err := services.GetMortalityDataVersions(year)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, versions)
}

func CheckExpRunName(c *gin.Context) {

	runName := c.Param("run_name")
	exists, err := services.CheckRunName(runName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, exists)
}

func ExpRunAnalysis(c *gin.Context) {
	var runJobs []models.ExpAnalysisRunSetting
	if err := c.ShouldBindJSON(&runJobs); err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := c.MustGet("user").(models.AppUser)
	go func() {
		err := services.RunExpAnalysis(runJobs, user)
		if err != nil {
			fmt.Println("Error running exp analysis: ", err)
		}
	}()

	c.JSON(http.StatusOK, gin.H{"message": "Analysis started successfully"})
}

func DeleteExpRunGroup(c *gin.Context) {
	runGroupId := utils.StringToInt(c.Param("exp_run_id"))
	// we need to run this in a goroutine to avoid blocking the request
	go func() {
		err := services.DeleteExpRunGroup(runGroupId)
		if err != nil {
			fmt.Println("Error deleting exp run group: ", err)
		}
	}()

	c.JSON(http.StatusOK, nil)

}

func DeleteExpConfiguration(c *gin.Context) {
	configId := utils.StringToInt(c.Param("config_id"))
	err := services.DeleteExpConfiguration(configId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, nil)
}

func ExpGetAnalysisRuns(c *gin.Context) {
	// runs, err := services.GetAnalysisRunSettings()
	runs, err := services.GetAnalysisRunGroups()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, runs)
}

func ExpGetAnalysisRunSettings(c *gin.Context) {
	runId := utils.StringToInt(c.Param("id"))
	runSettings, err := services.GetAnalysisRunSettingsById(runId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, runSettings)
}

func GetExposureModelPointsByRunId(c *gin.Context) {
	runId := utils.StringToInt(c.Param("exp_run_id"))
	points, err := services.GetExposureModelPointsByRunGroupId(runId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, points)
}

func GetExposureModelPointsByRunIdAndProduct(c *gin.Context) {
	runId := c.Param("exp_run_id")
	product := c.Param("product")
	points, err := services.GetExposureModelPointsByRunIdAndProduct(runId, product)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, points)
}

func GetExpCrudeResultsByRunGroupId(c *gin.Context) {
	runGroupId := utils.StringToInt(c.Param("exp_run_id"))
	results, err := services.GetExpCrudeResultsByRunGroupId(runGroupId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, results)
}

func GetExpLapseCrudeResultsByRunGroupId(c *gin.Context) {
	runGroupId := utils.StringToInt(c.Param("exp_run_id"))
	results, err := services.GetExpLapseCrudeResultsByRunGroupId(runGroupId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, results)
}

func CalculateActualsVsExpected(c *gin.Context) {
	runId := utils.StringToInt(c.Param("exp_run_id"))
	results, err := services.CalculateActualsVsExpected(runId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, results)
}

func GetExpTableMetaData(c *gin.Context) {
	metaData := services.GetExpTableMetaData()
	c.JSON(http.StatusOK, metaData)
}

func UploadExpTable(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	tableType := form.Value["table_type"][0]
	year, err := strconv.Atoi(form.Value["year"][0])
	version := form.Value["version"][0]
	if err != nil {
		year = 0
	}

	var month = 12
	if form.Value["month"] != nil {
		month, _ = strconv.Atoi(form.Value["month"][0])
	}
	file := form.File["file"][0]

	fmt.Println(file.Filename)

	err = services.SaveExpTables(file, tableType, year, month, version)

	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, nil)
}

func DeleteExpTable(c *gin.Context) {
	table := c.Param("table_type")
	year := c.Param("year")
	version := c.Param("version")
	err := services.DeleteExpTableData(table, year, version)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, nil)
}

func DeleteExpConfigurationData(c *gin.Context) {
	tableType := c.Param("table_type")
	portfolioId := c.Param("portfolioId")
	year := c.Param("year")
	version := c.Param("version")

	err := services.DeleteExpConfigurationData(tableType, portfolioId, year, version)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, nil)

}
