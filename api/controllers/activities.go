package controllers

import (
	"api/models"
	"api/services"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"io"
	"net/http"
	"time"
)

func CreateActivity(c *gin.Context) {
	reqbody, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Error().Err(err).Send()
	}
	var activity models.Activity
	err = json.Unmarshal(reqbody, &activity)
	if err != nil {
		fmt.Println(err)
	}
	activity.UserEmail = c.Keys["userEmail"].(string)
	activity.UserName = c.Keys["userName"].(string)
	activity.Date = time.Now()
	fmt.Println(activity)
	err = services.StoreActivity(activity)
	c.JSON(http.StatusOK, err)
}
