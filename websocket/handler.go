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
	"time"

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

		chatRoomID, ok := r.URL.Query()["room_id"]
		if !ok || len(chatRoomID[0]) < 1 {
			log.Println("URL Param 'room_id' is missing")
			return
		}

		cs.cm.AddConnection(chatRoomID[0], conn)
		defer cs.cm.RemoveConnection(chatRoomID[0], conn)

		for {
			_, p, err := conn.ReadMessage()
			if err != nil {
				log.Printf("read error: %v", err)
				break
			}

			var conv models.Conversation
			err = json.Unmarshal(p, &conv)
			if err != nil {
				log.Printf("Error unmarshaling message: %v, message: %s", err, string(p))
				continue
			}

			conv.ChatRoomID, err = strconv.Atoi(chatRoomID[0])
			if err != nil {
				log.Printf("Invalid chatRoomID: %v", err)
				continue
			}

			err = cs.SaveAndBroadcast(conv)
			if err != nil {
				log.Printf("error saving and broadcasting message: %v", err)
			} else {
				// After successfully saving and broadcasting the message,
				// send a notification to the chat room participants including sender and receiver.
				notifyParticipants(conv, cs)
			}
		}
	}
}

func notifyParticipants(conv models.Conversation, cs *ConversationService) {
	// Construct the message format as requested.
	msg := struct {
		UserID      int    `json:"user_id"`
		CompanyName string `json:"company_name"`
		Name        string `json:"name"`
		ChatRoomID  int    `json:"chat_room_id"`
		CreatedAt   string `json:"created_at"`
		LastMessage struct {
			English  string `json:"english"`
			Japanese string `json:"japanese"`
			UserID   int    `json:"user_id"`
			Date     int64  `json:"date"`
		} `json:"last_message"`
		LastMessageUser int `json:"last_message_user"`
	}{
		UserID:          conv.UserID,
		CompanyName:     conv.CompanyName,           // Assuming you fetch this from your database or have it in your conversation model
		Name:            conv.Speaker, // Assuming you fetch this from your database based on UserID
		ChatRoomID:      conv.ChatRoomID,
		CreatedAt:       time.Now().UTC().Format(time.RFC3339), // Use the actual creation time of the message
		LastMessageUser: conv.UserID,
	}

	msg.LastMessage.English = conv.EnglishText
	msg.LastMessage.Japanese = conv.JapaneseText
	msg.LastMessage.UserID = conv.UserID
	msg.LastMessage.Date = conv.Date // Ensure this is the timestamp in the correct format

	messageBytes, err := json.Marshal(msg)
	if err != nil {
		log.Printf("error marshaling notification message: %v", err)
		return
	}

	// Broadcast the formatted message to the chat room
	cs.cm.BroadcastMessage(strconv.Itoa(conv.ChatRoomID), messageBytes)
}

// Ensure that your ConnectionManager's BroadcastMessage method supports broadcasting
// JSON-encoded messages to all connections in a specific chat room.
