package models

import "time"

// ─── Notification ───────────────────────────────────────────────────────────

// Notification represents an in-app notification sent to a user.
// Type values: quote_submitted, quote_approved, quote_rejected, quote_accepted,
// message_received, mention, scheme_status_change,
// schedule_reviewed, schedule_approved, schedule_finalized, schedule_voided, schedule_cancelled,
// scheme_suspended, scheme_reinstated,
// submission_reviewed, submission_query_raised, submission_accepted, submission_rejected,
// bordereaux_reviewed, bordereaux_approved, bordereaux_submitted,
// ri_bordereaux_submitted, ri_bordereaux_acknowledged,
// csm_run_reviewed, csm_run_approved, amendment_approved,
// settlement_updated, settlement_dispute_resolved, claim_payment_summary,
// custom_tir_requested, custom_tir_uploaded
type Notification struct {
	ID             int        `gorm:"primaryKey;autoIncrement" json:"id"`
	RecipientEmail string     `json:"recipient_email" gorm:"not null;index"`
	SenderEmail    string     `json:"sender_email" gorm:"default:''"`
	SenderName     string     `json:"sender_name" gorm:"default:''"`
	Type           string     `json:"type" gorm:"not null;index"`
	Title          string     `json:"title" gorm:"not null"`
	Body           string     `json:"body" gorm:"type:text"`
	ObjectType     string     `json:"object_type" gorm:"default:''"`
	ObjectID       int        `json:"object_id" gorm:"default:0"`
	IsRead         bool       `json:"is_read" gorm:"default:false;index"`
	ReadAt         *time.Time `json:"read_at"`
	CreatedAt      time.Time  `json:"created_at" gorm:"autoCreateTime"`
}

// CreateNotificationRequest is the payload for creating a notification.
type CreateNotificationRequest struct {
	RecipientEmail string `json:"recipient_email" binding:"required"`
	SenderEmail    string `json:"sender_email"`
	SenderName     string `json:"sender_name"`
	Type           string `json:"type" binding:"required"`
	Title          string `json:"title" binding:"required"`
	Body           string `json:"body"`
	ObjectType     string `json:"object_type"`
	ObjectID       int    `json:"object_id"`
}

// NotificationFilter holds query parameters for listing notifications.
type NotificationFilter struct {
	IsRead   *bool  `form:"is_read"`
	Type     string `form:"type"`
	Page     int    `form:"page"`
	PageSize int    `form:"page_size"`
}

// NotificationListResponse wraps a paginated list of notifications.
type NotificationListResponse struct {
	Notifications []Notification `json:"notifications"`
	Total         int64          `json:"total"`
	Page          int            `json:"page"`
	PageSize      int            `json:"page_size"`
}

// ─── Conversation ───────────────────────────────────────────────────────────

// Conversation represents a threaded discussion tied to a business object.
type Conversation struct {
	ID           int                       `gorm:"primaryKey;autoIncrement" json:"id"`
	ObjectType   string                    `json:"object_type" gorm:"index"`
	ObjectID     int                       `json:"object_id" gorm:"index"`
	Title        string                    `json:"title" gorm:"default:''"`
	CreatedBy    string                    `json:"created_by" gorm:"not null"`
	IsArchived   bool                      `json:"is_archived" gorm:"default:false"`
	CreatedAt    time.Time                 `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time                 `json:"updated_at" gorm:"autoUpdateTime"`
	Participants []ConversationParticipant `json:"participants" gorm:"foreignKey:ConversationID"`
	Messages     []ConversationMessage     `json:"messages,omitempty" gorm:"foreignKey:ConversationID"`
}

// ConversationParticipant links a user to a conversation.
type ConversationParticipant struct {
	ID             int        `gorm:"primaryKey;autoIncrement" json:"id"`
	ConversationID int        `json:"conversation_id" gorm:"index;not null"`
	UserEmail      string     `json:"user_email" gorm:"index;not null"`
	UserName       string     `json:"user_name" gorm:"default:''"`
	JoinedAt       time.Time  `json:"joined_at" gorm:"autoCreateTime"`
	LastReadAt     *time.Time `json:"last_read_at"`
}

// ConversationMessage is a single message within a conversation.
type ConversationMessage struct {
	ID              int       `gorm:"primaryKey;autoIncrement" json:"id"`
	ConversationID  int       `json:"conversation_id" gorm:"index;not null"`
	SenderEmail     string    `json:"sender_email" gorm:"not null"`
	SenderName      string    `json:"sender_name" gorm:"default:''"`
	Body            string    `json:"body" gorm:"type:text;not null"`
	MessageType     string    `json:"message_type" gorm:"default:'text'"` // text, system, action
	ParentMessageID *int      `json:"parent_message_id"`
	IsEdited        bool      `json:"is_edited" gorm:"default:false"`
	CreatedAt       time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt       time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// CreateConversationRequest is the payload for starting a new conversation.
type CreateConversationRequest struct {
	ObjectType        string   `json:"object_type" binding:"required"`
	ObjectID          int      `json:"object_id"`
	Title             string   `json:"title"`
	ParticipantEmails []string `json:"participant_emails" binding:"required"`
}

// SendMessageRequest is the payload for posting a message to a conversation.
type SendMessageRequest struct {
	Body            string `json:"body" binding:"required"`
	MessageType     string `json:"message_type"`
	ParentMessageID *int   `json:"parent_message_id"`
}

// AddParticipantRequest is the payload for adding a user to a conversation.
type AddParticipantRequest struct {
	UserEmail string `json:"user_email" binding:"required"`
	UserName  string `json:"user_name"`
}
