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
        return true // atau logika yang lebih spesifik
    },
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
        cm: cm,
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


// Adjusted to accept ConversationService
func HandleWebSocket(cs *ConversationService) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        conn, err := upgrader.Upgrade(w, r, nil)
        if err != nil {
            // Log error dan kirim response HTTP di sini
            log.Errorf("Failed to upgrade to websocket: %v, with headers: %v", err, r.Header)
            http.Error(w, "Failed to upgrade to websocket", http.StatusBadRequest) // Ini mungkin sudah cukup
            return // Pastikan tidak ada kode tambahan yang mengirim header setelah baris ini
        }
        defer conn.Close()

        for {
            _, p, err := conn.ReadMessage()
            if err != nil {
                log.Errorf("Error reading websocket message: %v", err)
                break // Keluar dari loop jika ada error
            }

            var conv models.Conversation
            err = json.Unmarshal(p, &conv)
            if err != nil {
                log.Errorf("Error unmarshaling message: %v, with payload: %s", err, string(p))
                continue // Lanjutkan ke pesan berikutnya
            }

            err = cs.SaveAndBroadcast(conv)
            if err != nil {
                log.Errorf("Error saving and broadcasting message: %v, with conversation: %+v", err, conv)
            }
        }
    }
}

