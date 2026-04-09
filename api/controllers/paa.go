package controllers

import (
	"api/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetPAAPortfolioNames(c *gin.Context) {
	list, err := services.GetPAAPortfolioNames()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}

func GetPAARuns(c *gin.Context) {
	list, err := services.GetPAARuns()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}
