package http

import (
	"log"
	"sync"

	"github.com/gofiber/websocket/v2"
	"github.com/google/uuid"
)

// WebSocketHub manages all WebSocket connections
type WebSocketHub struct {
	// Map of team ID to WebSocket connection
	connections map[string]*websocket.Conn
	mu          sync.RWMutex

	// Broadcast channel
	broadcast chan *BroadcastMessage
}

// BroadcastMessage represents a message to broadcast
type BroadcastMessage struct {
	Type      string                 `json:"type"`
	TeamID    uuid.UUID              `json:"team_id,omitempty"`
	MatchID   uuid.UUID              `json:"match_id,omitempty"`
	Data      map[string]interface{} `json:"data,omitempty"`
	Timestamp string                 `json:"timestamp"`
}

// NewWebSocketHub creates a new WebSocket hub
func NewWebSocketHub() *WebSocketHub {
	hub := &WebSocketHub{
		connections: make(map[string]*websocket.Conn),
		broadcast:   make(chan *BroadcastMessage, 256),
	}

	// Start broadcast worker
	go hub.broadcastWorker()

	return hub
}

// Register adds a new WebSocket connection
func (h *WebSocketHub) Register(teamID string, conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.connections[teamID] = conn
	log.Printf("WebSocket registered for team: %s (total: %d)", teamID, len(h.connections))
}

// Unregister removes a WebSocket connection
func (h *WebSocketHub) Unregister(teamID string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if conn, ok := h.connections[teamID]; ok {
		conn.Close()
		delete(h.connections, teamID)
		log.Printf("WebSocket unregistered for team: %s (total: %d)", teamID, len(h.connections))
	}
}

// SendToTeam sends a message to a specific team
func (h *WebSocketHub) SendToTeam(teamID string, message interface{}) error {
	h.mu.RLock()
	conn, ok := h.connections[teamID]
	h.mu.RUnlock()

	if !ok {
		return nil // Team not connected, skip
	}

	return conn.WriteJSON(message)
}

// Broadcast sends a message to the broadcast channel
func (h *WebSocketHub) Broadcast(msg *BroadcastMessage) {
	select {
	case h.broadcast <- msg:
	default:
		log.Println("Broadcast channel full, dropping message")
	}
}

// BroadcastToClient sends a typed message to a specific client
func (h *WebSocketHub) BroadcastToClient(clientID string, message map[string]interface{}) {
	h.SendToTeam(clientID, message)
}

// broadcastWorker processes broadcast messages
func (h *WebSocketHub) broadcastWorker() {
	for msg := range h.broadcast {
		// Send to specific team if TeamID is set
		if msg.TeamID != uuid.Nil {
			if err := h.SendToTeam(msg.TeamID.String(), msg); err != nil {
				log.Printf("Failed to send to team %s: %v", msg.TeamID, err)
			}
		} else {
			// Broadcast to all connections
			h.mu.RLock()
			for teamID, conn := range h.connections {
				if err := conn.WriteJSON(msg); err != nil {
					log.Printf("Failed to broadcast to team %s: %v", teamID, err)
				}
			}
			h.mu.RUnlock()
		}
	}
}

// GetBroadcastChannel returns the broadcast channel
func (h *WebSocketHub) GetBroadcastChannel() chan<- *BroadcastMessage {
	return h.broadcast
}
