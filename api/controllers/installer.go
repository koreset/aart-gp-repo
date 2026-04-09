package controllers

import (
	"api/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

//func CreateDummyProduct(c *gin.Context) {
//	//services.LoadDummyData()
//	c.JSON(http.StatusOK, gin.H{"message": "Data Loaded"})
//}

func LoadBaseData(c *gin.Context) {
	services.BaseData(true)
	c.JSON(http.StatusOK, gin.H{"message": "Base feature data loaded"})
}
