package websocket

import (
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

// ConnectionManager bertanggung jawab atas manajemen koneksi WebSocket.
type ConnectionManager struct {
	Connections map[*websocket.Conn]struct{}
	mu          sync.Mutex
	// Map room ID ke daftar koneksi
	Rooms map[string]map[*websocket.Conn]struct{}
}

// NewConnectionManager membuat instance baru dari ConnectionManager.
func NewConnectionManager() *ConnectionManager {
	return &ConnectionManager{
		Connections: make(map[*websocket.Conn]struct{}),
	}
}

// AddConnection menambahkan koneksi ke ConnectionManager.
func (cm *ConnectionManager) AddConnection(conn *websocket.Conn) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.Connections[conn] = struct{}{}
}

// RemoveConnection menghapus koneksi dari ConnectionManager.
func (cm *ConnectionManager) RemoveConnection(conn *websocket.Conn) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	delete(cm.Connections, conn)
}

func (cm *ConnectionManager) BroadcastToRoom(roomID string, message []byte) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	for conn := range cm.Rooms[roomID] {
		if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
			log.Printf("Error broadcasting to room: %v", err)
			return err
		}
	}
	return nil
}

