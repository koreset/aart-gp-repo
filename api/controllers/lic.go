package controllers

import (
	"api/models"
	"api/services"
	"api/utils"
	"encoding/csv"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jszwec/csvutil"
	"mime/multipart"
	"net/http"
	"strconv"
)

func CreateIbnrPortfolio(c *gin.Context) {
	var licPortfolio models.LicPortfolio
	c.BindJSON(&licPortfolio)
	licPortfolio, err := services.CreateLicPortfolio(licPortfolio)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, licPortfolio)
}

func DeleteIbnrPortfolio(c *gin.Context) {
	portfolioId, err := strconv.Atoi(c.Param("portfolio_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	//delete related claims data
	err = services.DeleteIbnrPortfolio(portfolioId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "lic portfolio deleted"})
}

func GetIbnrPortfolios(c *gin.Context) {
	licPortfolios, err := services.GetIbnrPortfolios()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, licPortfolios)
}

func GetValidLicPortfolios(c *gin.Context) {
	validLicPortfolios, err := services.GetValidLicPortfolios()
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, validLicPortfolios)
}

func GetLicPortfolioClaimYears(c *gin.Context) {
	portfolioName := c.Param("portfolio_name")
	licPortfolioClaimYears, err := services.GetLicPortfolioClaimYears(portfolioName)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, licPortfolioClaimYears)
}

func GetLicPortfolioParameterYears(c *gin.Context) {
	portfolioName := c.Param("portfolio_name")
	licPortfolioParameterYears, err := services.GetLicPortfolioParameterYears(portfolioName)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, licPortfolioParameterYears)
}

func GetLicBaseVariables(c *gin.Context) {
	licVariables, err := services.GetLicBaseVariables()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, licVariables)
}

func GetIbnrTableMetaData(c *gin.Context) {
	licTableMetaData := services.GetIbnrTableMetaData()
	c.JSON(http.StatusOK, licTableMetaData)
}

func GetLicTableMetaData(c *gin.Context) {
	licTableMetaData := services.GetLicTableMetaData()
	c.JSON(http.StatusOK, licTableMetaData)
}

func DeleteLicTable(c *gin.Context) {
	tableType := c.Param("table_type")
	err := services.DeleteLicTable(tableType)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, nil)
}

func DeleteIbnrTable(c *gin.Context) {
	tableType := c.Param("table_type")
	err := services.DeleteIbnrTable(tableType)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, nil)
}

func DeleteIbnrTableData(c *gin.Context) {
	tableType := c.Param("table_type")
	portfolio := c.Param("portfolio")
	year := c.Param("year")
	version := c.Param("version")

	err := services.DeleteIbnrTableData(tableType, portfolio, year, version)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, nil)

}

func SaveLicTable(c *gin.Context) {
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
	file := form.File["file"][0]

	err = services.SaveLicTables(file, tableType, year, version)

	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, nil)
}

func SaveIbnrTable(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	tableType := form.Value["table_type"][0]
	year, err := strconv.Atoi(form.Value["year"][0])
	if err != nil {
		year = 0
	}

	var month = 12
	if form.Value["month"] != nil {
		month, _ = strconv.Atoi(form.Value["month"][0])
	}
	file := form.File["file"][0]
	//productName := form.Value["product_name"][0]
	yieldCurveCode := form.Value["yield_curve_code"][0]

	fmt.Println(file.Filename)

	err = services.SaveIbnrTables(file, tableType, year, month, yieldCurveCode)

	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, nil)
}

func UploadIbnrPortfolioData(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	year, _ := strconv.Atoi(form.Value["year"][0])
	file := form.File["file"][0]
	portfolioName := form.Value["portfolio_name"][0]
	tableType := form.Value["type"][0]
	portfolioId, _ := strconv.Atoi(form.Value["portfolio_id"][0])
	versionName := form.Value["version_name"][0]

	fmt.Println(year, file.Filename)

	err = services.SaveLicClaimsData(file, year, portfolioName, portfolioId, tableType, versionName)

	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, nil)
}

func GetIbnrClaimsData(c *gin.Context) {
	portfolioName := c.Param("portfolio_name")
	tableType := c.Param("type")
	year := c.Param("year")
	version := c.Param("version")
	licClaimsData, err := services.GetLicClaimsData(portfolioName, tableType, year, version)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, licClaimsData)
}

func GetLicTableData(c *gin.Context) {
	tableType := c.Param("table_type")
	licTableData, err := services.GetLicTableData(tableType)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, licTableData)
}

func GetIbnrTableData(c *gin.Context) {
	tableType := c.Param("table_type")
	licTableData, err := services.GetIbnrTableData(tableType)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, licTableData)
}

func RunIBNRJobs(c *gin.Context) {
	user := c.MustGet("user").(models.AppUser)
	var jobs []models.IBNRRunSetting
	err := c.BindJSON(&jobs)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	go func() {
		services.RunIBNR(jobs, user, nil, "")
	}()

	c.JSON(http.StatusOK, gin.H{"message": "jobs started"})
}

func RunManualIBNR(c *gin.Context) {
	user := c.MustGet("user").(models.AppUser)
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	file := form.File["file"][0]
	runId, err := strconv.Atoi(form.Value["run_id"][0])
	manualProdCode := form.Value["prod_code"][0]
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	fmt.Println(file.Filename)
	fmt.Println(runId)
	//We assume that the file is a csv file for LICDevelopmentFactors
	// check that the columns are correct and match LICDevelopmentFactors fields
	err = checkMatchingTags(file, models.LicDevelopmentFactor{})
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	go func() {
		err = services.RunIbnrWithManualFactors(file, runId, manualProdCode, user)
		if err != nil {
			fmt.Println(err)
		}
	}()

	c.JSON(http.StatusOK, gin.H{"message": "jobs started"})
}

func checkMatchingTags(fileHeader *multipart.FileHeader, factor models.LicDevelopmentFactor) error {
	var delimiter rune
	delimiterFile, err := fileHeader.Open()
	if err != nil {
		return err
	}
	defer delimiterFile.Close()
	delimiter, err = utils.GetDelimiter(delimiterFile)

	file, err := fileHeader.Open()
	if err != nil {
		return err
	}

	reader := csv.NewReader(file)
	reader.TrimLeadingSpace = true
	reader.Comma = delimiter
	dec, err := csvutil.NewDecoder(reader)
	if err != nil {
		fmt.Println(err)
	}
	dec.Header()

	if !utils.MatchCSVTags(dec.Header(), models.LicDevelopmentFactor{}) {
		return errors.New("the column names in the file does not match the column names in the database")
	}
	return nil
}

func GetAllIBNRRunJobs(c *gin.Context) {
	jobs, err := services.GetAllLICRunJobs()
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, jobs)
}

func DeleteIBNRRunJob(c *gin.Context) {
	jobId, err := strconv.Atoi(c.Param("job_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = services.DeleteIBNRRunJob(jobId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "lic run job deleted"})
}

func GetIBNRRunResults(c *gin.Context) {
	jobId, err := strconv.Atoi(c.Param("job_id"))
	resultType := c.Param("result_type")
	product := c.Param("product")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	results, err := services.GetIBNRRunResults(jobId, resultType, product)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, results)
}

func GetLicPortfolioProducts(c *gin.Context) {
	portfolioID, err := strconv.Atoi(c.Param("portfolio_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	productCodes, err := services.GetLicPortfolioProducts(portfolioID)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, productCodes)
}

// GetIbnrRunProducts returns the distinct product codes for the portfolio
// associated with a given IBNR run ID. RunDetail.vue uses this so it never
// needs to know the portfolio ID separately.
func GetIbnrRunProducts(c *gin.Context) {
	runID, err := strconv.Atoi(c.Param("run_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid run id"})
		return
	}
	productCodes, err := services.GetLicProductsForIbnrRun(runID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, productCodes)
}

func CheckIbnrRunName(c *gin.Context) {
	runName := c.Param("run_name")
	found := services.CheckIbnrRunName(runName)
	c.JSON(http.StatusOK, found)
}

func ChecklicRunName(c *gin.Context) {
	runName := c.Param("run_name")
	found := services.CheckLicRunName(runName)
	c.JSON(http.StatusOK, found)
}

func CheckLicConfigName(c *gin.Context) {
	configName := c.Param("config_name")
	found := services.CheckLicConfigName(configName)
	c.JSON(http.StatusOK, found)
}

func SaveLicConfiguration(c *gin.Context) {
	var licConfigurations models.LicVariableSet
	err := c.BindJSON(&licConfigurations)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = services.SaveLicConfiguration(licConfigurations)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, nil)
}

func GetLicVariableSets(c *gin.Context) {
	licVariableSets, err := services.GetLicVariableSets()
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, licVariableSets)
}

func RunLicJobs(c *gin.Context) {
	user := c.MustGet("user").(models.AppUser)
	var licRunSettings []models.LicRunSetting
	err := c.BindJSON(&licRunSettings)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	for i, _ := range licRunSettings {
		licRunSettings[i].UserName = user.UserName
		licRunSettings[i].UserEmail = user.UserEmail
	}
	go func() {
		services.RunLic(licRunSettings)
	}()
	c.JSON(http.StatusOK, nil)
}

func CreateIBNRShockSetting(c *gin.Context) {
	var shockSetting models.IBNRShockSetting
	err := c.BindJSON(&shockSetting)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	shockSetting, err = services.SaveIBNRShockSetting(shockSetting)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, shockSetting)
}

func GetIBNRShockSettings(c *gin.Context) {
	settings, err := services.GetIBNRShockSettings()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, settings)
}

func DeleteIBNRShockSetting(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("setting_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	err = services.DeleteIBNRShockSetting(id)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, nil)
}

// ── IBNR Method Assignments ───────────────────────────────────────────────────

func GetMethodAssignments(c *gin.Context) {
	portfolioID, err := strconv.Atoi(c.Param("portfolio_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid portfolio_id"})
		return
	}
	assignments, err := services.GetMethodAssignments(portfolioID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": assignments})
}

func SaveMethodAssignment(c *gin.Context) {
	var assignment models.LicIbnrMethodAssignment
	if err := c.BindJSON(&assignment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Inject audit fields from authenticated user context
	if email, ok := c.Get("userEmail"); ok {
		assignment.CreatedBy = fmt.Sprintf("%v", email)
	}
	if name, ok := c.Get("userName"); ok {
		assignment.UserName = fmt.Sprintf("%v", name)
	}
	saved, err := services.SaveMethodAssignment(assignment)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": saved})
}

func DeleteMethodAssignment(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := services.DeleteMethodAssignment(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

type autoAssignRequest struct {
	PortfolioID    int     `json:"portfolio_id"`
	PortfolioName  string  `json:"portfolio_name"`
	ClaimsYear     string  `json:"claims_year"`
	Version        string  `json:"version"`
	MinDataYears   int     `json:"min_data_years"`
	MaxCVThreshold float64 `json:"max_cv_threshold"`
	FallbackMethod string  `json:"fallback_method"`
}

func AutoAssignMethods(c *gin.Context) {
	var req autoAssignRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if req.MinDataYears == 0 {
		req.MinDataYears = 5
	}
	if req.MaxCVThreshold == 0 {
		req.MaxCVThreshold = 0.5
	}
	if req.FallbackMethod == "" {
		req.FallbackMethod = "bornhuetter-ferguson"
	}
	createdBy, _ := c.Get("userEmail")
	userName, _ := c.Get("userName")
	results, err := services.AutoAssignMethods(
		req.PortfolioID, req.PortfolioName,
		req.ClaimsYear, req.Version,
		req.MinDataYears, req.MaxCVThreshold, req.FallbackMethod,
		fmt.Sprintf("%v", createdBy), fmt.Sprintf("%v", userName),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": results})
}

// ─────────────────────────────────────────────────────────────────────────────

func GetModelPointRows(c *gin.Context) {
	modelPointRows := services.GetModelPointRows()
	c.JSON(http.StatusOK, modelPointRows)
}

func GetLicBuildUps(c *gin.Context) {
	runId, err := strconv.Atoi(c.Param("run_id"))
	portfolioName := c.Param("portfolio_name")
	productCode := c.Param("product_code")

	buildUps, err := services.GetLicBuildUp(runId, portfolioName, productCode)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, buildUps)
}

func GetLicRuns(c *gin.Context) {
	runs, err := services.GetLicRuns()
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, runs)
}

func DeleteLicRun(c *gin.Context) {
	runId, err := strconv.Atoi(c.Param("run_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	err = services.DeleteLicRun(runId)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, nil)
}
func DeleteLicConfig(c *gin.Context) {
	configId, err := strconv.Atoi(c.Param("config_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	err = services.DeleteLicConfig(configId)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, nil)
}

func CheckExistingLicRun(c *gin.Context) {
	//var run models.CsmRun
	var run models.LicRun
	err := c.BindJSON(&run)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, err)
		return
	}

	result := services.CheckExistingLicRun(run)
	c.JSON(http.StatusOK, gin.H{"result": result})
}

func GetAvailableInputVersions(c *gin.Context) {
	portfolioName := c.Param("portfolio_name")
	year := c.Param("year")
	versions, err := services.GetAvailableInputVersions(portfolioName, year)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, versions)
}

func GetAvailableLicParameterYears(c *gin.Context) {
	years, err := services.GetAvailableLicParameterYears()
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, years)
}

func GetAvailableLicParameterVersions(c *gin.Context) {
	year := c.Param("year")
	versions, err := services.GetAvailableLicParameterVersions(year)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, versions)
}
