package controllers

import (
	"api/models"
	"api/services"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateProductFamilies(c *gin.Context) {
	var productFamilies []models.ProductFamily
	err := c.Bind(&productFamilies)
	if err != nil {
		fmt.Println(err)
	}

	services.PopulateProductFamilies(productFamilies)

}

func GetAppVersion(c *gin.Context) {
	appVersion := services.GetAppVersion()
	c.JSON(http.StatusOK, gin.H{"version": appVersion})
}

func LoadDummyData(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Output Results"})
}

func CheckHealth(c *gin.Context) {
	appVersion := services.GetAppVersion()
	c.JSON(http.StatusOK, gin.H{"status": "healthy", "version": appVersion})
}
