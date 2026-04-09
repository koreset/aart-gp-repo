package controllers

import (
	"api/models"
	"api/services"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreateConversation handles POST /conversations
func CreateConversation(c *gin.Context) {
	user := c.MustGet("user").(models.AppUser)
	var req models.CreateConversationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, err)
		return
	}
	conv, err := services.CreateConversation(req, user)
	if err != nil {
		InternalError(c, err)
		return
	}
	Created(c, conv)
}

// GetUserConversations handles GET /conversations (inbox)
func GetUserConversations(c *gin.Context) {
	user := c.MustGet("user").(models.AppUser)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	conversations, total, err := services.GetConversationsForUser(user.UserEmail, page, pageSize)
	if err != nil {
		InternalError(c, err)
		return
	}
	OK(c, gin.H{"conversations": conversations, "total": total, "page": page, "page_size": pageSize})
}

// GetObjectConversations handles GET /conversations/by-object?object_type=quote&object_id=123
func GetObjectConversations(c *gin.Context) {
	objectType := c.Query("object_type")
	objectID, err := strconv.Atoi(c.Query("object_id"))
	if err != nil || objectType == "" {
		BadRequestMsg(c, "object_type and object_id are required")
		return
	}
	conversations, err := services.GetConversationsForObject(objectType, objectID)
	if err != nil {
		InternalError(c, err)
		return
	}
	OK(c, conversations)
}

// GetConversation handles GET /conversations/:id
func GetConversation(c *gin.Context) {
	user := c.MustGet("user").(models.AppUser)
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id == 0 {
		BadRequestMsg(c, "invalid conversation id")
		return
	}
	conv, err := services.GetConversation(id, user.UserEmail)
	if err != nil {
		BadRequest(c, err)
		return
	}
	OK(c, conv)
}

// SendMessage handles POST /conversations/:id/messages
func SendMessage(c *gin.Context) {
	user := c.MustGet("user").(models.AppUser)
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id == 0 {
		BadRequestMsg(c, "invalid conversation id")
		return
	}
	var req models.SendMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, err)
		return
	}
	msg, err := services.SendConversationMessage(id, req, user)
	if err != nil {
		InternalError(c, err)
		return
	}
	Created(c, msg)
}

// EditMessage handles PATCH /conversations/messages/:message_id
func EditMessage(c *gin.Context) {
	user := c.MustGet("user").(models.AppUser)
	messageID, err := strconv.Atoi(c.Param("message_id"))
	if err != nil || messageID == 0 {
		BadRequestMsg(c, "invalid message id")
		return
	}
	var req struct {
		Body string `json:"body" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, err)
		return
	}
	if err := services.EditConversationMessage(messageID, req.Body, user.UserEmail); err != nil {
		BadRequest(c, err)
		return
	}
	OKMsg(c, "message updated")
}

// DeleteMessage handles DELETE /conversations/messages/:message_id
func DeleteMessage(c *gin.Context) {
	user := c.MustGet("user").(models.AppUser)
	messageID, err := strconv.Atoi(c.Param("message_id"))
	if err != nil || messageID == 0 {
		BadRequestMsg(c, "invalid message id")
		return
	}
	if err := services.DeleteConversationMessage(messageID, user.UserEmail); err != nil {
		BadRequest(c, err)
		return
	}
	OKMsg(c, "message deleted")
}

// AddParticipant handles POST /conversations/:id/participants
func AddParticipant(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id == 0 {
		BadRequestMsg(c, "invalid conversation id")
		return
	}
	var req models.AddParticipantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, err)
		return
	}
	if err := services.AddConversationParticipant(id, req); err != nil {
		InternalError(c, err)
		return
	}
	OKMsg(c, "participant added")
}

// MarkConversationRead handles POST /conversations/:id/read
func MarkConversationRead(c *gin.Context) {
	user := c.MustGet("user").(models.AppUser)
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id == 0 {
		BadRequestMsg(c, "invalid conversation id")
		return
	}
	if err := services.UpdateConversationLastRead(id, user.UserEmail); err != nil {
		InternalError(c, err)
		return
	}
	OKMsg(c, "conversation marked as read")
}
