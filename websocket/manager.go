// File: connection_manager.go
package websocket

import (
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

// Pastikan ini ada di dalam package `websocket`.
type Message struct {
	RoomID  string `json:"roomID"`
	Content string `json:"content"`
	Sender  int    `json:"sender"`
}

type ConnectionManager struct {
	Connections map[string]map[*websocket.Conn]struct{}
	mu          sync.Mutex
}

func NewConnectionManager() *ConnectionManager {
	return &ConnectionManager{
		Connections: make(map[string]map[*websocket.Conn]struct{}),
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

// BroadcastMessage mengirimkan pesan ke semua koneksi di room tertentu.
func (cm *ConnectionManager) BroadcastMessage(chatRoomID string, message []byte) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	// Kirim pesan ke semua koneksi di room tersebut
	for conn := range cm.Connections[chatRoomID] {
		if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
			// Log error atau handle jika diperlukan
			for conn := range cm.Connections[chatRoomID] {
				if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
					log.Printf("Failed to send message to connection in room %s: %v", chatRoomID, err)

					// Opsional: Hapus koneksi yang error dari daftar koneksi.
					// Ini mungkin memerlukan penguncian dan pengecekan ulang karena Anda mengubah map saat iterasi.
					cm.mu.Lock()
					if _, ok := cm.Connections[chatRoomID][conn]; ok {
						delete(cm.Connections[chatRoomID], conn)
						// Jika perlu, lakukan tindakan tambahan seperti menutup koneksi.
						conn.Close()
					}
					cm.mu.Unlock()

					// Anda juga bisa memutuskan untuk melakukan tindakan lain, seperti mencoba mengirim pesan error ke klien,
					// atau melakukan upaya koneksi ulang, tergantung pada kasus penggunaan spesifik Anda.
				}
			}

		}
	}
}
