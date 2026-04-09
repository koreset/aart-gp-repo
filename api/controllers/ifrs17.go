package controllers

import (
	"api/services"
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetAOSVariables(c *gin.Context) {
	variables, err := services.GetAOSVariables(context.Background())
	if err != nil {
		c.JSON(http.StatusNotFound, err)
		return
	}
	c.JSON(http.StatusOK, variables)
}

func GetAosVariableSets(c *gin.Context) {
	configs, err := services.GetAosVariableSets()

	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, configs)
}

func GetFinancialReport(c *gin.Context) {
	prodCode := c.Param("prod_code")
	runDate := c.Param("run_date")
	if prodCode == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No product supplied"})
		return
	}
	report := services.CreateFinancialReport(prodCode,runDate)
	c.JSON(http.StatusOK, report)
}
