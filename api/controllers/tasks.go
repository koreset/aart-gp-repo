package controllers

import (
	"api/models"
	"api/services"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func GetTasksForUser(c *gin.Context) {
	var user models.AppUser
	err := c.BindJSON(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	tasks := services.GetActiveTasks(user)
	c.JSON(http.StatusOK, tasks)
}

func CreateTask(c *gin.Context) {
	var task models.Task
	err := c.BindJSON(&task)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	services.GenerateTask(task)
	fmt.Println(task)
	c.JSON(http.StatusCreated, nil)
}

func DeleteTask(c *gin.Context) {
	taskId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	services.DeleteTask(taskId)
	c.JSON(http.StatusOK, nil)
}
