package controllers

import (
	"api/models"
	"api/services"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetNotifications handles GET /notifications
func GetNotifications(c *gin.Context) {
	user := c.MustGet("user").(models.AppUser)
	var filter models.NotificationFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		BadRequest(c, err)
		return
	}
	resp, err := services.GetNotifications(user.UserEmail, filter)
	if err != nil {
		InternalError(c, err)
		return
	}
	OK(c, resp)
}

// GetUnreadCount handles GET /notifications/unread-count
func GetUnreadCount(c *gin.Context) {
	user := c.MustGet("user").(models.AppUser)
	count, err := services.GetUnreadNotificationCount(user.UserEmail)
	if err != nil {
		InternalError(c, err)
		return
	}
	OK(c, gin.H{"unread_count": count})
}

// MarkNotificationAsRead handles PATCH /notifications/:id/read
func MarkNotificationAsRead(c *gin.Context) {
	user := c.MustGet("user").(models.AppUser)
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id == 0 {
		BadRequestMsg(c, "invalid notification id")
		return
	}
	if err := services.MarkNotificationAsRead(id, user.UserEmail); err != nil {
		InternalError(c, err)
		return
	}
	OKMsg(c, "notification marked as read")
}

// MarkAllNotificationsAsRead handles POST /notifications/read-all
func MarkAllNotificationsAsRead(c *gin.Context) {
	user := c.MustGet("user").(models.AppUser)
	if err := services.MarkAllNotificationsAsRead(user.UserEmail); err != nil {
		InternalError(c, err)
		return
	}
	OKMsg(c, "all notifications marked as read")
}

// DeleteNotification handles DELETE /notifications/:id
func DeleteNotification(c *gin.Context) {
	user := c.MustGet("user").(models.AppUser)
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id == 0 {
		BadRequestMsg(c, "invalid notification id")
		return
	}
	if err := services.DeleteNotification(id, user.UserEmail); err != nil {
		InternalError(c, err)
		return
	}
	OKMsg(c, "notification deleted")
}
