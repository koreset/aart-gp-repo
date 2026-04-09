package controllers

import (
	"api/models"
	"api/services"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
)

func PostModelPointVariable(c *gin.Context) {
	var mpv []models.ModelPointVariable

	c.Bind(&mpv)
	fmt.Println(mpv)
	//Do validation here...

	for _, mp := range mpv {
		err := services.DB.Create(&mp).Error
		if err != nil {
			fmt.Println(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}
	}

	c.JSON(http.StatusOK, mpv)
}

// GetAllModelPointVariables godoc
// @Summary Get a list of available model point variables
// @Description get all model point variables
// @ID getAllModelPointvariables
// @Accept json
// @Produce json
// @Success 200 {array} models.ModelPointVariable
// @Router /modelpointvariables [get]
func GetAllModelPointVariables(c *gin.Context) {
	var mpvs []models.ModelPointVariable
	err := services.DB.Find(&mpvs).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
	}
	c.JSON(http.StatusOK, mpvs)
}

// GetAllModelPointVariablesForProduct godoc
// @Summary Get a list of model point variables for product
// @Description get all applicable model point variables for a product
// @Param prodId path string true "Product Code"
// @Accept json
// @Produce json
// @Success 200 {array} models.ModelPointVariable
// @Router /modelpointvariables/{prodId} [get]
func GetModelPointVariablesForProduct(c *gin.Context) {
	prodId := c.Param("prodId")
	if len(strings.TrimSpace(prodId)) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Supply a product code"})
	}
	var mpvs []models.ProductModelpointVariable

	services.DB.Where("product_id=?", prodId).Find(&mpvs)
	fmt.Println(mpvs)
	c.JSON(http.StatusOK, mpvs)
}

func GetModelPointStructureForProduct(c *gin.Context) {
	prodcode := c.Param("prodId")
	if len(strings.TrimSpace(prodcode)) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Supply a product code"})
	}

	tableName := strings.ToLower(prodcode) + "_modelpoints"

	mps := services.GetTableStructure(tableName)
	c.JSON(http.StatusOK, mps)

}

func GetTableStructureForAssociatedTable(c *gin.Context) {
	tableId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "a valid product table id is required"})
	}
	csvArray := services.GetAssociatedTableStructure(tableId)

	c.JSON(http.StatusOK, csvArray)

}

func GetModelPointCountForProduct(c *gin.Context) {
	prodId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	mpCount, err := services.GetModelPointCountForProduct(prodId)

	c.JSON(http.StatusOK, gin.H{"results": mpCount})
}

func GetModelPointsForProduct(c *gin.Context) {
	prodId, err := strconv.Atoi(c.Param("id"))
	year, err := strconv.Atoi(c.Param("year"))
	version := c.Param("version")
	mps, err := services.GetModelPointsForProductAndYear(prodId, year, version, 5000)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, mps)
}

func GetModelPointsForProductExcel(c *gin.Context) {
	prodId, err := strconv.Atoi(c.Param("id"))
	year, err := strconv.Atoi(c.Param("year"))
	version := c.Param("version")
	excelData, err := services.GetModelPointsForProductAndYearExcel(prodId, year, version)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to zip files: %v", err)})
		return
	}

	c.Header("Content-Disposition", `attachment; filename="export.xlsx"`)
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", excelData)

}

func DeleteModelPointsForProduct(c *gin.Context) {
	prodId, _ := strconv.Atoi(c.Param("id"))
	year, _ := strconv.Atoi(c.Param("year"))
	version := c.Param("version")
	err := services.DeleteModelPointsForProductAndYear(prodId, year, version)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"results": "success"})
}

func GetModelPointStats(c *gin.Context) {
	prodId, _ := strconv.Atoi(c.Param("id"))
	runId, _ := strconv.Atoi(c.Param("runId"))
	modelPointStats := services.GetModelPointVariableStatsForProduct(prodId, runId)
	c.JSON(http.StatusOK, modelPointStats)
}

func GetModelPointVersionsForYear(c *gin.Context) {
	prodId, _ := strconv.Atoi(c.Param("id"))
	year, _ := strconv.Atoi(c.Param("year"))
	if year == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "year is required"})
		return
	}
	modelPointVersions, err := services.GetModelPointVersionsForProductAndYear(prodId, year)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, modelPointVersions)
}

func GetModelPointStatsForYear(c *gin.Context) {
	prodId, _ := strconv.Atoi(c.Param("id"))
	year, _ := strconv.Atoi(c.Param("year"))
	version := c.Param("version")
	if year == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "year is required"})
		return
	}
	modelPointStats := services.GetModelPointVariableStatsForProductAndYear(prodId, year, version)
	c.JSON(http.StatusOK, modelPointStats)
}

func GetPaaModelPointStats(c *gin.Context) {
	portfolioName := c.Param("portfolio")
	year, _ := strconv.Atoi(c.Param("year"))
	version := c.Param("version")
	modelPointStats := services.GetModelPointVariableStatsForPortfolio(portfolioName, year, version)
	c.JSON(http.StatusOK, modelPointStats)
}

func PostModelPointVariablesForProduct(c *gin.Context) {
	var mpvs []models.ProductModelpointVariable

	c.Bind(&mpvs)
	fmt.Println(mpvs)
	//Do validation here...

	for _, mpv := range mpvs {
		err := services.DB.Create(&mpv).Error
		if err != nil {
			fmt.Println(err)
			//c.JSON(http.StatusBadRequest, gin.H{"error": err})
		}
	}
	c.JSON(http.StatusOK, mpvs)
}

// GetAssumptionVariables godoc
// @Description This will return all assumption variables available for product configuration
// @Summary Get all available assumption variables
// @Tags products
// @Accept json
// @Produce json
// @Success 200 {array} models.BaseAssumptionVariable
// @Router /products/assumptions/variables [get]
func GetAssumptionVariables(c *gin.Context) {
	assumptionVars := services.GetAssumptionVariables()
	c.JSON(http.StatusOK, assumptionVars)
}

// GetBaseFeatures godoc
// @Description This will return all benefit features available for product configuration
// @Summary Get all available benefit features
// @Tags products
// @Accept json
// @Produce json
// @Success 200 {array} models.BaseFeature
// @Router /products/features [get]
func GetBenefitFeatures(c *gin.Context) {
	benefitFeatures, _ := services.GetBaseFeatures()
	c.JSON(http.StatusOK, benefitFeatures)
}

func GetBasisForValuations(c *gin.Context) {
	productCode := c.Param("id")
	year := c.Param("year")
	basisForValuations, _ := services.GetBasisForValuations(productCode, year)
	c.JSON(http.StatusOK, basisForValuations)
}

func GetModelPointVersions(c *gin.Context) {
	productId, _ := strconv.Atoi(c.Param("id"))
	year, _ := strconv.Atoi(c.Param("year"))
	modelPointVersions, err := services.GetModelPointVersions(productId, year)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, modelPointVersions)
}
