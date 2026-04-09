package services

import (
	"encoding/json"
	"sync"
	"time"

	"api/log"

	"github.com/gorilla/websocket"
)

// ─── Message types for typed WebSocket channels ─────────────────────────────

// WSMessageType identifies the kind of payload in a WebSocket envelope.
type WSMessageType string

const (
	WSNotification        WSMessageType = "notification"
	WSConversationMessage WSMessageType = "conversation_message"
	WSPresence            WSMessageType = "presence"
	WSTyping              WSMessageType = "typing"
)

// WSEnvelope is the JSON wrapper sent over the wire.
type WSEnvelope struct {
	Type    WSMessageType `json:"type"`
	Payload interface{}   `json:"payload"`
}

// ─── Client ─────────────────────────────────────────────────────────────────

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 4096
	sendBufSize    = 256
)

// wsClient represents a single WebSocket connection for a user.
type wsClient struct {
	hub   *Hub
	email string
	conn  *websocket.Conn
	send  chan []byte
}

// readPump reads messages from the client (typing indicators, acks, etc.)
// and unregisters the client on disconnect.
func (c *wsClient) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			break
		}
		// Parse incoming client messages (e.g. typing indicators)
		var env WSEnvelope
		if err := json.Unmarshal(message, &env); err != nil {
			continue
		}
		c.hub.handleClientMessage(c, env)
	}
}

// writePump sends messages from the send channel to the WebSocket connection.
func (c *wsClient) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// ─── Hub ────────────────────────────────────────────────────────────────────

// Hub manages all active WebSocket clients keyed by user email.
type Hub struct {
	mu         sync.RWMutex
	clients    map[string]map[*wsClient]bool // email → set of clients
	register   chan *wsClient
	unregister chan *wsClient
}

var (
	hubInstance *Hub
	hubOnce    sync.Once
)

// InitWSHub creates and starts the singleton hub.
func InitWSHub() {
	hubOnce.Do(func() {
		hubInstance = &Hub{
			clients:    make(map[string]map[*wsClient]bool),
			register:   make(chan *wsClient),
			unregister: make(chan *wsClient),
		}
		go hubInstance.run()
		log.Info("WebSocket hub initialized")
	})
}

// GetHub returns the singleton hub instance.
func GetHub() *Hub {
	return hubInstance
}

// run is the main event loop for the hub.
func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			if h.clients[client.email] == nil {
				h.clients[client.email] = make(map[*wsClient]bool)
			}
			h.clients[client.email][client] = true
			h.mu.Unlock()
			log.WithField("email", client.email).Info("WebSocket client connected")

		case client := <-h.unregister:
			h.mu.Lock()
			if conns, ok := h.clients[client.email]; ok {
				if _, exists := conns[client]; exists {
					delete(conns, client)
					close(client.send)
					if len(conns) == 0 {
						delete(h.clients, client.email)
					}
				}
			}
			h.mu.Unlock()
			log.WithField("email", client.email).Info("WebSocket client disconnected")
		}
	}
}

// RegisterClient creates a wsClient for the given user and starts its pumps.
func (h *Hub) RegisterClient(email string, conn *websocket.Conn) {
	client := &wsClient{
		hub:   h,
		email: email,
		conn:  conn,
		send:  make(chan []byte, sendBufSize),
	}
	h.register <- client
	go client.writePump()
	go client.readPump()
}

// SendToUser sends a typed envelope to all connections for the given email.
func (h *Hub) SendToUser(email string, envelope WSEnvelope) {
	data, err := json.Marshal(envelope)
	if err != nil {
		log.WithField("error", err).Error("Failed to marshal WebSocket envelope")
		return
	}

	// If Redis is available, publish for multi-instance delivery
	if RedisAvailable() {
		h.publishToRedis(email, data)
		return
	}

	h.deliverLocal(email, data)
}

// BroadcastToUsers sends a typed envelope to multiple users.
func (h *Hub) BroadcastToUsers(emails []string, envelope WSEnvelope) {
	data, err := json.Marshal(envelope)
	if err != nil {
		log.WithField("error", err).Error("Failed to marshal WebSocket envelope")
		return
	}
	for _, email := range emails {
		if RedisAvailable() {
			h.publishToRedis(email, data)
		} else {
			h.deliverLocal(email, data)
		}
	}
}

// IsOnline returns true if the user has at least one active connection.
func (h *Hub) IsOnline(email string) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.clients[email]) > 0
}

// GetOnlineUsers returns the emails of all connected users.
func (h *Hub) GetOnlineUsers() []string {
	h.mu.RLock()
	defer h.mu.RUnlock()
	emails := make([]string, 0, len(h.clients))
	for email := range h.clients {
		emails = append(emails, email)
	}
	return emails
}

// deliverLocal sends raw bytes to all local connections for a user.
func (h *Hub) deliverLocal(email string, data []byte) {
	h.mu.RLock()
	conns := h.clients[email]
	h.mu.RUnlock()
	for client := range conns {
		select {
		case client.send <- data:
		default:
			// client buffer full — disconnect
			go func(c *wsClient) { h.unregister <- c }(client)
		}
	}
}

// handleClientMessage processes messages received from a client (e.g. typing).
func (h *Hub) handleClientMessage(sender *wsClient, env WSEnvelope) {
	switch env.Type {
	case WSTyping:
		// Broadcast typing indicator to other participants
		// Payload is expected to contain a "conversation_id" field
		data, ok := env.Payload.(map[string]interface{})
		if !ok {
			return
		}
		convIDFloat, ok := data["conversation_id"].(float64)
		if !ok {
			return
		}
		convID := int(convIDFloat)

		// Look up conversation participants and forward to them
		go h.forwardTypingIndicator(sender.email, convID, env)
	}
}

// forwardTypingIndicator sends a typing event to all other participants of a conversation.
func (h *Hub) forwardTypingIndicator(senderEmail string, conversationID int, env WSEnvelope) {
	participants, err := getConversationParticipantEmails(conversationID)
	if err != nil {
		return
	}
	data, err := json.Marshal(env)
	if err != nil {
		return
	}
	for _, email := range participants {
		if email == senderEmail {
			continue
		}
		h.deliverLocal(email, data)
	}
}

// getConversationParticipantEmails fetches participant emails for a conversation.
func getConversationParticipantEmails(conversationID int) ([]string, error) {
	var emails []string
	err := DB.Table("conversation_participants").
		Where("conversation_id = ?", conversationID).
		Pluck("user_email", &emails).Error
	return emails, err
}

// ─── Redis Pub/Sub for multi-instance scaling ───────────────────────────────

const redisWSChannel = "ws:messages"

type redisWSMessage struct {
	Email string `json:"email"`
	Data  []byte `json:"data"`
}

func (h *Hub) publishToRedis(email string, data []byte) {
	msg := redisWSMessage{Email: email, Data: data}
	bs, err := json.Marshal(msg)
	if err != nil {
		return
	}
	if redisClient != nil {
		redisClient.Publish(redisCtx, redisWSChannel, bs)
	}
}

// StartRedisWSSubscriber listens for WebSocket messages published via Redis
// and delivers them to local clients. Call after InitWSHub and InitRedis.
func StartRedisWSSubscriber() {
	if !RedisAvailable() || hubInstance == nil {
		return
	}
	go func() {
		sub := redisClient.Subscribe(redisCtx, redisWSChannel)
		ch := sub.Channel()
		for msg := range ch {
			var wsMsg redisWSMessage
			if err := json.Unmarshal([]byte(msg.Payload), &wsMsg); err != nil {
				continue
			}
			hubInstance.deliverLocal(wsMsg.Email, wsMsg.Data)
		}
	}()
	log.Info("Redis WebSocket subscriber started")
}
