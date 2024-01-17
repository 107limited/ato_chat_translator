
package websocket

import (
	"sync"

	"github.com/gorilla/websocket"
)

// ConnectionManager bertanggung jawab atas manajemen koneksi WebSocket.
type ConnectionManager struct {
	Connections map[*websocket.Conn]struct{}
	mu          sync.Mutex
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
