package controllers

import (
    "api/models"
    "api/services"
    "net/http"
    "strconv"

    "github.com/gin-gonic/gin"
)

// CreateBordereauxTemplate handles POST to create a new template
func CreateBordereauxTemplate(c *gin.Context) {
    var payload models.BordereauxTemplate
    if err := c.ShouldBindJSON(&payload); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    if payload.Name == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "name is required"})
        return
    }
    if err := services.CreateBordereauxTemplate(&payload); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusCreated, payload)
}

// GetBordereauxTemplates handles GET to list all templates
func GetBordereauxTemplates(c *gin.Context) {
    list, err := services.GetBordereauxTemplates()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, list)
}

// GetBordereauxTemplate handles GET by id
func GetBordereauxTemplate(c *gin.Context) {
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
        return
    }
    tpl, err := services.GetBordereauxTemplateByID(id)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, tpl)
}

// UpdateBordereauxTemplate handles PUT by id
func UpdateBordereauxTemplate(c *gin.Context) {
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
        return
    }
    var payload models.BordereauxTemplate
    if err := c.ShouldBindJSON(&payload); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    updated, err := services.UpdateBordereauxTemplate(id, payload)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, updated)
}

// DeleteBordereauxTemplate handles DELETE by id
func DeleteBordereauxTemplate(c *gin.Context) {
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
        return
    }
    if err := services.DeleteBordereauxTemplate(id); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.Status(http.StatusNoContent)
}
