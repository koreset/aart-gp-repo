package controllers

import (
	"api/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetBalanceSheetDates(c *gin.Context) {
	dates, err := services.GetBalanceSheetDates()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"dates": dates})
}
