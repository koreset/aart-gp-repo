package services

import (
	"api/models"
	"context"
	"fmt"
	"regexp"
	"time"

	"gorm.io/gorm"
)

// CreateConversation starts a new conversation tied to a business object with the given participants.
func CreateConversation(req models.CreateConversationRequest, user models.AppUser) (models.Conversation, error) {
	conv := models.Conversation{
		ObjectType: req.ObjectType,
		ObjectID:   req.ObjectID,
		Title:      req.Title,
		CreatedBy:  user.UserEmail,
	}

	err := DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&conv).Error; err != nil {
			return err
		}

		// Add creator as participant
		allEmails := append([]string{user.UserEmail}, req.ParticipantEmails...)
		seen := make(map[string]bool)
		for _, email := range allEmails {
			if seen[email] {
				continue
			}
			seen[email] = true
			p := models.ConversationParticipant{
				ConversationID: conv.ID,
				UserEmail:      email,
				UserName:       email, // name resolved from user directory later
			}
			if email == user.UserEmail {
				p.UserName = user.UserName
			}
			if err := tx.Create(&p).Error; err != nil {
				return err
			}
			conv.Participants = append(conv.Participants, p)
		}
		return nil
	})

	return conv, err
}

// GetConversationsForUser returns conversations the user is participating in (inbox).
func GetConversationsForUser(email string, page, pageSize int) ([]models.Conversation, int64, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 50 {
		pageSize = 20
	}

	var conversations []models.Conversation
	var total int64

	err := DBReadWithResilience(context.Background(), func(db *gorm.DB) error {
		// Get conversation IDs for this user
		var convIDs []int
		if err := db.Model(&models.ConversationParticipant{}).
			Where("user_email = ?", email).
			Pluck("conversation_id", &convIDs).Error; err != nil {
			return err
		}
		if len(convIDs) == 0 {
			return nil
		}

		if err := db.Model(&models.Conversation{}).Where("id IN ?", convIDs).Count(&total).Error; err != nil {
			return err
		}

		offset := (page - 1) * pageSize
		return db.Preload("Participants").
			Where("id IN ?", convIDs).
			Order("updated_at DESC").
			Offset(offset).Limit(pageSize).
			Find(&conversations).Error
	})

	return conversations, total, err
}

// GetConversationsForObject returns all conversations tied to a specific business object.
func GetConversationsForObject(objectType string, objectID int) ([]models.Conversation, error) {
	var conversations []models.Conversation
	err := DBReadWithResilience(context.Background(), func(db *gorm.DB) error {
		return db.Preload("Participants").Preload("Messages", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at ASC")
		}).Where("object_type = ? AND object_id = ?", objectType, objectID).
			Find(&conversations).Error
	})
	return conversations, err
}

// GetConversation loads a single conversation with participants and messages.
// Verifies the requesting user is a participant.
func GetConversation(id int, email string) (models.Conversation, error) {
	var conv models.Conversation
	err := DBReadWithResilience(context.Background(), func(db *gorm.DB) error {
		// Verify membership
		var count int64
		if err := db.Model(&models.ConversationParticipant{}).
			Where("conversation_id = ? AND user_email = ?", id, email).
			Count(&count).Error; err != nil {
			return err
		}
		if count == 0 {
			return fmt.Errorf("not a participant of this conversation")
		}

		return db.Preload("Participants").Preload("Messages", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at ASC")
		}).First(&conv, id).Error
	})
	return conv, err
}

// SendConversationMessage posts a message to a conversation and notifies participants.
func SendConversationMessage(conversationID int, req models.SendMessageRequest, user models.AppUser) (models.ConversationMessage, error) {
	msg := models.ConversationMessage{
		ConversationID: conversationID,
		SenderEmail:    user.UserEmail,
		SenderName:     user.UserName,
		Body:           req.Body,
		MessageType:    req.MessageType,
	}
	if msg.MessageType == "" {
		msg.MessageType = "text"
	}
	if req.ParentMessageID != nil {
		msg.ParentMessageID = req.ParentMessageID
	}

	err := DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&msg).Error; err != nil {
			return err
		}
		// Touch conversation updated_at
		return tx.Model(&models.Conversation{}).Where("id = ?", conversationID).
			Update("updated_at", time.Now()).Error
	})
	if err != nil {
		return msg, err
	}

	// Push message via WebSocket to all other participants
	go notifyConversationParticipants(conversationID, user.UserEmail, msg)

	// Handle @mentions
	go handleMentions(conversationID, msg, user)

	return msg, nil
}

// EditConversationMessage edits a message (only by the original sender).
func EditConversationMessage(messageID int, body string, email string) error {
	result := DB.Model(&models.ConversationMessage{}).
		Where("id = ? AND sender_email = ?", messageID, email).
		Updates(map[string]interface{}{"body": body, "is_edited": true})
	if result.RowsAffected == 0 {
		return fmt.Errorf("message not found or not owned by user")
	}
	return result.Error
}

// DeleteConversationMessage deletes a message (only by the original sender).
func DeleteConversationMessage(messageID int, email string) error {
	result := DB.Where("id = ? AND sender_email = ?", messageID, email).
		Delete(&models.ConversationMessage{})
	if result.RowsAffected == 0 {
		return fmt.Errorf("message not found or not owned by user")
	}
	return result.Error
}

// AddConversationParticipant adds a user to a conversation.
func AddConversationParticipant(conversationID int, req models.AddParticipantRequest) error {
	p := models.ConversationParticipant{
		ConversationID: conversationID,
		UserEmail:      req.UserEmail,
		UserName:       req.UserName,
	}
	return DB.Create(&p).Error
}

// UpdateConversationLastRead updates the last_read_at timestamp for a participant.
func UpdateConversationLastRead(conversationID int, email string) error {
	now := time.Now()
	return DB.Model(&models.ConversationParticipant{}).
		Where("conversation_id = ? AND user_email = ?", conversationID, email).
		Update("last_read_at", now).Error
}

// ─── Internal helpers ───────────────────────────────────────────────────────

// notifyConversationParticipants pushes a message via WebSocket and creates notifications for offline users.
func notifyConversationParticipants(conversationID int, senderEmail string, msg models.ConversationMessage) {
	emails, err := getConversationParticipantEmails(conversationID)
	if err != nil {
		return
	}

	hub := GetHub()
	envelope := WSEnvelope{
		Type:    WSConversationMessage,
		Payload: msg,
	}

	for _, email := range emails {
		if email == senderEmail {
			continue
		}
		if hub != nil && hub.IsOnline(email) {
			hub.SendToUser(email, envelope)
		} else {
			// Create a persisted notification for offline users
			CreateNotification(models.CreateNotificationRequest{
				RecipientEmail: email,
				SenderEmail:    senderEmail,
				SenderName:     msg.SenderName,
				Type:           "message_received",
				Title:          "New message",
				Body:           truncate(msg.Body, 100),
				ObjectType:     "conversation",
				ObjectID:       conversationID,
			})
		}
	}
}

var mentionRegex = regexp.MustCompile(`@([\w.+-]+@[\w.-]+\.\w+)`)

// handleMentions extracts @email mentions from a message and creates notifications.
func handleMentions(conversationID int, msg models.ConversationMessage, sender models.AppUser) {
	matches := mentionRegex.FindAllStringSubmatch(msg.Body, -1)
	if len(matches) == 0 {
		return
	}
	seen := make(map[string]bool)
	for _, m := range matches {
		email := m[1]
		if email == sender.UserEmail || seen[email] {
			continue
		}
		seen[email] = true
		CreateNotification(models.CreateNotificationRequest{
			RecipientEmail: email,
			SenderEmail:    sender.UserEmail,
			SenderName:     sender.UserName,
			Type:           "mention",
			Title:          fmt.Sprintf("%s mentioned you", sender.UserName),
			Body:           truncate(msg.Body, 100),
			ObjectType:     "conversation",
			ObjectID:       conversationID,
		})
	}
}

// truncate shortens a string to maxLen, appending "…" if truncated.
func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "…"
}
