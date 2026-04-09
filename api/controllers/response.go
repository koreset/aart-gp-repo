package controllers

import (
	"api/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// OK sends 200 with success=true and optional data payload.
func OK(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, models.PremiumResponse{Success: true, Data: data})
}

// OKMsg sends 200 with success=true and a plain message (no data).
func OKMsg(c *gin.Context, msg string) {
	c.JSON(http.StatusOK, models.PremiumResponse{Success: true, Message: msg})
}

// Created sends 201 with success=true and the created resource.
func Created(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, models.PremiumResponse{Success: true, Data: data})
}

// BadRequest sends 400 with success=false and the error message.
func BadRequest(c *gin.Context, err error) {
	c.JSON(http.StatusBadRequest, models.PremiumResponse{Success: false, Message: err.Error()})
}

// BadRequestMsg sends 400 with success=false and a plain string message.
func BadRequestMsg(c *gin.Context, msg string) {
	c.JSON(http.StatusBadRequest, models.PremiumResponse{Success: false, Message: msg})
}

// NotFound sends 404 with success=false and the given message.
func NotFound(c *gin.Context, msg string) {
	c.JSON(http.StatusNotFound, models.PremiumResponse{Success: false, Message: msg})
}

// InternalError sends 500 with success=false and the error message.
func InternalError(c *gin.Context, err error) {
	c.JSON(http.StatusInternalServerError, models.PremiumResponse{Success: false, Message: err.Error()})
}
