package controllers

import (
	"api/models"
	"api/services"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var wsUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

// HandleWebSocketUpgrade upgrades an HTTP connection to a WebSocket and
// registers the authenticated user with the hub.
func HandleWebSocketUpgrade(c *gin.Context) {
	// First try user from middleware (Authorization header present)
	userVal, exists := c.Get("user")
	if exists {
		user := userVal.(models.AppUser)
		upgradeAndRegister(c, user)
		return
	}

	// Fallback: parse token from query parameter (browser WebSocket can't send headers)
	token := c.Query("token")
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// Strip "Bearer " prefix if present
	token = strings.TrimPrefix(token, "Bearer ")

	claims := jwt.MapClaims{}
	jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
		return nil, nil
	})

	if claims == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return
	}

	activeUser, ok := claims["user"].(map[string]interface{})
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user claim not found"})
		return
	}

	email, _ := activeUser["Email"].(string)
	fullName, _ := activeUser["FullName"].(string)
	if email == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "email not found in token"})
		return
	}

	user := models.AppUser{UserEmail: email, UserName: fullName}
	upgradeAndRegister(c, user)
}

func upgradeAndRegister(c *gin.Context, user models.AppUser) {
	conn, err := wsUpgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	hub := services.GetHub()
	if hub == nil {
		conn.Close()
		return
	}
	hub.RegisterClient(user.UserEmail, conn)
}
