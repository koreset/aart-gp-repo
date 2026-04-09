package controllers

import (
	"api/services"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

func GetMonthlyConsolidated(c *gin.Context) {
	results := services.GetMonthlyConsolidatedResults()
	c.JSON(http.StatusOK, results)
}

func GetAnnualConsolidated(c *gin.Context) {
	results := services.GetAnnualConsolidatedResults()
	c.JSON(http.StatusOK, results)
}

func GetReportingDisclosures(c *gin.Context) {
	//We need current year and past year values
	prodId, err := strconv.Atoi(c.Param("prodId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "A valid product id is required"})
		return
	}
	prod, err := services.GetProductById(prodId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "A valid product matching the supplied id was not found"})
		return
	}
	currentYear := time.Now().Year()
	pastYear := currentYear - 1
	//yield := services.GetYieldCurveDisclosure(currentYear, pastYear)
	lapse := services.GetLapseDisclosure(prod.ProductCode, currentYear, pastYear)
	mortality := services.GetMortalityDisclosure(prod.ProductCode, currentYear, pastYear)
	mortalityAccidental := services.GetMortalityAccidentalDisclosure(prod.ProductCode, currentYear, pastYear)
	// basis="" → picks first matching row (backward-compatible default)
	parameters := services.GetParameterDisclosure(prod.ProductCode, currentYear, pastYear, "")
	modelPointYears := services.GetAvailableModelPointYears(prod.ProductCode)
	mortalityYears := services.GetAvailableMortalityYears(prod.ProductCode)
	mortalityAccidentalYears := services.GetAvailableMortalityAccidentalYears(prod.ProductCode)

	results := make(map[string]interface{})
	//results["yield"] = yield
	results["lapse"] = lapse
	results["mortality"] = mortality
	results["mortality_accidental"] = mortalityAccidental
	results["parameters"] = parameters
	results["model_point_years"] = modelPointYears
	results["mortality_years"] = mortalityYears
	results["mortality_accidental_years"] = mortalityAccidentalYears

	c.JSON(http.StatusOK, results)
}

func GetReserveSummary(c *gin.Context) {
	runDate := c.Param("run_date")
	items := services.GetReserveSummary(runDate)
	c.JSON(http.StatusOK, items)
}

func GetPAAReserveSummary(c *gin.Context) {
	runDate := c.Param("run_date")
	items := services.GetPAAReserveSummary(runDate)
	c.JSON(http.StatusOK, items)
}

func GetReserveSummaryForProduct(c *gin.Context) {
	runDate := c.Param("run_date")
	prodCode := c.Param("product")
	items := services.GetReserveSummaryForProduct(runDate, prodCode)

	c.JSON(http.StatusOK, items)
}

func GetReserveSummaryForPortfolio(c *gin.Context) {
	runDate := c.Param("run_date")
	portfolio := c.Param("portfolio")
	items := services.GetReserveSummaryForPortfolio(runDate, portfolio)

	c.JSON(http.StatusOK, items)
}

func GetReserveSummaryForPortfolioProduct(c *gin.Context) {
	runDate := c.Param("run_date")
	portfolio := c.Param("portfolio")
	prodCode := c.Param("product")
	items := services.GetReserveSummaryForPortfolioProduct(runDate, portfolio, prodCode)

	c.JSON(http.StatusOK, items)
}

func GetReserveSummaryForPortfolioProductGroup(c *gin.Context) {
	runDate := c.Param("run_date")
	portfolio := c.Param("portfolio")
	prodCode := c.Param("product")
	group := c.Param("group")

	items := services.GetReserveSummaryForPortfolioProductGroup(runDate, portfolio, prodCode, group)

	c.JSON(http.StatusOK, items)
}
func GetReserveSummaryForProductAndBasis(c *gin.Context) {
	runDate := c.Param("run_date")
	prodCode := c.Param("product")
	basis := c.Param("basis")
	//spCode, err := strconv.Atoi(c.Param("spcode"))
	items := services.GetReserveSummaryForProductAndBasis(runDate, prodCode, basis)

	c.JSON(http.StatusOK, items)
}

func GetReserveSummaryForProductBasisAndSpcode(c *gin.Context) {
	runDate := c.Param("run_date")
	prodCode := c.Param("product")
	spCode, err := strconv.Atoi(c.Param("spcode"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "A valid spcode is required"})
		return
	}
	basis := c.Param("basis")
	items := services.GetReserveSummaryForProductBasisAndSpcode(runDate, prodCode, spCode, basis)

	c.JSON(http.StatusOK, items)
}

func GetReserveSummaryForBasis(c *gin.Context) {
	runDate := c.Param("run_date")
	basis := c.Param("basis")
	items := services.GetReserveSummaryForBasis(runDate, basis)
	c.JSON(http.StatusOK, items)
}

func GetReserveSummaryForBasisAndProduct(c *gin.Context) {
	runDate := c.Param("run_date")
	basis := c.Param("basis")
	product := c.Param("product")
	items := services.GetReserveSummaryForBasisAndProduct(runDate, basis, product)
	c.JSON(http.StatusOK, items)
}

func GetReserveSummaryForBasisProductAndSpcode(c *gin.Context) {
	runDate := c.Param("run_date")
	basis := c.Param("basis")
	product := c.Param("product")
	spCode, err := strconv.Atoi(c.Param("spcode"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "A valid spcode is required"})
		return
	}
	items := services.GetReserveSummaryForBasisProductAndSpcode(runDate, basis, product, spCode)
	c.JSON(http.StatusOK, items)
}

func GetReserveSummaryRunlist(c *gin.Context) {
	results := services.GetReserveSummaryRunlist()
	c.JSON(http.StatusOK, results)
}

func GetPAAReserveSummaryRunlist(c *gin.Context) {
	results := services.GetPAAReserveSummaryRunlist()
	c.JSON(http.StatusOK, results)
}

// GetParameterMeta returns distinct bases and years available for a product's parameters.
func GetParameterMeta(c *gin.Context) {
	prodId, err := strconv.Atoi(c.Param("prodId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "A valid product id is required"})
		return
	}
	prod, err := services.GetProductById(prodId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Product not found"})
		return
	}
	c.JSON(http.StatusOK, services.GetParameterMeta(prod.ProductCode))
}

// GetParametersForYearsAndBasis returns parameter disclosure for user-chosen years and basis.
func GetParametersForYearsAndBasis(c *gin.Context) {
	prodId, err := strconv.Atoi(c.Param("prodId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "A valid product id is required"})
		return
	}
	prod, err := services.GetProductById(prodId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Product not found"})
		return
	}
	year1, err := strconv.Atoi(c.Query("year1"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "year1 must be a valid integer"})
		return
	}
	year2, err := strconv.Atoi(c.Query("year2"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "year2 must be a valid integer"})
		return
	}
	basis := c.Query("basis")
	items := services.GetParameterDisclosure(prod.ProductCode, year1, year2, basis)
	c.JSON(http.StatusOK, gin.H{"parameters": items})
}

// GetMortalityForYears returns mortality data for two user-chosen years.
func GetMortalityForYears(c *gin.Context) {
	prodId, err := strconv.Atoi(c.Param("prodId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "A valid product id is required"})
		return
	}
	prod, err := services.GetProductById(prodId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Product not found"})
		return
	}
	year1, err := strconv.Atoi(c.Query("year1"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "year1 must be a valid integer"})
		return
	}
	year2, err := strconv.Atoi(c.Query("year2"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "year2 must be a valid integer"})
		return
	}
	c.JSON(http.StatusOK, services.GetMortalityForYears(prod.ProductCode, year1, year2))
}

// GetAllYieldCurveCodes returns every distinct yield_curve_code in the database.
func GetAllYieldCurveCodes(c *gin.Context) {
	codes := services.GetAllYieldCurveCodes()
	c.JSON(http.StatusOK, gin.H{"codes": codes})
}

// GetYieldCurveDisclosureByCodes returns a side-by-side comparison of two yield curve codes.
func GetYieldCurveDisclosureByCodes(c *gin.Context) {
	code1 := c.Query("code1")
	code2 := c.Query("code2")
	if code1 == "" || code2 == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "code1 and code2 query parameters are required"})
		return
	}
	items := services.GetYieldCurveDisclosureByCodes(code1, code2)
	c.JSON(http.StatusOK, gin.H{"items": items})
}
