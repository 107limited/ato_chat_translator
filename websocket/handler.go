// File: websocket/handler.go
package websocket

import (
	"ato_chat/chat"
	"ato_chat/models"
	"encoding/json"
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	// Upgrade HTTP request ke koneksi WebSocket
	conn, err := upgrader.Upgrade(w, r, nil) 
	if err != nil {
		log.Error("Upgrade to websocket failed:", err)
		return
	}
	defer conn.Close()

	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Warn("read:", err)
			break
		}
		log.Printf("recv: %s", p)

		// Echo pesan yang diterima kembali ke klien
		err = conn.WriteMessage(messageType, p)
		if err != nil {
			log.Warn("write:", err)
			break
		}
	}
}

// MessageFormat mewakili format pesan yang Anda harapkan melalui WebSocket
type MessageFormat struct {
	RoomID  string `json:"roomID"`
	Message string `json:"message"`
}

type Message struct {
	RoomID  string `json:"roomID"`  // Identifies the chat room
	Content string `json:"content"` // The message content
	Sender  int    `json:"sender"`  // The ID of the user sending the message
	// Add any additional fields as needed
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
	// Save the conversation
	err := cs.repo.SaveConversation(&conv)
	if err != nil {
		return err
	}

	// Siarkan pesan ke room
	roomID := fmt.Sprintf("%d", conv.ChatRoomID)
	messageBytes, _ := json.Marshal(conv)
	cs.cm.BroadcastToRoom(roomID, messageBytes)

	return nil
}

// GetRoomIDFromMessage mengurai pesan dan mengembalikan roomID.
func GetRoomIDFromMessage(p []byte) (string, error) {
	var msg MessageFormat
	err := json.Unmarshal(p, &msg)
	if err != nil {
		return "", err
	}
	return msg.RoomID, nil
}
