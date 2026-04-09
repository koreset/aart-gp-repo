package controllers

import (
	"api/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetWinners(c *gin.Context) {
	var body map[string]interface{}
	err := c.BindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	winners, err := services.GetOlympicWinners(body)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, winners)
}
