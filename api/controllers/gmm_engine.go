package controllers

import (
	"api/models"
	"api/services"
	"api/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"net/http"
	"strconv"
	"time"
)

func RunPAAProjections(c *gin.Context) {
	user := c.MustGet("user").(models.AppUser)
	fmt.Println(user)
	var runJobGroups []models.MgmmRunPayload
	err := c.BindJSON(&runJobGroups)

	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	go func() {
		services.RunGMMProjectionGroups(runJobGroups, user)
	}()

	c.JSON(http.StatusOK, gin.H{"message": "projection runs have been submitted for processing"})

}

func CreateMGMMShockSetting(c *gin.Context) {
	var shockSetting models.ModifiedGMMShockSetting
	err := c.BindJSON(&shockSetting)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	shockSetting, err = services.SaveMGMMShockSetting(shockSetting)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, shockSetting)
}

func GetMGMMShockSettings(c *gin.Context) {
	settings, err := services.GetMGMMShockSettings()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, settings)
}

func DeleteMGMMShockSetting(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("setting_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	err = services.DeleteMGMMShockSetting(id)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, nil)
}
func GetMGMMShockBases(c *gin.Context) {

	bases, err := services.GetMGMMShockBases()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
	}

	c.JSON(http.StatusOK, bases)
}

func CreatePAAPortfolio(c *gin.Context) {
	var portfolio models.PaaPortfolio
	err := c.BindJSON(&portfolio)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	result := services.CreatePaaPortfolio(portfolio)
	c.JSON(http.StatusCreated, result)

}

func DeletePortfolio(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	err = services.DeletePortfolio(id)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err)
		return
	}

	c.JSON(http.StatusOK, nil)
}

func GetPortfolios(c *gin.Context) {
	results, err := services.GetPortfolios()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, results)
}

func GetRunPortfolios(c *gin.Context) {
	runDate := c.Param("run_date")
	results := services.GetRunPortfolios(runDate)
	c.JSON(http.StatusOK, results)
}

func GetRunPortfolioProducts(c *gin.Context) {
	runDate := c.Param("run_date")
	portfolioName := c.Param("portfolio_name")
	results := services.GetRunPortfolioProducts(runDate, portfolioName)
	c.JSON(http.StatusOK, results)
}

func GetRunPortfolioProductGroups(c *gin.Context) {
	runDate := c.Param("run_date")
	portfolioName := c.Param("portfolio_name")
	productCode := c.Param("product_code")
	results := services.GetRunPortfolioProductGroups(runDate, portfolioName, productCode)
	c.JSON(http.StatusOK, results)
}

func GetRunPortfolioProductGroup(c *gin.Context) {
	runDate := c.Param("run_date")
	portfolioName := c.Param("portfolio_name")
	productCode := c.Param("product_code")
	groupName := c.Param("group")
	results := services.GetRunPortfolioProductGroup(runDate, portfolioName, productCode, groupName)

	c.JSON(http.StatusOK, results)
}

func GetValidPortfolios(c *gin.Context) {
	results := services.GetValidPortfolios()
	c.JSON(http.StatusCreated, results)
}

func GetModelpointYearsForPortfolio(c *gin.Context) {
	portfolioId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
	}
	years := services.GetModelpointYearsForPortfolio(portfolioId)
	c.JSON(http.StatusOK, years)
}

func UploadMGMMParameters(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	year, _ := strconv.Atoi(form.Value["year"][0])
	file := form.File["file"][0]

	fmt.Println(year, file.Filename)

	err = services.SaveGMMParameters(file, year)

	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	years := services.GetGMMParameterYears()
	c.JSON(http.StatusOK, gin.H{"parameterYears": years})
}

func UploadMGMMModelPoints(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	fmt.Println("====received form")
	year, _ := strconv.Atoi(form.Value["year"][0])
	user := c.MustGet("user").(models.AppUser)
	file := form.File["file"][0]
	portfolioName := form.Value["portfolio_name"][0]
	portfolioId, _ := strconv.Atoi(form.Value["portfolio_id"][0])
	mpVersion := form.Value["mp_version"][0]

	fmt.Println(year, file.Filename)
	log.Print(year, file.Filename)

	yearVersion, err := services.SaveGMMModelPoints(file, year, portfolioName, portfolioId, mpVersion, user)

	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"year_version": yearVersion})
}

func UploadMGMMTables(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	user := c.MustGet("user").(models.AppUser)

	tableType := form.Value["table_type"][0]
	year, _ := strconv.Atoi(form.Value["year"][0])
	file := form.File["file"][0]
	var portfolioName string
	var yieldCurveCode string
	if form.Value["portfolio_name"] != nil {
		portfolioName = form.Value["portfolio_name"][0]
	}

	if form.Value["yield_curve_code"] != nil {
		yieldCurveCode = form.Value["yield_curve_code"][0]
	}

	var month int = 12
	if form.Value["month"] != nil {
		month, _ = strconv.Atoi(form.Value["month"][0])
	}

	fmt.Println(file.Filename)

	err = services.SaveGMMTables(file, tableType, year, month, portfolioName, yieldCurveCode, user)

	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, nil)
}

func DeleteMGMMTables(c *gin.Context) {
	tableType := c.Param("table_type")

	yieldCurveCode := c.Param("yield_curve_code")
	yieldCurveYear := utils.StringToInt(c.Param("yield_curve_year"))
	yieldCurveMonth := utils.StringToInt(c.Param("yield_curve_month"))
	year := utils.StringToInt(c.Param("year"))

	err := services.DeleteGMMTables(tableType, yieldCurveCode, yieldCurveYear, yieldCurveMonth, year)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
}

func GetGMMTableMetadata(c *gin.Context) {
	result := services.GetAllGMMTableMetaData()
	c.JSON(http.StatusOK, result)
}
func GetGMMTableData(c *gin.Context) {
	tableType := c.Param("table_type")
	results := services.GetGMMTableData(tableType)

	c.JSON(http.StatusOK, results)
}

func GetAllMGMMRunJobs(c *gin.Context) {
	results := services.GetAllRunJobs()
	c.JSON(http.StatusOK, results)
}

func GetAllMGMMRunJobsv2(c *gin.Context) {
	results := services.GetAllRunJobsV2()
	c.JSON(http.StatusOK, results)
}

func GetMGMMProjectionsForGroup(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		fmt.Println(err)
	}
	group := c.Param("group")

	gmmRunJob := services.GetGMMRunJob(id)
	if gmmRunJob.ID == 0 {
		c.JSON(http.StatusNotFound, nil)
		return
	}

	results := services.GetMGMMProjectionsByRun(id)

	//results := services.GetGMMProjectionsForGroup(id, group)
	scopedResults := services.GetMGMMScopedProjectionsByRunAndGroup(id, group)

	c.JSON(http.StatusOK, gin.H{"results": results, "scoped_results": scopedResults})
}

func GetMGMMProjections(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		fmt.Println(err)
	}

	gmmRunJob := services.GetGMMRunJob(id)
	if gmmRunJob.ID == 0 {
		c.JSON(http.StatusNotFound, nil)
		return
	}

	startResults := time.Now()
	results := services.GetMGMMProjectionsByRun(gmmRunJob.ID)
	endResults := time.Since(startResults)

	startGroups := time.Now()
	groups := services.GetMGMMIFRS17GroupsForProjection(gmmRunJob.ID)
	endGroups := time.Since(startGroups)

	startScopedResults := time.Now()
	scopedResults := services.GetMGMMScopedProjectionsByRun(gmmRunJob.ID)
	endScopedResults := time.Since(startScopedResults)

	startRunSettings := time.Now()
	runSettings := services.GetGMMRunSettings(id)
	endRunSettings := time.Since(startRunSettings)

	fmt.Println("Results: ", endResults)
	fmt.Println("Groups: ", endGroups)
	fmt.Println("Scoped Results: ", endScopedResults)
	fmt.Println("Run Settings: ", endRunSettings)

	c.JSON(http.StatusOK, gin.H{"results": results, "scoped_results": scopedResults, "run_settings": runSettings, "groups": groups})
}

func DeleteMGMMProjection(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		fmt.Println(err)
	}
	err = services.DeleteMGMMProjectionsRun(id)
	c.JSON(http.StatusOK, err)
}

func CheckFinanceYearExists(c *gin.Context) {
	portfolioName := c.Param("id")
	year, err := strconv.Atoi(c.Param("finance_year"))
	if err != nil {
		fmt.Println(err)
	}
	exists := services.CheckFinanceYearExists(portfolioName, year)
	c.JSON(http.StatusOK, exists)
}

func GetAvailableModelPointVersions(c *gin.Context) {
	portfolioName := c.Param("portfolio_name")
	year := c.Param("year")
	versions := services.GetAvailableModelPointVersions(portfolioName, year)
	c.JSON(http.StatusOK, versions)
}

func GetAvailablePAAYieldCurveYears(c *gin.Context) {
	results, err := services.GetAvailablePAAYieldCurveYears()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
	}
	c.JSON(http.StatusOK, results)
}

func GetAvailablePAAYieldCurveMonths(c *gin.Context) {
	year, err := strconv.Atoi(c.Param("year"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
	}
	results, err := services.GetAvailableYieldCurveMonths(year)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
	}
	c.JSON(http.StatusOK, results)
}

func GetModelPointsForPortfolio(c *gin.Context) {
	portfolioId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	//Don't need to check for error here as the year is optional
	var year int
	year, err = strconv.Atoi(c.Param("year"))
	if err != nil {
		year = 0
	}

	version := c.Param("version")

	data, err := services.GetPortfolioModelPoints(portfolioId, year, version)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, data)
}

func DeleteModelPointsForPortfolioAndYear(c *gin.Context) {
	portfolioId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	year, err := strconv.Atoi(c.Param("year"))
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	version := c.Param("version")

	err = services.DeleteModelPointsForPortfolioAndYear(portfolioId, year, version)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, nil)
}

func GetPAAYieldCurveCodes(c *gin.Context) {
	year := c.Param("year")
	results := services.GetPAAYieldCurveCodes(year)
	c.JSON(http.StatusOK, results)
}

func GeAvailablePAAParamYears(c *gin.Context) {
	portfolioName := c.Param("portfolio_name")
	results, err := services.GetAvailablePAAParameterYears(portfolioName)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, results)
}

func GetAvailablePAAPremiumPatternYears(c *gin.Context) {
	portfolioName := c.Param("portfolio_name")
	results, err := services.GetAvailablePAAPremiumPatternYears(portfolioName)
	if err != nil {
		//c.JSON(http.StatusBadRequest, err)
		//return
		fmt.Println(err)
	}
	c.JSON(http.StatusOK, results)
}

func CheckExistingGMMRunName(c *gin.Context) {
	runName := c.Param("run_name")
	results, err := services.CheckExistingGMMRunName(runName)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, results)
}

func CheckExistingPortfolioName(c *gin.Context) {
	portfolioName := c.Param("name")
	results, err := services.CheckExistingPortfolioName(portfolioName)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, results)
}
