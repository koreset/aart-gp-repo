package controllers

import (
    "api/log"
    "api/models"
    "api/services"
    "context"
    "net/http"
    "strconv"

    "github.com/gin-gonic/gin"
)

func CreateSchemeType(c *gin.Context) {
    requestID, exists := c.Get("requestID")
    var ctx context.Context
    if exists {
        ctx = context.WithValue(context.Background(), log.RequestIDKey, requestID.(string))
    } else {
        ctx = context.Background()
    }
    userEmail, emailExists := c.Get("userEmail")
    userName, nameExists := c.Get("userName")
    if emailExists && nameExists {
        ctx = log.ContextWithUserInfo(ctx, userEmail.(string), userName.(string))
    }
    logger := log.WithContext(ctx)

    var payload models.SchemeType
    if err := c.ShouldBindJSON(&payload); err != nil {
        logger.WithField("error", err.Error()).Error("invalid payload")
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    st, err := services.CreateSchemeType(&payload)
    if err != nil {
        logger.WithField("error", err.Error()).Error("failed to create scheme type")
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusCreated, st)
}

func GetSchemeTypes(c *gin.Context) {
    list, err := services.GetSchemeTypes()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, list)
}

func GetSchemeType(c *gin.Context) {
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
        return
    }
    st, err := services.GetSchemeType(id)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, st)
}

func UpdateSchemeType(c *gin.Context) {
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
        return
    }
    var payload models.SchemeType
    if err := c.ShouldBindJSON(&payload); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    st, err := services.UpdateSchemeType(id, &payload)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, st)
}

func DeleteSchemeType(c *gin.Context) {
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
        return
    }
    if err := services.DeleteSchemeType(id); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}
