package websocket

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"social-network/pkg/models"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer
	maxMessageSize = 512 * 1024
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins in development
	},
}

type Client struct {
	hub      *Hub
	conn     *websocket.Conn
	send     chan []byte
	userID   int
	username string
}

type Hub struct {
	clients    map[int]*Client
	broadcast  chan *Message
	register   chan *Client
	unregister chan *Client
	mu         sync.RWMutex
	repo       *models.Repository
}

// BroadcastNewPost sends a notification to all connected clients about a new post
func (h *Hub) BroadcastNewPost(postID int, userID int) {
	message := &Message{
		Type:      "new_post",
		SenderID:  userID,
		Content:   "",
		Timestamp: time.Now().Format(time.RFC3339),
	}

	data, _ := json.Marshal(message)
	
	h.mu.RLock()
	defer h.mu.RUnlock()
	
	// Send to all connected clients
	for _, client := range h.clients {
		select {
		case client.send <- data:
		default:
			// Client's send channel is full, skip
		}
	}
}

type Message struct {
	Type       string `json:"type"` // "private", "group", "notification"
	SenderID   int    `json:"sender_id"`
	ReceiverID *int   `json:"receiver_id,omitempty"`
	GroupID    *int   `json:"group_id,omitempty"`
	Content    string `json:"content"`
	Timestamp  string `json:"timestamp"`
	Username   string `json:"username"`
}

func NewHub(repo *models.Repository) *Hub {
	return &Hub{
		clients:    make(map[int]*Client),
		broadcast:  make(chan *Message, 256),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		repo:       repo,
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client.userID] = client
			h.mu.Unlock()
			log.Printf("Client registered: %d", client.userID)

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client.userID]; ok {
				delete(h.clients, client.userID)
				close(client.send)
			}
			h.mu.Unlock()
			log.Printf("Client unregistered: %d", client.userID)

		case message := <-h.broadcast:
			// Save message to database
			h.repo.CreateMessage(message.SenderID, message.ReceiverID, message.GroupID, message.Content)

			// Send to appropriate recipients
			if message.Type == "private" && message.ReceiverID != nil {
				// Send real-time notification to receiver if they're online
				sender, err := h.repo.GetUserByID(message.SenderID)
				if err == nil {
					notifContent := fmt.Sprintf("%s %s sent you a message", sender.FirstName, sender.LastName)
					h.SendNotificationToUser(*message.ReceiverID, "message", notifContent, message.SenderID)
				}
				
				// Only send to receiver, not back to sender (sender handles optimistically)
				h.sendToUser(*message.ReceiverID, message)
			} else if message.Type == "group" && message.GroupID != nil {
				h.sendToGroup(*message.GroupID, message)
			}
		}
	}
}

func (h *Hub) sendToUser(userID int, message *Message) {
	h.mu.RLock()
	client, ok := h.clients[userID]
	h.mu.RUnlock()

	if ok {
		data, _ := json.Marshal(message)
		select {
		case client.send <- data:
			log.Printf("Message sent to user %d", userID)
		default:
			// Channel full, but don't close connection - just log it
			log.Printf("Warning: Send channel full for user %d, message may be delayed", userID)
			// Try one more time with a timeout
			select {
			case client.send <- data:
				log.Printf("Message sent to user %d on retry", userID)
			case <-time.After(2 * time.Second):
				log.Printf("Failed to send message to user %d after timeout", userID)
			}
		}
	} else {
		// User is offline, create a notification
		if message.Type == "private" && message.SenderID != userID {
			sender, err := h.repo.GetUserByID(message.SenderID)
			if err == nil {
				content := fmt.Sprintf("%s %s sent you a message", sender.FirstName, sender.LastName)
				h.repo.CreateNotification(userID, "message", content, &message.SenderID)
			}
		}
	}
}

// SendNotificationToUser sends a real-time notification to a specific user
func (h *Hub) SendNotificationToUser(userID int, notifType, content string, senderID int) {
	h.mu.RLock()
	client, ok := h.clients[userID]
	h.mu.RUnlock()

	if ok {
		notification := &Message{
			Type:      "notification",
			SenderID:  senderID,
			Content:   content,
			Timestamp: time.Now().Format(time.RFC3339),
		}
		data, _ := json.Marshal(notification)
		select {
		case client.send <- data:
		default:
			// Channel full, skip
		}
	}
}

// SendFollowStatusUpdate sends a real-time follow status update to a user
func (h *Hub) SendFollowStatusUpdate(userID int) {
	h.mu.RLock()
	client, ok := h.clients[userID]
	h.mu.RUnlock()

	if ok {
		statusUpdate := &Message{
			Type:      "follow_status_update",
			SenderID:  0,
			Content:   "Follow status updated",
			Timestamp: time.Now().Format(time.RFC3339),
		}
		data, _ := json.Marshal(statusUpdate)
		select {
		case client.send <- data:
		default:
			// Channel full, skip
		}
	}
}

func (h *Hub) sendToGroup(groupID int, message *Message) {
	// Get all group members
	members, err := h.repo.GetGroupMembers(groupID)
	if err != nil {
		return
	}

	data, _ := json.Marshal(message)
	h.mu.RLock()
	defer h.mu.RUnlock()

	for _, member := range members {
		if client, ok := h.clients[member.UserID]; ok {
			select {
			case client.send <- data:
			default:
				close(client.send)
				delete(h.clients, member.UserID)
			}
		}
	}
}

func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
		log.Printf("Client %d disconnected from readPump", c.userID)
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		log.Printf("✅ Received pong from client %d", c.userID)
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, messageData, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("❌ WebSocket error for client %d: %v", c.userID, err)
			} else {
				log.Printf("🔌 Client %d connection closed: %v", c.userID, err)
			}
			break
		}

		// Reset read deadline on any message received
		c.conn.SetReadDeadline(time.Now().Add(pongWait))

		var msg Message
		if err := json.Unmarshal(messageData, &msg); err != nil {
			log.Printf("⚠️ Error unmarshaling message from client %d: %v", c.userID, err)
			continue
		}

		msg.SenderID = c.userID
		msg.Username = c.username
		msg.Timestamp = time.Now().Format(time.RFC3339)

		log.Printf("📨 Broadcasting message from client %d: type=%s", c.userID, msg.Type)
		c.hub.broadcast <- &msg
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
		log.Printf("🔌 Client %d disconnected from writePump", c.userID)
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				log.Printf("❌ Error getting writer for client %d: %v", c.userID, err)
				return
			}
			w.Write(message)

			// Add queued messages to the current websocket message
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				log.Printf("❌ Error closing writer for client %d: %v", c.userID, err)
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Printf("❌ Error sending ping to client %d: %v", c.userID, err)
				return
			}
			log.Printf("🏓 Sent ping to client %d", c.userID)
		}
	}
}

func ServeWs(hub *Hub, repo *models.Repository, w http.ResponseWriter, r *http.Request) {
	// Get user from session
	cookie, err := r.Cookie("session_id")
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	session, err := repo.GetSession(cookie.Value)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	user, err := repo.GetUserByID(session.UserID)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := &Client{
		hub:      hub,
		conn:     conn,
		send:     make(chan []byte, 1024), // Increased buffer size
		userID:   user.ID,
		username: user.FirstName + " " + user.LastName,
	}

	client.hub.register <- client

	go client.writePump()
	go client.readPump()
}
