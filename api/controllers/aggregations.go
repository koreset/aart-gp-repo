package controllers

import (
	"api/services"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func GetAggregationResults(c *gin.Context) {
	runId, err := strconv.Atoi(c.Param("run_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	productCode := c.Param("prod_code")

	spCode := c.Param("sp_code")

	var variables []string
	err = c.ShouldBindJSON(&variables)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	aggregations, err := services.GetAggregations(runId, productCode, spCode, variables)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, aggregations)
}

func GetAggregationSpCodes(c *gin.Context) {
	runId := c.Param("run_id")
	productCode := c.Param("prod_code")
	spCodes := services.GetAggregationSpCodes(runId, productCode)
	c.JSON(http.StatusOK, spCodes)
}

func GetAggregationVariables(c *gin.Context) {
	variables := services.GetAggregationVariables()
	c.JSON(http.StatusOK, variables)
}
