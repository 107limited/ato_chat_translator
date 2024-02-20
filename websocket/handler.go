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

// Adjusted to accept ConversationService
func HandleWebSocket(cs *ConversationService) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        conn, err := upgrader.Upgrade(w, r, nil)
        if err != nil {
            log.Println("failed to upgrade to websocket:", err)
            return
        }
        defer conn.Close()

        for {
            _, p, err := conn.ReadMessage()
            if err != nil {
                log.Println("error reading websocket message:", err)
                break
            }

            // Logic to extract conversation details from message `p` and create a models.Conversation object
            var conv models.Conversation
            // Assuming `p` is JSON that can be unmarshaled into a models.Conversation
            err = json.Unmarshal(p, &conv)
            if err != nil {
                log.Println("error unmarshaling message:", err)
                continue // or handle error differently
            }

            // Use ConversationService to save and broadcast the message
            err = cs.SaveAndBroadcast(conv)
            if err != nil {
                log.Println("error saving and broadcasting message:", err)
                // Decide how to handle the error; continue to read next messages
            }
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
