// File: handler.go
package websocket

import (
	"ato_chat/chat"
	"ato_chat/models"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

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
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("Upgrade to websocket failed:", err)
			return
		}
		defer conn.Close()

		chatRoomID := r.URL.Query().Get("room")
		cs.cm.AddConnection(chatRoomID, conn)
		defer cs.cm.RemoveConnection(chatRoomID, conn)

		for {
			_, p, err := conn.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				break
			}
			log.Printf("%v", p)

			var conv models.Conversation
			err = json.Unmarshal(p, &conv)
			if err != nil {
				log.Printf("Error unmarshaling message: %v, message: %s", err, string(p))
				continue // atau handle error sesuai kebutuhan
			}

			err = cs.SaveAndBroadcast(conv)
			if err != nil {
				log.Println("error saving and broadcasting message:", err)
			}
		}
	}
}
