package websocket

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"social-network/pkg/models"

	"github.com/gorilla/websocket"
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
	Type       string `json:"type"` // "private", "group"
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
				h.sendToUser(*message.ReceiverID, message)
				h.sendToUser(message.SenderID, message) // Echo back to sender
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
		default:
			h.mu.Lock()
			close(client.send)
			delete(h.clients, userID)
			h.mu.Unlock()
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
	}()

	c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		_, messageData, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		var msg Message
		if err := json.Unmarshal(messageData, &msg); err != nil {
			log.Printf("error unmarshaling message: %v", err)
			continue
		}

		msg.SenderID = c.userID
		msg.Username = c.username
		msg.Timestamp = time.Now().Format(time.RFC3339)

		c.hub.broadcast <- &msg
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(54 * time.Second)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued messages
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
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
		send:     make(chan []byte, 256),
		userID:   user.ID,
		username: user.FirstName + " " + user.LastName,
	}

	client.hub.register <- client

	go client.writePump()
	go client.readPump()
}
