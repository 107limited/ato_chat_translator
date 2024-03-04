// File: connection_manager.go
package websocket

import (
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// Struktur Message sesuai yang Anda berikan
type Message struct {
	ID                int    `json:"id"`
	JapaneseText      string `json:"japanese_text"`
	EnglishText       string `json:"english_text"`
	Speaker           string `json:"speaker"`
	UserID            int    `json:"user_id"` // Asumsikan ini adalah User1ID dari TranslationRequest
	CompanyID         int    `json:"company_id"`
	ChatRoomID        int    `json:"chat_room_id"`
	OriginalMessage   string `json:"original_message"`
	TranslatedMessage string `json:"translated_message"`
	CreatedAt         string `json:"created_at"`
	Date              int64  `json:"date"`
}

// Messagefmt represents the format of a message being parsed.
type Messagefmt struct {
	Id        string    `json:"id"`
	ToId      string    `json:"to_id"`
	UserId    string    `json:"user_id"`
	Message   string    `json:"message"`
	Date      time.Time `json:"date"`
	RoomId    int       `json:"room_id"`
	CompanyId int       `json:"company_id"`
	CreatedAt string    `json:"created_at"`
	English   string    `json:"english"`
	Japanese  string    `json:"japanese"`
	Speaker   string    `json:"speaker"`
}

type ConnectionManager struct {
    Connections         map[string]map[*websocket.Conn]struct{} // Room-specific connections
    GlobalConnections   map[*websocket.Conn]struct{}            // Global connections
    UserConnections     map[int]*websocket.Conn                 // User ID to connection
    connectionsMu       sync.Mutex                              // Mutex for room-specific connections
    globalConnectionsMu sync.Mutex                              // Mutex for global connections
    userConnectionsMu   sync.Mutex                              // Mutex for user connections
}

// NewConnectionManager initializes and returns a new instance of ConnectionManager.
func NewConnectionManager() *ConnectionManager {
    return &ConnectionManager{
        Connections:       make(map[string]map[*websocket.Conn]struct{}),
        GlobalConnections: make(map[*websocket.Conn]struct{}),
        UserConnections:   make(map[int]*websocket.Conn),
    }
}


// AddConnection adds a new connection to a specific chat room.
func (cm *ConnectionManager) AddConnection(chatRoomID string, conn *websocket.Conn) {
	cm.connectionsMu.Lock()
	defer cm.connectionsMu.Unlock()

	if _, ok := cm.Connections[chatRoomID]; !ok {
		cm.Connections[chatRoomID] = make(map[*websocket.Conn]struct{})
	}
	cm.Connections[chatRoomID][conn] = struct{}{}
}

// RemoveConnection removes a connection from a specific chat room.
func (cm *ConnectionManager) RemoveConnection(chatRoomID string, conn *websocket.Conn) {
	cm.connectionsMu.Lock()
	defer cm.connectionsMu.Unlock()

	if connections, ok := cm.Connections[chatRoomID]; ok {
		delete(connections, conn)
		if len(connections) == 0 {
			delete(cm.Connections, chatRoomID)
		}
	}
	conn.Close()
}

// AddGlobalConnection adds a new connection to the global list.
func (cm *ConnectionManager) AddGlobalConnection(conn *websocket.Conn) {
	cm.globalConnectionsMu.Lock()
	defer cm.globalConnectionsMu.Unlock()
	cm.GlobalConnections[conn] = struct{}{}
}

// RemoveGlobalConnection removes a connection from the global list.
func (cm *ConnectionManager) RemoveGlobalConnection(conn *websocket.Conn) {
	cm.globalConnectionsMu.Lock()
	defer cm.globalConnectionsMu.Unlock()
	delete(cm.GlobalConnections, conn)
	conn.Close()
}

// BroadcastMessage sends a message to all connections within a specific chat room.
func (cm *ConnectionManager) BroadcastMessage(chatRoomID string, message []byte) {
	cm.connectionsMu.Lock()
	connections, ok := cm.Connections[chatRoomID]
	cm.connectionsMu.Unlock()

	if !ok {
		return // Chat room does not exist or has no connections
	}

	for conn := range connections {
		if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
			log.Printf("Failed to send message to connection in room %s: %v", chatRoomID, err)
			// Remove the failing connection
			cm.RemoveConnection(chatRoomID, conn)
		}
	}
}

func (cm *ConnectionManager) SetUserOnline(userID int, conn *websocket.Conn) {
    cm.userConnectionsMu.Lock()
    cm.UserConnections[userID] = conn
    cm.userConnectionsMu.Unlock()
    cm.broadcastOnlineStatus(userID, true)
}

func (cm *ConnectionManager) SetUserOffline(userID int) {
    cm.userConnectionsMu.Lock()
    delete(cm.UserConnections, userID)
    cm.userConnectionsMu.Unlock()
    cm.broadcastOnlineStatus(userID, false)
}

func (cm *ConnectionManager) broadcastOnlineStatus(userID int, isOnline bool) {
    status := OnlineStatus{
        UserID:   userID,
        IsOnline: isOnline,
    }
    statusMsg, err := json.Marshal(status)
    if err != nil {
        log.Printf("Error marshaling online status: %v", err)
        return
    }
    cm.globalConnectionsMu.Lock()
    defer cm.globalConnectionsMu.Unlock()
    for conn := range cm.GlobalConnections {
        if err := conn.WriteMessage(websocket.TextMessage, statusMsg); err != nil {
            log.Printf("Error sending online status to global connections: %v", err)
        }
    }
}
