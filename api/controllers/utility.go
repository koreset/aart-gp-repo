package controllers

import (
	"api/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetDataTableYearVersions(c *gin.Context) {
	// Get the year from the URL parameter
	tableType := c.Param("table_type")

	// Get the data table versions for the specified year
	yearVersions, err := services.GetDataTableYearVersions(tableType)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "The requested data table versions were not found"})
		return
	}

	// Return the versions as JSON
	c.JSON(http.StatusOK, yearVersions)
}

func GetDataTableYears(c *gin.Context) {
	// Get the year from the URL parameter
	tableType := c.Param("table_type")

	// Get the data table versions for the specified year
	years, err := services.GetDataTableYears(tableType)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "The requested data table years were not found"})
		return
	}

	// Return the versions as JSON
	c.JSON(http.StatusOK, years)
}
