package websocket

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"ecom-golang-clean-architecture/internal/domain/entities"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// Hub maintains the set of active clients and broadcasts messages to the clients
type Hub struct {
	// Registered clients
	clients map[*Client]bool

	// Inbound messages from the clients
	broadcast chan []byte

	// Register requests from the clients
	register chan *Client

	// Unregister requests from clients
	unregister chan *Client

	// User-specific client mapping
	userClients map[uuid.UUID][]*Client

	// Mutex for thread safety
	mu sync.RWMutex

	// Context for graceful shutdown
	ctx    context.Context
	cancel context.CancelFunc
}

// Client is a middleman between the websocket connection and the hub
type Client struct {
	// The websocket connection
	conn *websocket.Conn

	// Buffered channel of outbound messages
	send chan []byte

	// User ID for this client
	userID uuid.UUID

	// Client ID for tracking
	id uuid.UUID

	// Hub reference
	hub *Hub

	// Last activity time
	lastActivity time.Time
}

// NotificationMessage represents a real-time notification message
type NotificationMessage struct {
	Type         string                 `json:"type"`
	Event        string                 `json:"event"`
	Notification *entities.Notification `json:"notification,omitempty"`
	Data         map[string]interface{} `json:"data,omitempty"`
	Timestamp    time.Time              `json:"timestamp"`
}

const (
	// Time allowed to write a message to the peer
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer
	maxMessageSize = 512
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// Allow connections from any origin in development
		// In production, you should check the origin properly
		return true
	},
}

// NewHub creates a new WebSocket hub
func NewHub(ctx context.Context) *Hub {
	hubCtx, cancel := context.WithCancel(ctx)
	return &Hub{
		clients:     make(map[*Client]bool),
		broadcast:   make(chan []byte),
		register:    make(chan *Client),
		unregister:  make(chan *Client),
		userClients: make(map[uuid.UUID][]*Client),
		ctx:         hubCtx,
		cancel:      cancel,
	}
}

// Run starts the hub
func (h *Hub) Run() {
	defer h.cancel()

	// Start cleanup routine
	go h.cleanupRoutine()

	for {
		select {
		case <-h.ctx.Done():
			log.Println("🔌 WebSocket hub shutting down...")
			return

		case client := <-h.register:
			h.registerClient(client)

		case client := <-h.unregister:
			h.unregisterClient(client)

		case message := <-h.broadcast:
			h.broadcastMessage(message)
		}
	}
}

// registerClient registers a new client
func (h *Hub) registerClient(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.clients[client] = true
	if h.userClients[client.userID] == nil {
		h.userClients[client.userID] = make([]*Client, 0)
	}
	h.userClients[client.userID] = append(h.userClients[client.userID], client)

	log.Printf("🔌 Client %s connected for user %s (total: %d)", 
		client.id, client.userID, len(h.clients))

	// Send welcome message
	welcomeMsg := NotificationMessage{
		Type:      "system",
		Event:     "connected",
		Data:      map[string]interface{}{"message": "Connected to notification service"},
		Timestamp: time.Now(),
	}
	client.sendMessage(welcomeMsg)
}

// unregisterClient unregisters a client
func (h *Hub) unregisterClient(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if _, ok := h.clients[client]; ok {
		delete(h.clients, client)
		close(client.send)

		// Remove from user clients
		userClients := h.userClients[client.userID]
		for i, c := range userClients {
			if c.id == client.id {
				h.userClients[client.userID] = append(userClients[:i], userClients[i+1:]...)
				break
			}
		}

		// Clean up empty user client list
		if len(h.userClients[client.userID]) == 0 {
			delete(h.userClients, client.userID)
		}

		log.Printf("🔌 Client %s disconnected for user %s (total: %d)", 
			client.id, client.userID, len(h.clients))
	}
}

// broadcastMessage broadcasts a message to all clients
func (h *Hub) broadcastMessage(message []byte) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	for client := range h.clients {
		select {
		case client.send <- message:
		default:
			close(client.send)
			delete(h.clients, client)
		}
	}
}

// SendToUser sends a notification to a specific user
func (h *Hub) SendToUser(userID uuid.UUID, notification *entities.Notification) {
	h.mu.RLock()
	clients := h.userClients[userID]
	h.mu.RUnlock()

	if len(clients) == 0 {
		log.Printf("📱 No connected clients for user %s", userID)
		return
	}

	message := NotificationMessage{
		Type:         "notification",
		Event:        "new_notification",
		Notification: notification,
		Timestamp:    time.Now(),
	}

	for _, client := range clients {
		client.sendMessage(message)
	}

	log.Printf("📱 Sent real-time notification to %d clients for user %s", len(clients), userID)
}

// SendToAll broadcasts a notification to all connected clients
func (h *Hub) SendToAll(notification *entities.Notification) {
	message := NotificationMessage{
		Type:         "notification",
		Event:        "broadcast_notification",
		Notification: notification,
		Timestamp:    time.Now(),
	}

	messageBytes, err := json.Marshal(message)
	if err != nil {
		log.Printf("❌ Failed to marshal broadcast message: %v", err)
		return
	}

	h.broadcast <- messageBytes
	log.Printf("📢 Broadcast notification to all %d connected clients", len(h.clients))
}

// GetConnectedUsers returns list of connected user IDs
func (h *Hub) GetConnectedUsers() []uuid.UUID {
	h.mu.RLock()
	defer h.mu.RUnlock()

	users := make([]uuid.UUID, 0, len(h.userClients))
	for userID := range h.userClients {
		users = append(users, userID)
	}
	return users
}

// GetStats returns hub statistics
func (h *Hub) GetStats() map[string]interface{} {
	h.mu.RLock()
	defer h.mu.RUnlock()

	return map[string]interface{}{
		"total_clients":    len(h.clients),
		"connected_users":  len(h.userClients),
		"users_with_clients": h.userClients,
	}
}

// cleanupRoutine periodically cleans up inactive clients
func (h *Hub) cleanupRoutine() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-h.ctx.Done():
			return
		case <-ticker.C:
			h.cleanupInactiveClients()
		}
	}
}

// cleanupInactiveClients removes clients that haven't been active
func (h *Hub) cleanupInactiveClients() {
	h.mu.Lock()
	defer h.mu.Unlock()

	cutoff := time.Now().Add(-10 * time.Minute)
	var toRemove []*Client

	for client := range h.clients {
		if client.lastActivity.Before(cutoff) {
			toRemove = append(toRemove, client)
		}
	}

	for _, client := range toRemove {
		log.Printf("🧹 Cleaning up inactive client %s for user %s", client.id, client.userID)
		delete(h.clients, client)
		close(client.send)
	}
}

// sendMessage sends a message to the client
func (c *Client) sendMessage(message NotificationMessage) {
	messageBytes, err := json.Marshal(message)
	if err != nil {
		log.Printf("❌ Failed to marshal message: %v", err)
		return
	}

	select {
	case c.send <- messageBytes:
		c.lastActivity = time.Now()
	default:
		log.Printf("⚠️ Client %s send channel is full, closing connection", c.id)
		c.hub.unregister <- c
	}
}

// readPump pumps messages from the websocket connection to the hub
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		c.lastActivity = time.Now()
		return nil
	})

	for {
		_, _, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("❌ WebSocket error: %v", err)
			}
			break
		}
		c.lastActivity = time.Now()
	}
}

// writePump pumps messages from the hub to the websocket connection
func (c *Client) writePump() {
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

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued chat messages to the current websocket message
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
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

// HandleWebSocket handles websocket requests from the peer
func (h *Hub) HandleWebSocket(c *gin.Context) {
	// Try to get user ID from context (set by auth middleware)
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		// If not in context, try to get from token query parameter
		token := c.Query("token")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized - no token provided"})
			return
		}

		// For now, use admin user ID for testing
		// In production, you should validate the token and extract user ID
		userIDInterface = uuid.MustParse("5aa91738-ea42-4a8e-ae2e-836ec492ad41")
	}

	userID, ok := userIDInterface.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Upgrade HTTP connection to WebSocket
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("❌ Failed to upgrade connection: %v", err)
		return
	}

	// Create new client
	client := &Client{
		conn:         conn,
		send:         make(chan []byte, 256),
		userID:       userID,
		id:           uuid.New(),
		hub:          h,
		lastActivity: time.Now(),
	}

	// Register client with hub
	h.register <- client

	// Start goroutines for reading and writing
	go client.writePump()
	go client.readPump()
}
