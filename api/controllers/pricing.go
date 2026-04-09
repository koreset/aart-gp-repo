package controllers

import (
	"api/services"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func GetScenarioData(c *gin.Context) {
	scenarioID, err := strconv.Atoi(c.Param("scenario_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	data := services.GetPricingData(scenarioID)
	c.JSON(http.StatusOK, data)
}

func DeleteScenario(c *gin.Context) {
	scenarioID, err := strconv.Atoi(c.Param("scenario_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	err = services.DeletePricingScenario(scenarioID)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, nil)
}

func CheckPricingRunName(c *gin.Context) {
	runName := c.Param("run_name")
	found := services.PricingRunNameDoesNotExist(runName)
	c.JSON(http.StatusOK, found)
}

func GetProductPricingParameters(c *gin.Context) {
	prodCode := c.Param("prod_code")

	pricingParams, bases, shockBases, err := services.GetProductPricingParameters(prodCode)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"pricing_parameters": pricingParams, "bases": bases, "shock_bases": shockBases, "error": nil})
}

func GetPricingParameters(c *gin.Context) {
	prodCode := c.Param("prod_code")
	pricingParams, err := services.GetPricingParameters(prodCode)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"data": pricingParams, "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": pricingParams, "error": nil})
}

func GetPricingDemographics(c *gin.Context) {
	prodCode := c.Param("prod_code")
	pricingParams, err := services.GetPricingDemographicsData(prodCode)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"data": pricingParams, "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": pricingParams, "error": nil})
}

func DeletePricingParameters(c *gin.Context) {
	prodCode := c.Param("prod_code")
	if prodCode == "" {
		c.JSON(http.StatusBadRequest, "No product code provided")
		return
	}
	err := services.DeletePricingParameters(prodCode)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, nil)
}

func DeletePricingDemographics(c *gin.Context) {
	prodCode := c.Param("prod_code")
	if prodCode == "" {
		c.JSON(http.StatusBadRequest, "No product code provided")
		return
	}
	err := services.DeletePricingDemographics(prodCode)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, nil)
}

func UploadPricingParameters(c *gin.Context) {
	form, err := c.MultipartForm()
	values := form.Value
	prodCode := values["product_code"][0]
	file := form.File["file"]

	err = services.ProcessPricingParameters(prodCode, file[0])
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, nil)
}

func UploadPricingDemographics(c *gin.Context) {
	form, err := c.MultipartForm()
	values := form.Value
	prodCode := values["product_code"][0]
	file := form.File["file"]
	fmt.Println(file)

	err = services.ProcessPricingDemographics(prodCode, file[0])
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, nil)
}

func UploadBulkPricingParameters(c *gin.Context) {
	form, err := c.MultipartForm()
	file := form.File["file"]
	fmt.Println(file)

	err = services.ProcessBulkPricingParameters(file[0])
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, nil)
}

func UpdatePricingRatingTable(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	values := form.Value
	tableIdString := values["table_id"][0]
	prodCode := values["product_code"][0]
	var yieldCurveCode string
	if values["yield_curve_code"] != nil {
		yieldCurveCode = values["yield_curve_code"][0]
	}

	file := form.File["file"]

	tableId, err := strconv.Atoi(tableIdString)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	ctx := context.WithValue(c, "keys", c.Keys)

	// We need to separate the pricing table from the product table
	err = services.ProcessPricingTable(ctx, tableId, prodCode, file[0], 0, yieldCurveCode) //year as zero should trigger pricing update.
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, nil)
}

func GetPricingModelPointsAndCount(c *gin.Context) {
	prodCode := c.Param("prod_code")

	mps, versions := services.GetPricingModelPoints(prodCode)

	c.JSON(http.StatusOK, gin.H{"model_point_sets": mps, "versions": versions})
}

func GetPricingModelPointVersions(c *gin.Context) {
	prodCode := c.Param("prod_code")

	versions := services.GetPricingModelPointVersions(prodCode)

	c.JSON(http.StatusOK, gin.H{"versions": versions})
}

func GetPricingModelPointsForVersion(c *gin.Context) {
	prodCode := c.Param("prod_code")
	version := c.Param("version")

	mps, err := services.GetPricingModelPointsForVersion(prodCode, version)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, mps)
}

func GetPricingJobExcel(c *gin.Context) {
	pricingRunId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Something wrong with the ID"})
	}

	mpps, err := services.GetModelPointsForPricing(pricingRunId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	services.MakePricingFile("pricing.xlsx", mpps)
	c.FileAttachment("pricing.xlsx", "job.xlsx")

}

func GetPricingControlExcel(c *gin.Context) {
	pricingRunId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Something wrong with the ID"})
	}

	scenarioId, err := strconv.Atoi(c.Param("scenario_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Something wrong with the ID"})
	}

	productCode := c.Param("product_code")

	points, err := services.GetPricingPoints(pricingRunId, scenarioId, productCode)
	services.MakePricingControlExcelFile("pricing-control.xlsx", points)
	c.FileAttachment("pricing-control.xlsx", "pricing-control.xlsx")

}
