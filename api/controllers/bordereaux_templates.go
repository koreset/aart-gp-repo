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
    user := c.MustGet("user").(models.AppUser)
    var payload models.BordereauxTemplate
    if err := c.ShouldBindJSON(&payload); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    if payload.Name == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "name is required"})
        return
    }
    if err := services.CreateBordereauxTemplate(&payload, user); err != nil {
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
    user := c.MustGet("user").(models.AppUser)
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
    updated, err := services.UpdateBordereauxTemplate(id, payload, user)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, updated)
}

// DeleteBordereauxTemplate handles DELETE by id
func DeleteBordereauxTemplate(c *gin.Context) {
    user := c.MustGet("user").(models.AppUser)
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
        return
    }
    if err := services.DeleteBordereauxTemplate(id, user); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.Status(http.StatusNoContent)
}

// TestBordereauxTemplate handles POST /bordereaux/templates/:id/test — applies
// the template's field mappings to a sample of live snapshot data and returns
// a preview so operators can validate mappings before a full generation.
func TestBordereauxTemplate(c *gin.Context) {
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
        return
    }
    var req services.TestBordereauxTemplateRequest
    // Body is optional — defaults apply when omitted.
    _ = c.ShouldBindJSON(&req)
    result, svcErr := services.TestBordereauxTemplate(id, req)
    if svcErr != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": svcErr.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"success": true, "data": result})
}
