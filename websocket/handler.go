// File: handler.go
package websocket

import (
	"ato_chat/chat"
	"ato_chat/models"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type MessageFormat struct {
	RoomID  string `json:"roomID"`
	Message string `json:"message"`
}

type ConversationService struct {
	repo chat.ConversationRepository
	cm   *ConnectionManager
}

func NewConversationService(repo chat.ConversationRepository, cm *ConnectionManager) *ConversationService {
	return &ConversationService{
		repo: repo,
		cm:   cm,
	}
}

func (cs *ConversationService) SaveAndBroadcast(conv models.Conversation) error {
	// Simpan percakapan
	err := cs.repo.SaveConversation(&conv)
	if err != nil {
		return err
	}

	// Siapkan roomID dan pesan untuk broadcast
	roomID := fmt.Sprintf("%d", conv.ChatRoomID)
	messageBytes, err := json.Marshal(conv)
	if err != nil {
		return err // Jika gagal melakukan marshal, return error
	}

	// Broadcast pesan ke semua koneksi dalam room yang sesuai
	cs.cm.BroadcastMessage(roomID, messageBytes) // Gunakan hanya roomID dan messageBytes

	return nil
}

func HandleWebSocket(cs *ConversationService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Upgrade HTTP to WebSocket connection
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("Upgrade to websocket failed:", err)
			return
		}
		defer conn.Close()

		// Safely get chatRoomID from URL query parameters
		chatRoomIDs, ok := r.URL.Query()["room_id"]
		if !ok || len(chatRoomIDs[0]) < 1 {
			log.Println("URL Param 'room_id' is missing")
			return // Optionally send an error message through WebSocket before returning
		}
		chatRoomID := chatRoomIDs[0]

		// Add connection to ConnectionManager
		cs.cm.AddConnection(chatRoomID, conn)
		defer cs.cm.RemoveConnection(chatRoomID, conn)

		// Continuously read messages from WebSocket
		for {
			_, p, err := conn.ReadMessage()
			if err != nil {
				log.Printf("read error: %v", err)
				break // Exit loop on read error (e.g., client disconnect)
			}

			var conv models.Conversation
			err = json.Unmarshal(p, &conv)
			if err != nil {
				log.Printf("Error unmarshaling message: %v, message: %s", err, string(p))
				continue // Skip this message but continue listening for new ones
			}

			// Set ChatRoomID from URL parameter, assuming conv.ChatRoomID is an int
			if chatRoomID, err := strconv.Atoi(chatRoomID); err == nil {
				conv.ChatRoomID = chatRoomID
			} else {
				log.Printf("Invalid chatRoomID: %v", err)
				continue // Skip this message but continue listening for new ones
			}

			// Save and broadcast the message
			err = cs.SaveAndBroadcast(conv)
			if err != nil {
				log.Printf("error saving and broadcasting message: %v", err)
				// Consider how to handle broadcast errors, possibly notify sender
			}
		}
	}
}
