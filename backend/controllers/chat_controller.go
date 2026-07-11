package controllers

import (
	"backend/database"
	"backend/models"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Boleh diakses dari origin mana saja (karena kita pakai CORS di gin)
	},
}

// Hub untuk menyimpan koneksi aktif
type ClientManager struct {
	clients map[uint]map[*websocket.Conn]bool
	mu      sync.Mutex
}

var manager = ClientManager{
	clients: make(map[uint]map[*websocket.Conn]bool),
}

// WS struct untuk parsing JSON dari client
type WsMessage struct {
	ReceiverID uint   `json:"receiver_id"`
	Content    string `json:"content"`
}

// Handle WebSocket
func HandleChatWebSocket(c *gin.Context) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID := userIDVal.(uint)

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	defer conn.Close()

	// Daftarkan client
	manager.mu.Lock()
	if manager.clients[userID] == nil {
		manager.clients[userID] = make(map[*websocket.Conn]bool)
	}
	manager.clients[userID][conn] = true
	manager.mu.Unlock()

	defer func() {
		manager.mu.Lock()
		if manager.clients[userID] != nil {
			delete(manager.clients[userID], conn)
			if len(manager.clients[userID]) == 0 {
				delete(manager.clients, userID)
			}
		}
		manager.mu.Unlock()
	}()

	for {
		var msg WsMessage
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Println("Read Error:", err)
			break
		}

		// Handle Keep-Alive Ping
		if msg.Content == "PING" && msg.ReceiverID == 0 {
			continue
		}

		// Simpan ke database
		chatMsg := models.ChatMessage{
			SenderID:   userID,
			ReceiverID: msg.ReceiverID,
			Content:    msg.Content,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}
		database.DB.Create(&chatMsg)

		// Kirim ke penerima jika online
		manager.mu.Lock()
		receiverConns, ok := manager.clients[msg.ReceiverID]
		if ok {
			for rc := range receiverConns {
				err = rc.WriteJSON(chatMsg)
				if err != nil {
					log.Println("Write error:", err)
					rc.Close()
					delete(manager.clients[msg.ReceiverID], rc)
				}
			}
			if len(manager.clients[msg.ReceiverID]) == 0 {
				delete(manager.clients, msg.ReceiverID)
			}
		}

		// Kirim balik ke SEMUA perangkat pengirim sebagai echo
		senderConns, okSender := manager.clients[userID]
		if okSender {
			for sc := range senderConns {
				sc.WriteJSON(chatMsg)
			}
		}
		manager.mu.Unlock()
	}
}

// Get Chat History
func GetChatHistory(c *gin.Context) {
	userIDVal, _ := c.Get("user_id")
	userID := userIDVal.(uint)

	otherUserIDStr := c.Param("userId")
	otherUserID, _ := strconv.Atoi(otherUserIDStr)

	var messages []models.ChatMessage
	database.DB.Where(
		"(sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)",
		userID, otherUserID, otherUserID, userID,
	).Order("created_at asc").Find(&messages)

	// Tandai pesan sudah dibaca
	database.DB.Model(&models.ChatMessage{}).Where("receiver_id = ? AND sender_id = ? AND is_read = ?", userID, otherUserID, false).Update("is_read", true)

	c.JSON(http.StatusOK, messages)
}

// Get Contacts
func GetContacts(c *gin.Context) {
	userIDVal, _ := c.Get("user_id")
	userID := userIDVal.(uint)
	roleVal, _ := c.Get("role")
	role := roleVal.(string)

	var users []models.User
	if role == "student" {
		// Student can see teachers
		database.DB.Select("id, username, role").Where("role = ? OR role = ?", "teacher", "admin").Find(&users)
	} else {
		// Tutor/Admin can see everyone except themselves
		database.DB.Select("id, username, role").Where("id != ?", userID).Find(&users)
	}

	type ContactResponse struct {
		ID          uint   `json:"id"`
		Username    string `json:"username"`
		Role        string `json:"role"`
		UnreadCount int64  `json:"unread_count"`
	}

	var contacts []ContactResponse
	for _, u := range users {
		var count int64
		database.DB.Model(&models.ChatMessage{}).Where("sender_id = ? AND receiver_id = ? AND is_read = ?", u.ID, userID, false).Count(&count)
		
		contacts = append(contacts, ContactResponse{
			ID:          u.ID,
			Username:    u.Username,
			Role:        u.Role,
			UnreadCount: count,
		})
	}

	c.JSON(http.StatusOK, contacts)
}

// Get Unread Count
func GetUnreadCount(c *gin.Context) {
	userIDVal, _ := c.Get("user_id")
	userID := userIDVal.(uint)

	var count int64
	database.DB.Model(&models.ChatMessage{}).Where("receiver_id = ? AND is_read = ?", userID, false).Count(&count)

	c.JSON(http.StatusOK, gin.H{"count": count})
}

// Delete Message
func DeleteMessage(c *gin.Context) {
	userIDVal, _ := c.Get("user_id")
	userID := userIDVal.(uint)

	messageIDStr := c.Param("id")
	messageID, err := strconv.Atoi(messageIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid message ID"})
		return
	}

	var message models.ChatMessage
	if err := database.DB.First(&message, messageID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Message not found"})
		return
	}

	// Ensure the user is the sender
	if message.SenderID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only delete your own messages"})
		return
	}

	// Hard delete from database
	if err := database.DB.Delete(&message).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete message"})
		return
	}

	// Notify receiver via WebSocket if they are online
	manager.mu.Lock()
	deletedMsgPayload := models.ChatMessage{
		ID:         message.ID,
		SenderID:   message.SenderID,
		ReceiverID: message.ReceiverID,
		IsDeleted:  true,
	}

	receiverConns, ok := manager.clients[message.ReceiverID]
	if ok {
		for rc := range receiverConns {
			rc.WriteJSON(deletedMsgPayload)
		}
	}
	
	senderConns, okSender := manager.clients[userID]
	if okSender {
		for sc := range senderConns {
			sc.WriteJSON(deletedMsgPayload)
		}
	}
	manager.mu.Unlock()

	c.JSON(http.StatusOK, gin.H{"message": "Message deleted successfully"})
}

