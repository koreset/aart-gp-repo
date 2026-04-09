package controllers

import (
	"api/models"
	"api/services"
	"api/utils"
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func GetSpCodes(c *gin.Context) {
	productCode := c.Param("prod-code")

	result := services.GetSpCodesForProduct(productCode)

	c.JSON(http.StatusOK, result)
}

func GetIfrs17Groups(c *gin.Context) {
	productCode := c.Param("prod-code")

	result := services.GetIfrs17GroupsForProduct(productCode)

	c.JSON(http.StatusOK, result)
}

func SaveRAConfiguration(c *gin.Context) {
	var configs []models.RAConfiguration
	err := c.BindJSON(&configs)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, err)
		return
	}

	err = services.SaveRiskAdjustmentConfigs(configs)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err)
		return
	}

	c.JSON(http.StatusCreated, nil)
}

func SaveRiskDrivers(c *gin.Context) {
	var drivers []models.RiskDriver
	err := c.BindJSON(&drivers)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	err = services.SaveRiskDrivers(drivers)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err)
		return
	}
	c.JSON(http.StatusCreated, nil)
}

func DeleteRiskDriver(c *gin.Context) {
	productCode := c.Param("prod-code")
	if productCode == "" {
		c.JSON(http.StatusBadRequest, "product code is required")
		return
	}

	err := services.DeleteRiskDriver(productCode)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err)
		return
	}
	c.JSON(http.StatusOK, nil)

}

func GetRiskDrivers(c *gin.Context) {
	results, err := services.GetRiskDrivers()
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, results)
}

// GetRiskAdjustmentDriverTypes returns an array of objects: { risk_type: "..." }
func GetRiskAdjustmentDriverTypes(c *gin.Context) {
	types, err := services.GetRiskAdjustmentDriverTypes()
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	// Wrap into required shape
	resp := make([]map[string]string, 0, len(types))
	for _, t := range types {
		resp = append(resp, map[string]string{"risk_type": t})
	}
	c.JSON(http.StatusOK, resp)
}

// GetRiskAdjustmentFactorVersions returns the available versions for a given year
func GetRiskAdjustmentFactorVersions(c *gin.Context) {
	yearStr := c.Param("year")
	year, err := strconv.Atoi(yearStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid year"})
		return
	}

	versions, err := services.GetRiskAdjustmentFactorVersions(year)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, versions)
}

func GetRiskAdjustmentConfigs(c *gin.Context) {
	configs, err := services.GetRiskAdjustmentConfigs()
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, configs)
}

func SaveAosConfigurations(c *gin.Context) {
	var configSet models.AosVariableSet
	err := c.BindJSON(&configSet)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	configSet, err = services.SaveAosConfigurations(configSet)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err)
		return
	}
	c.JSON(http.StatusCreated, configSet)
}

func DeleteAosConfiguration(c *gin.Context) {
	configId, err := strconv.Atoi(c.Param("config_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	err = services.DeleteAosConfiguration(configId)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, nil)
}

func GetAosConfigurations(c *gin.Context) {
	configs, err := services.GetAosConfigurations()
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, configs)
}

func RunCsm(c *gin.Context) {
	var csmRuns []models.CsmRun
	user := c.MustGet("user").(models.AppUser)
	fmt.Println(user)
	err := c.BindJSON(&csmRuns)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, err)
		return
	}

	for i, _ := range csmRuns {
		res := strings.Split(csmRuns[i].RunDate, "-")
		csmRuns[i].RunDate = res[0] + "-" + res[1]
		csmRuns[i].CreationDate = time.Now()
		csmRuns[i].UserName = user.UserName
		csmRuns[i].UserEmail = user.UserEmail
		csmRuns[i].ProcessingStatus = "queued"

		services.DeleteExistingCsmRun(csmRuns[i])
		err = services.CreateCsmRun(&csmRuns[i])
	}

	go func() {
		for _, csmRun := range csmRuns {
			if csmRun.MeasurementType == "GMM" || csmRun.MeasurementType == "VFA" {
				services.RunCsmEngine(csmRun, false)
			} else {
				services.RunCsmEngineForPAA(csmRun)
			}
		}
	}()
	c.JSON(http.StatusOK, gin.H{"message": "jobs have been successfully submitted"})
}

func RerunCsm(c *gin.Context) {

	runId := utils.StringToInt(c.Param("id"))

	fmt.Println(runId)

	csmRun, err := services.GetCsmRunById(runId)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	fmt.Println(csmRun)

	if csmRun.MeasurementType == "GMM" || csmRun.MeasurementType == "VFA" {
		go func() {
			services.RerunCsmEngine(csmRun)

		}()
	}

	c.JSON(http.StatusOK, gin.H{"message": "job has been successfully submitted"})
}

func CheckExistingCsmRunName(c *gin.Context) {
	name := c.Param("name")

	result := services.CheckExistingCsmRunName(name)
	c.JSON(http.StatusOK, gin.H{"result": result})
}
func CheckExistingRun(c *gin.Context) {
	var run models.CsmRun
	err := c.BindJSON(&run)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	result := services.CheckExistingRun(run)
	c.JSON(http.StatusOK, gin.H{"result": result})
}

func GetCsmTableData(c *gin.Context) {
	tableName := c.Param("table")
	fmt.Println(tableName)
	results := services.GetCsmTableData(tableName)
	c.JSON(http.StatusOK, results)
}

func DeleteCsmTableData(c *gin.Context) {
	tableName := c.Param("table")
	year := c.Param("year")
	version := c.Param("version")
	err := services.DeleteCsmTableData(tableName, year, version)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, nil)
}

func GetCsmRunList(c *gin.Context) {
	results := services.GetCSMRunList()
	c.JSON(http.StatusOK, results)
}

func GetAllCsmRuns(c *gin.Context) {
	results, err := services.GetCSMRuns()
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, results)
}

func DeleteCsmRun(c *gin.Context) {
	runId, err := strconv.Atoi(c.Param("run_id"))
	runDate := c.Param("run_date")
	measurementType := c.Param("measurement_type")
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	err = services.DeleteCsmRun(runId, measurementType, runDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, nil)
}

func GetAosRunList(c *gin.Context) {
	results := services.GetAosRunList()
	c.JSON(http.StatusOK, results)
}

func GetPAARunList(c *gin.Context) {
	eligibilityMode := c.Param("eligibility-mode")
	results := services.GetPAARunList(eligibilityMode)
	c.JSON(http.StatusOK, results)

}

func GetCombinedRunListForDisclosures(c *gin.Context) {
	results := services.GetCombinedRunListForDisclosures()
	c.JSON(http.StatusOK, results)
}

func GetPortfoliosForDisclosures(c *gin.Context) {
	runDate := c.Param("run_date")
	results := services.GetPortfoliosForDisclosures(runDate)
	c.JSON(http.StatusOK, results)
}

func GetPortfolioProductsForDisclosures(c *gin.Context) {
	runDate := c.Param("run_date")
	portfolioName := c.Param("portfolio_name")
	results := services.GetPortfolioProductsForDisclosures(runDate, portfolioName)
	c.JSON(http.StatusOK, results)
}

func GetPortfolioProductGroupsForDisclosures(c *gin.Context) {
	runDate := c.Param("run_date")
	portfolioName := c.Param("portfolio_name")
	productCode := c.Param("product_code")
	results := services.GetPortfolioProductGroupsForDisclosures(runDate, portfolioName, productCode)
	c.JSON(http.StatusOK, results)
}

func GetCsmReportGroups(c *gin.Context) {
	runDate := c.Param("run-date")
	productCode := c.Param("product-code")
	results := services.GetReportGroups(runDate, productCode)
	c.JSON(http.StatusOK, results)
}

func GetCSMProductList(c *gin.Context) {
	runDate := c.Param("run-date")
	results := services.GetCSMProductList(runDate)
	c.JSON(http.StatusOK, results)
}

func GetAosProductList(c *gin.Context) {
	runDate := c.Param("run-date")
	measurementType := c.Param("measurement-type")
	results := services.GetAosProductList(runDate, measurementType)
	c.JSON(http.StatusOK, results)
}

func GetResultsForProduct(c *gin.Context) {
	prodCode := c.Param("prod_code")
	runId := c.Param("run-id")
	var results map[string]interface{}
	if prodCode != "All Products" {
		results = services.GetStepResultsForProduct(prodCode, runId)
	} else {
		results = services.GetAosStepResultsForAllProducts(runId)
		delete(results, "id")
	}
	c.JSON(http.StatusOK, results)
}

func GetStepResultsForProductByRunId(c *gin.Context) {
	prodCode := c.Param("prod_code")
	runId, err := strconv.Atoi(c.Param("run_id"))
	runDate := c.Param("run_date")
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	results := services.GetStepResultsForAosProductByRunId(prodCode, runId, runDate)
	c.JSON(http.StatusOK, results)
}

func GetAosResultsForCsmRun(c *gin.Context) {
	prodCode := c.Param("prod_code")
	runId, err := strconv.Atoi(c.Param("run_id"))
	runDate := c.Param("run_date")
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	var results map[string]interface{}
	if prodCode != "All Products" && prodCode != "" {
		results = services.GetAosStepResultsForProductByRunId(prodCode, runId)
	} else {
		results = services.GetAosStepResultsForAllProductsByRunId(runId, runDate)
		delete(results, "id")
	}
	c.JSON(http.StatusOK, results)
}

func GetPaaResultsForCsmRun(c *gin.Context) {
	prodCode := c.Param("prod_code")
	runId, err := strconv.Atoi(c.Param("run_id"))
	runDate := c.Param("run_date")
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	var results map[string]interface{}
	if prodCode != "All Products" && prodCode != "" {
		results = services.GetPaaStepResultsForProductByRunId(runId, prodCode, runDate)
	} else {
		results = services.GetPaaStepResultsForAllProductsByRunId(runId, runDate)
		delete(results, "id")
	}
	//results["run_settings"] = services.GetPaaRunSettings(runId)

	c.JSON(http.StatusOK, results)
}

func GetResultsForPAAProduct(c *gin.Context) {
	prodCode := c.Param("prod_code")
	runId := c.Param("run-id")
	var results map[string]interface{}
	if prodCode != "All Products" {
		results = services.GetStepResultsForPAAProduct(prodCode, runId)
	} else {
		results = services.GetPAAResultsForAllProducts(runId)
	}
	c.JSON(http.StatusOK, results)
}

func GetResultsForGroup(c *gin.Context) {
	group := c.Param("group")
	runDate := c.Param("run_date")
	results := services.GetStepResultsForGroup(runDate, group)
	c.JSON(http.StatusOK, results)
}

func GetStepResultsForAosProductGroupByRunId(c *gin.Context) {
	prodCode := c.Param("prod_code")
	group := c.Param("group")
	runId, err := strconv.Atoi(c.Param("run_id"))
	runDate := c.Param("run_date")
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	results := services.GetStepResultsForAosProductGroupByRunId(runId, runDate, prodCode, group)
	c.JSON(http.StatusOK, results)
}

func GetAllResultsForAosSteps(c *gin.Context) {
	runDate := c.Param("run-date")
	results := services.GetAosStepResultsForAllProducts(runDate)
	c.JSON(http.StatusOK, results)
}

func GetAllResultsForPAA(c *gin.Context) {
	runDate := c.Param("run-date")
	results := services.GetPAAResultsForAllProducts(runDate)
	c.JSON(http.StatusOK, results)
}

func GetResultsForPAAGroup(c *gin.Context) {
	group := c.Param("group")
	prodCode := c.Param("prod_code")
	runId := c.Param("run-id")
	results := services.GetStepResultsForPAAProductGroup(prodCode, runId, group)
	c.JSON(http.StatusOK, results)
}

func GetStepResultsForPAASingleProductGroupByRunId(c *gin.Context) {
	prodCode := c.Param("prod_code")
	group := c.Param("group")
	runDate := c.Param("run_date")
	runId, err := strconv.Atoi(c.Param("run_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	results := services.GetStepResultsForPAASingleProductGroup(runId, prodCode, group, runDate)
	c.JSON(http.StatusOK, results)
}

func GetRAFactors(c *gin.Context) {
	factors, err := services.GetRAFactors()
	if err != nil {
		c.JSON(http.StatusNotFound, nil)
		return
	}
	c.JSON(http.StatusOK, factors)
}

func UploadRAFactors(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	file := form.File["file"]
	year, err := strconv.Atoi(form.Value["year"][0])
	version := form.Value["version"][0]
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	ctx := context.WithValue(c, "keys", c.Keys)

	err = services.ProcessRiskAdjustmentFactorTables(ctx, file[0], year, version)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	factors, err := services.GetRAFactors()

	c.JSON(http.StatusOK, factors)
}

func UploadTransitionAdjustments(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	file := form.File["file"]
	year, err := strconv.Atoi(form.Value["year"][0])
	version := form.Value["version"][0]

	ctx := context.WithValue(c, "keys", c.Keys)

	err = services.ProcessTransitionAdjustments(ctx, file[0], year, version)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	factors, err := services.GetRAFactors()

	c.JSON(http.StatusOK, factors)
}

func UploadFinanceFile(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	file := form.File["file"]
	year, err := strconv.Atoi(form.Value["year"][0])
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	version := form.Value["version"][0]
	ctx := context.WithValue(c, "keys", c.Keys)

	err = services.ProcessFinanceVariables(ctx, file[0], year, version)
	results, err := services.GetFinanceFile()
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, results)
}

func UploadPaaFinanceFile(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	file := form.File["file"]
	year, err := strconv.Atoi(form.Value["year"][0])
	version := form.Value["version"][0]
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	user := c.MustGet("user").(models.AppUser)
	ctx := context.WithValue(c, "keys", c.Keys)
	ctx = context.WithValue(ctx, "user", user)

	err = services.ProcessPAAFinance(ctx, file[0], year, version)
	results, err := services.GetFinanceFile()
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, results)
}

func GetSapFileList(c *gin.Context) {
	results := services.GetSapFileList()
	c.JSON(http.StatusOK, results)
}

func GetSapResultsForRunName(c *gin.Context) {
	runName := c.Param("run_name")
	results, err := services.GetSapResultsForRunName(runName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, results)
}

func DeleteSapData(c *gin.Context) {
	runName := c.Param("run_name")
	err := services.DeleteSapData(runName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, "success")
}

func UploadSapFile(c *gin.Context) {
	user := c.MustGet("user").(models.AppUser)
	fmt.Println(user)

	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	file := form.File["file"]
	ctx := context.WithValue(c, "keys", c.Keys)

	err = services.ProcessSAPFile(ctx, file[0], user)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	fileList := services.GetSapFileList()

	c.JSON(http.StatusCreated, fileList)
}

func GetFinanceFile(c *gin.Context) {
	results, err := services.GetFinanceFile()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
	}
	c.JSON(http.StatusOK, results)
}

func GetFinanceAndRaYears(c *gin.Context) {
	measure := c.Param("measure")
	paaRunId := c.Param("paa_run_id")
	years, err := services.GetFinanceAndRaYears(measure, paaRunId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
	}
	c.JSON(http.StatusOK, years)
}

func GetAvailablePaaFinanceVersions(c *gin.Context) {
	year := c.Param("year")
	versions, err := services.GetAvailablePaaFinanceVersions(year)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
	}
	c.JSON(http.StatusOK, versions)
}

func GetAvailableGmmFinanceVersions(c *gin.Context) {
	year := c.Param("year")
	versions, err := services.GetAvailableGmmFinanceVersions(year)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
	}
	c.JSON(http.StatusOK, versions)
}

func GetLiabilityMovementReport(c *gin.Context) {
	runDate := c.Param("run_date")
	productCode := c.Param("product_code")
	ifrs17Group := c.Param("ifrs17_group")
	movements := services.GetLiabilityMovements(runDate, productCode, ifrs17Group)
	c.JSON(http.StatusOK, movements)
}

func GetLiabilityMovementRunList(c *gin.Context) {
	runlist := services.GetLiabilityMovementRunList()
	c.JSON(http.StatusOK, runlist)
}

func GetLiabilityMovementsProducts(c *gin.Context) {
	runDate := c.Param("run_date")
	products := services.GetLiabilityMovementsProducts(runDate)
	c.JSON(http.StatusOK, products)
}

func GetLiabilityMovementsProductGroups(c *gin.Context) {
	productCode := c.Param("product_code")
	runDate := c.Param("run_date")
	products := services.GetLiabilityMovementsProductGroups(runDate, productCode)
	c.JSON(http.StatusOK, products)
}

func GetLiabilityMovementsGroups(c *gin.Context) {
	productCode := c.Param("product_code")
	runDate := c.Param("run_date")
	group := c.Param("ifrs17_group")
	products := services.GetLiabilityMovements(runDate, productCode, group)
	c.JSON(http.StatusOK, products)
}

func GetInsuranceRevenueAnalysis(c *gin.Context) {
	runDate := c.Param("run_date")
	productCode := c.Param("product_code")
	ifrs17Group := c.Param("ifrs17_group")
	insuranceRevenues := services.GetInsuranceRevenueAnalysis(runDate, productCode, ifrs17Group)
	c.JSON(http.StatusOK, insuranceRevenues)
}

func GetInsuranceRevenueAnalysisRuns(c *gin.Context) {
	runs := services.GetInsuranceRevenueAnalysisRuns()
	c.JSON(http.StatusOK, runs)
}

func GetInitialRecognitionAnalysisRuns(c *gin.Context) {
	runs := services.GetInitialRecognitionAnalysisRuns()
	c.JSON(http.StatusOK, runs)
}

func GetInitialRecognitionAnalysis(c *gin.Context) {
	runDate := c.Param("run_date")
	productCode := c.Param("product_code")
	ifrs17Group := c.Param("ifrs17_group")
	initialRecognition := services.GetInitialRecognitionAnalysis(runDate, productCode, ifrs17Group)
	c.JSON(http.StatusOK, initialRecognition)
}

func GetCSMProjectionsRunList(c *gin.Context) {
	runs := services.GetCSMProjectionsRunList()
	c.JSON(http.StatusOK, runs)
}

func GetCsmProjections(c *gin.Context) {
	runDate := c.Param("run_date")
	productCode := c.Param("product_code")
	ifrs17Group := c.Param("ifrs17_group")
	results := services.GetCsmProjections(runDate, productCode, ifrs17Group)
	c.JSON(http.StatusOK, results)
}

func GenerateTransactionReports(c *gin.Context) {
	//productCode := c.Param("product_code")
	var report []models.JTReportEntry
	ifrs17Group := c.Param("ifrs17_group")
	runDate := c.Param("run_date")
	if ifrs17Group != "" {
		//assume report by group
		report = services.GenerateTransactionReports(ifrs17Group, runDate)
	} else {
		//report = services.GenerateAllTransactionReports()
	}

	c.JSON(http.StatusOK, report)
}

func GenerateSubLedgerReportsforAllProducts(c *gin.Context) {
	runDate := c.Param("run_date")
	report := services.GenerateSubLedgerReportsforAllProducts(runDate)

	c.JSON(http.StatusOK, report)
}
func GenerateSubLedgerReportsforOneProduct(c *gin.Context) {
	runDate := c.Param("run_date")
	prodCode := c.Param("product_code")
	report := services.GenerateSubLedgerReportsforOneProduct(runDate, prodCode)

	c.JSON(http.StatusOK, report)
}

func GenerateSubLedgerReports(c *gin.Context) {
	//productCode := c.Param("product_code")
	ifrs17Group := c.Param("ifrs17_group")
	runDate := c.Param("run_date")
	report := services.GenerateSubLedgerReports(ifrs17Group, runDate)

	c.JSON(http.StatusOK, report)
}

func GenerateTrialBalanceReports(c *gin.Context) {
	//productCode := c.Param("product_code")
	ifrs17Group := c.Param("ifrs17_group")
	runDate := c.Param("run_date")
	report := services.GenerateTrialBalanceReports(ifrs17Group, runDate)

	c.JSON(http.StatusOK, report)
}

func GenerateBalanceSheetReportsforAllProducts(c *gin.Context) {
	runDate := c.Param("run_date")
	report := services.GenerateBalanceSheetReportsforAllProducts(runDate)

	c.JSON(http.StatusOK, report)
}

func GenerateBalanceSheetReportForProduct(c *gin.Context) {
	runDate := c.Param("run_date")
	prodCode := c.Param("product_code")
	report := services.GenerateBalanceSheetReportForProduct(runDate, prodCode)

	c.JSON(http.StatusOK, report)
}

func GenerateBalanceSheetReportForProductGroup(c *gin.Context) {
	runDate := c.Param("run_date")
	prodCode := c.Param("product_code")
	group := c.Param("group")
	report := services.GenerateBalanceSheetReportForProductGroup(runDate, prodCode, group)

	c.JSON(http.StatusOK, report)
}

func GenerateBalanceSheetSummariesForAll(c *gin.Context) {
	runDate := c.Param("run_date")
	productCode := ""
	group := ""
	report := services.GenerateBalanceSheetSummariesForAll(runDate, productCode, group)

	c.JSON(http.StatusOK, report)
}

func GenerateBalanceSheetSummariesForProduct(c *gin.Context) {
	runDate := c.Param("run_date")
	prodCode := c.Param("product_code")
	group := ""
	report := services.GenerateBalanceSheetSummariesForAll(runDate, prodCode, group)

	c.JSON(http.StatusOK, report)
}

func GenerateBalanceSheetSummariesForProductGroup(c *gin.Context) {
	runDate := c.Param("run_date")
	prodCode := c.Param("product_code")
	group := c.Param("group")
	report := services.GenerateBalanceSheetSummariesForAll(runDate, prodCode, group)

	c.JSON(http.StatusOK, report)

}

func GenerateTrialBalanceReportsforAllProducts(c *gin.Context) {
	runDate := c.Param("run_date")
	report := services.GenerateTrialBalanceReportsforAllProducts(runDate)
	c.JSON(http.StatusOK, report)
}

func GenerateTrialBalanceReportsforOneProduct(c *gin.Context) {
	runDate := c.Param("run_date")
	prodCode := c.Param("product_code")
	report := services.GenerateTrialBalanceReportsforOneProduct(runDate, prodCode)
	c.JSON(http.StatusOK, report)
}

func GetJournalEntryProducts(c *gin.Context) {
	runDate := c.Param("run-date")
	results := services.GetJournalEntryProducts(runDate)
	c.JSON(http.StatusOK, results)
}

func GetJournalEntryForAll(c *gin.Context) {
	runDate := c.Param("run-date")
	results := services.GenerateJournalsforAll(runDate)
	c.JSON(http.StatusOK, results)
}

func GetJournalEntryForOneProduct(c *gin.Context) {
	runDate := c.Param("run-date")
	prodCode := c.Param("product_code")
	results := services.GenerateJournalsforOneProduct(runDate, prodCode)
	c.JSON(http.StatusOK, results)
}

func GetBelBuildupForAll(c *gin.Context) {
	runDate := c.Param("run_date")
	results := services.GenerateBelBuildupforAll(runDate)
	c.JSON(http.StatusOK, results)
}

func GetBelBuildupForOneProduct(c *gin.Context) {
	runDate := c.Param("run_date")
	prodCode := c.Param("product_code")
	results := services.GenerateBelBuildupOneProduct(runDate, prodCode)
	c.JSON(http.StatusOK, results)
}

func GetBelBuildupForOneProductGroup(c *gin.Context) {
	runDate := c.Param("run_date")
	productCode := c.Param("product_code")
	ifrs17Group := c.Param("ifrs17_group")
	if ifrs17Group != "" {
		//assume report by group
		results := services.GenerateBelBuildupOneGroup(ifrs17Group, productCode, runDate)
		c.JSON(http.StatusOK, results)
	}

}

func UploadIFRS17Tables(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	file := form.File["file"]
	ctx := context.WithValue(c, "keys", c.Keys)
	tableType := form.Value["table_type"][0]
	year := form.Value["year"][0]
	version := form.Value["version"][0]

	err = services.ProcessIFRS17Tables(ctx, file[0], tableType, year, version)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	//factors, err := services.GetRAFactors()

	c.JSON(http.StatusOK, nil)
}

// ─── Run Approval Workflow Handlers ─────────────────────────────────────────

func ReviewCsmRun(c *gin.Context) {
	runID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "invalid run id"})
		return
	}
	var req models.CsmRunApprovalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	user := c.MustGet("user").(models.AppUser)
	if err := services.ReviewCsmRun(runID, req.Notes, user); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Run marked as reviewed"})
}

func ApproveCsmRun(c *gin.Context) {
	runID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "invalid run id"})
		return
	}
	var req models.CsmRunApprovalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	user := c.MustGet("user").(models.AppUser)
	if err := services.ApproveCsmRun(runID, req.Notes, user); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Run approved"})
}

func LockCsmRun(c *gin.Context) {
	runID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "invalid run id"})
		return
	}
	user := c.MustGet("user").(models.AppUser)
	if err := services.LockCsmRun(runID, user); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Run locked"})
}

func ReturnCsmRunToDraft(c *gin.Context) {
	runID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "invalid run id"})
		return
	}
	var req models.CsmRunApprovalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	user := c.MustGet("user").(models.AppUser)
	if err := services.ReturnCsmRunToDraft(runID, req.Reason, user); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Run returned to draft"})
}

func CompareRuns(c *gin.Context) {
	runAId, _ := strconv.Atoi(c.Query("run_a"))
	runBId, _ := strconv.Atoi(c.Query("run_b"))
	if runAId == 0 || runBId == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "run_a and run_b query params required"})
		return
	}
	result, err := services.CompareRuns(runAId, runBId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": result})
}

//func GetCsmExcelExport(c *gin.Context) {
//	runId := c.Param("csm_run_id")
//	//productCode := c.Param("product_code")
//	//ifrs17Group := c.Param("ifrs17_group")
//	//exportType := c.Param("export_type")
//
//	data, err := services.GetCsmExcelExport(runId, productCode, ifrs17Group)
//	if err != nil {
//		c.JSON(http.StatusInternalServerError, err.Error())
//		return
//	}
//	//c.Data("csm_export.xlsx", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", data)
//	//} else if exportType == "paa" {
//	//	data, err := services.GetPaaExcelExport(runId, productCode, ifrs17Group)
//	//	if err != nil {
//	//		c.JSON(http.StatusInternalServerError, err.Error())
//	//		return
//	//	}
//	//	c.Data("paa_export.xlsx", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", data)
//	//} else {
//	//	c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid export type"})
//	//}
//}

// GetIFRS17AuditLog returns audit log entries for the IFRS 17 module.
// Query params: event_type (optional), from (optional, YYYY-MM-DD), to (optional, YYYY-MM-DD).
func GetIFRS17AuditLog(c *gin.Context) {
	eventType := c.Query("event_type")
	from := c.Query("from")
	to := c.Query("to")
	logs, err := services.GetIFRS17AuditLog(eventType, from, to)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": logs})
}

// ─── Task 2.4: Transition Adjustments Review ────────────────────────────────

// GetTransitionAdjustments returns all transition adjustments, optionally filtered by status.
func GetTransitionAdjustments(c *gin.Context) {
	status := c.Query("status")
	records, err := services.GetTransitionAdjustments(status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": records})
}

// ApproveTransitionAdjustment approves a single transition adjustment.
func ApproveTransitionAdjustment(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "invalid id"})
		return
	}
	var req models.TransitionAdjustmentReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	user := c.MustGet("user").(models.AppUser)
	if err := services.ApproveTransitionAdjustment(id, req.Notes, user); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Adjustment approved"})
}

// RejectTransitionAdjustment rejects a single transition adjustment.
func RejectTransitionAdjustment(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "invalid id"})
		return
	}
	var req models.TransitionAdjustmentReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	user := c.MustGet("user").(models.AppUser)
	if err := services.RejectTransitionAdjustment(id, req.Notes, user); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Adjustment rejected"})
}

// BulkApproveTransitionAdjustments approves multiple transition adjustments.
func BulkApproveTransitionAdjustments(c *gin.Context) {
	var body struct {
		IDs []int `json:"ids"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	user := c.MustGet("user").(models.AppUser)
	if err := services.BulkApproveTransitionAdjustments(body.IDs, user); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Adjustments approved"})
}

// ─── Task 2.5: IFRS17 Restatement / Amendment Workflow ──────────────────────

// CreateIFRS17Amendment creates a new amendment record.
func CreateIFRS17Amendment(c *gin.Context) {
	var req models.CreateAmendmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	user := c.MustGet("user").(models.AppUser)
	amendment, err := services.CreateIFRS17Amendment(req, user)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": amendment})
}

// GetAmendmentsForRun returns all amendments for a given run ID.
func GetAmendmentsForRun(c *gin.Context) {
	runID, err := strconv.Atoi(c.Param("run_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "invalid run_id"})
		return
	}
	records, err := services.GetAmendmentsForRun(runID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": records})
}

// GetAllIFRS17Amendments returns all amendments optionally filtered by status.
func GetAllIFRS17Amendments(c *gin.Context) {
	status := c.Query("status")
	records, err := services.GetAllAmendments(status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": records})
}

// ApproveIFRS17Amendment approves an amendment.
func ApproveIFRS17Amendment(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "invalid id"})
		return
	}
	user := c.MustGet("user").(models.AppUser)
	if err := services.ApproveIFRS17Amendment(id, user); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Amendment approved"})
}

// GetConsolidatedDisclosure returns aggregated IFRS 17 metrics across all approved/locked runs.
func GetConsolidatedDisclosure(c *gin.Context) {
	data, err := services.GetConsolidatedDisclosure()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.PremiumResponse{
			Success: false, Message: err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": data})
}

func GetSCRRABridge(c *gin.Context) {
	runIDStr := c.Param("run_id")
	runID, _ := strconv.Atoi(runIDStr)
	rows, err := services.GetSCRRABridge(runID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": rows})
}

func CreateSCREntry(c *gin.Context) {
	var req models.CreateSCREntryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	user := c.MustGet("user").(models.AppUser)
	entry, err := services.CreateSCREntry(req, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": entry})
}

func GetSCREntries(c *gin.Context) {
	product := c.Query("product_code")
	entries, err := services.GetSCREntries(product)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": entries})
}

func DeleteSCREntry(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)
	if err := services.DeleteSCREntry(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func GetSARBMappings(c *gin.Context) {
	mappings, err := services.GetSARBMappings()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": mappings})
}

func UpsertSARBMapping(c *gin.Context) {
	var req models.UpsertSARBMappingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	user := c.MustGet("user").(models.AppUser)
	mapping, err := services.UpsertSARBMapping(req, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": mapping})
}

func DeleteSARBMapping(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)
	if err := services.DeleteSARBMapping(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func GenerateSARBReport(c *gin.Context) {
	runIDStr := c.Param("run_id")
	runID, _ := strconv.Atoi(runIDStr)
	rows, err := services.GenerateSARBReport(runID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": rows})
}

func SeedSARBMappings(c *gin.Context) {
	user := c.MustGet("user").(models.AppUser)
	if err := services.SeedDefaultSARBMappings(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Default mappings seeded"})
}

func ComputeDeferredTax(c *gin.Context) {
	var req models.ComputeDeferredTaxRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	entries, summary, err := services.ComputeDeferredTax(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": entries, "summary": summary})
}

func GetDeferredTax(c *gin.Context) {
	runIDStr := c.Param("run_id")
	runID, _ := strconv.Atoi(runIDStr)
	entries, summary, err := services.GetDeferredTaxEntries(runID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": entries, "summary": summary})
}
