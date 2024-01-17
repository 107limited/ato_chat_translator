// File: websocket/handler.go
package websocket

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// HandleWebSocket mengelola permintaan WebSocket.
func HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	// Mengelola koneksi melalui ConnectionManager
	connectionManager := NewConnectionManager()
	connectionManager.AddConnection(conn)
	defer connectionManager.RemoveConnection(conn)

	// Loop terus menerus membaca pesan dari koneksi
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		// Mengirim pesan yang sama ke semua koneksi terdaftar
		for c := range connectionManager.Connections {
			if err := c.WriteMessage(messageType, p); err != nil {
				log.Println(err)
				return
			}
		}
	}
}
