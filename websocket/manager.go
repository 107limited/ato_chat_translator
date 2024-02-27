// File: connection_manager.go
package websocket

import (
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

type Message struct {
	ID                int    `json:"id"`
	JapaneseText      string `json:"japanese_text"`
	EnglishText       string `json:"english_text"`
	Speaker           string `json:"speaker"`
	UserID            int    `json:"user_id"`
	CompanyID         int    `json:"company_id"`
	ChatRoomID        int    `json:"chat_room_id"`
	OriginalMessage   string `json:"original_message"`
	TranslatedMessage string `json:"translated_message"`
	CreatedAt         string `json:"created_at"`
	Date              int64  `json:"date"`
}

type ConnectionManager struct {
	Connections map[string]map[*websocket.Conn]struct{}
	GlobalConnections map[*websocket.Conn]struct{}
	mu          sync.Mutex
}

func NewConnectionManager() *ConnectionManager {
	return &ConnectionManager{
		Connections: make(map[string]map[*websocket.Conn]struct{}),
		GlobalConnections: make(map[*websocket.Conn]struct{}),
	}
}

func (cm *ConnectionManager) AddConnection(chatRoomID string, conn *websocket.Conn) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	if _, ok := cm.Connections[chatRoomID]; !ok {
		cm.Connections[chatRoomID] = make(map[*websocket.Conn]struct{})
	}
	cm.Connections[chatRoomID][conn] = struct{}{}
}

func (cm *ConnectionManager) RemoveConnection(chatRoomID string, conn *websocket.Conn) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	if connections, ok := cm.Connections[chatRoomID]; ok {
		if _, ok := connections[conn]; ok {
			delete(connections, conn)
			if len(connections) == 0 {
				delete(cm.Connections, chatRoomID)
			}
		}
	}
}

// AddGlobalConnection adds a new connection to the global list.
func (cm *ConnectionManager) AddGlobalConnection(conn *websocket.Conn) {
    cm.mu.Lock()
    defer cm.mu.Unlock()
    cm.GlobalConnections[conn] = struct{}{}
}

// RemoveGlobalConnection removes a connection from the global list.
func (cm *ConnectionManager) RemoveGlobalConnection(conn *websocket.Conn) {
    cm.mu.Lock()
    defer cm.mu.Unlock()
    delete(cm.GlobalConnections, conn)
}

// AddConnectionToRoom adds a new connection to a specific room.
func (cm *ConnectionManager) AddConnectionToRoom(chatRoomID string, conn *websocket.Conn) {
    cm.mu.Lock()
    defer cm.mu.Unlock()
    if _, ok := cm.Connections[chatRoomID]; !ok {
        cm.Connections[chatRoomID] = make(map[*websocket.Conn]struct{})
    }
    cm.Connections[chatRoomID][conn] = struct{}{}
}

// RemoveConnectionFromRoom removes a connection from a specific room.
func (cm *ConnectionManager) RemoveConnectionFromRoom(chatRoomID string, conn *websocket.Conn) {
    cm.mu.Lock()
    defer cm.mu.Unlock()
    if connections, ok := cm.Connections[chatRoomID]; ok {
        delete(connections, conn)
        if len(connections) == 0 {
            delete(cm.Connections, chatRoomID)
        }
    }
}

// BroadcastMessage sends a message to all connections in a specific room.
func (cm *ConnectionManager) BroadcastMessage(chatRoomID string, message []byte) {
    log.Printf("called: %v", chatRoomID)
    cm.mu.Lock()
    defer cm.mu.Unlock()

    // Send the message to all connections in the room
    log.Printf("connection: %v", cm.Connections)
    for conn := range cm.Connections[chatRoomID] {
        if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
            // Log the error or handle it as needed
            // Optionally: Remove the connection that errored from the list of connections.
            // This might require locking and re-checking since you're modifying the map during iteration.
            cm.mu.Lock()
            if _, ok := cm.Connections[chatRoomID][conn]; ok {
                delete(cm.Connections[chatRoomID], conn)
                // If needed, perform additional actions such as closing the connection.
                conn.Close()
            }
            cm.mu.Unlock()

            // You could also decide to take other actions, such as trying to send an error message to the client,
            // or attempting a reconnect, depending on your specific use case.
        }
    }
}
