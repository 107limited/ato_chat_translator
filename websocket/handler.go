// File: handler.go
package websocket

import (
	// Importing necessary packages
	"ato_chat/chat"
	"ato_chat/models"
	"encoding/json"
	"fmt"
	"log"

	"net/http"
	"strconv"
	"time"

	// Importing the Gorilla WebSocket package
	"github.com/gorilla/websocket"
)

// Setting up the WebSocket upgrader with read and write buffer sizes
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// CheckOrigin returns true to allow all connections regardless of the origin
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// MessageFormat defines the structure of a message
type MessageFormat struct {
	RoomID  string `json:"roomID"`
	Message string `json:"message"`
}

// ConversationService represents a service for managing conversations
type ConversationService struct {
	repo chat.ConversationRepository // Repository for conversation data storage
	cm   *ConnectionManager          // Manager for WebSocket connections
}

// LastMessage structure for the last message in a conversation
type LastMessage struct {
	English  string `json:"english"`
	Japanese string `json:"japanese"`
	UserID   int    `json:"user_id"`
	Date     int64  `json:"date"`
}

// SidebarMessage structure for a message to be sent to the sidebar
type SidebarMessage struct {
	UserID          int         `json:"user_id"`
	CompanyName     string      `json:"company_name"`
	Name            string      `json:"name"`
	ChatRoomID      int         `json:"chat_room_id"`
	CreatedAt       string      `json:"created_at"`
	LastMessage     LastMessage `json:"last_message"`
	LastMessageUser int         `json:"last_message_user"`
}

// TypingMessage represents the structure of a typing notification.
type TypingMessage struct {
	ChatRoomID string `json:"chatRoomID"`
	UserID     int    `json:"userID"`
	IsTyping   bool   `json:"isTyping"`
}

type OnlineStatus struct {
    UserID    int  `json:"user_id"`
    IsOnline  bool `json:"is_online"`
}


// NewConversationService creates a new instance of ConversationService
func NewConversationService(repo chat.ConversationRepository, cm *ConnectionManager) *ConversationService {
	return &ConversationService{
		repo: repo,
		cm:   cm,
	}
}

// SaveAndBroadcast saves a conversation and broadcasts it to relevant users
func (cs *ConversationService) SaveAndBroadcast(conv models.Conversation) error {
	// Save the conversation
	_, err := cs.repo.SaveConversation(&conv)
	if err != nil {
		return err
	}

	// Prepare roomID and message for broadcasting
	roomID := fmt.Sprintf("%d", conv.ChatRoomID)
	messageBytes, err := json.Marshal(conv)
	if err != nil {
		return err // Return error if marshaling fails
	}

	// Broadcast the message to all connections in the appropriate room
	cs.cm.BroadcastMessage(roomID, messageBytes) // Use only roomID and messageBytes

	return nil
}

// SendSidebarMessage is a function to send sidebar messages
func (cs *ConversationService) SendSidebarMessage(conn *websocket.Conn, userID int, companyName, name string, chatRoomID int, lastMessage LastMessage) error {
	sidebarMessage := SidebarMessage{
		UserID:          userID,
		CompanyName:     companyName,
		Name:            name,
		ChatRoomID:      chatRoomID,
		CreatedAt:       time.Now().UTC().Format(time.RFC3339),
		LastMessage:     lastMessage,
		LastMessageUser: lastMessage.UserID,
	}

	// Encode the message into JSON
	messageBytes, err := json.Marshal(sidebarMessage)
	if err != nil {
		log.Printf("Failed to encode sidebar message: %v", err)
		return err
	}

	// Send the message through WebSocket
	if err := conn.WriteMessage(websocket.TextMessage, messageBytes); err != nil {
		log.Printf("Failed to send sidebar message: %v", err)
		return err
	}

	log.Printf("Sidebar message sent: %+v", sidebarMessage)
	return nil
}

// BroadcastTypingStatus broadcasts typing status to users in a chat room.
func (cs *ConversationService) BroadcastTypingStatus(typingMsg TypingMessage) {
	cs.cm.connectionsMu.Lock()
	defer cs.cm.connectionsMu.Unlock()

	// Iterate through all connections in the chat room and send the typing status.
	for conn := range cs.cm.Connections[typingMsg.ChatRoomID] {
		err := conn.WriteJSON(typingMsg)
		if err != nil {
			log.Printf("error broadcasting typing status: %v", err)
		}
	}
}

// HandleWebSocket creates an HTTP handler function to manage WebSocket connections
func HandleWebSocket(cs *ConversationService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Upgrade the HTTP server connection to the WebSocket protocol.
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("Upgrade to websocket failed:", err)
			return
		}
		defer conn.Close()

		// Retrieve the chat room ID from the URL query parameters.
		chatRoomID := r.URL.Query().Get("room_id")
		// If there is no chatRoomID, add the connection to all chat rooms.
		if chatRoomID == "" {
			cs.cm.AddGlobalConnection(conn)
			defer cs.cm.RemoveGlobalConnection(conn)
		} else {
			// If a chatRoomID is provided, add the connection to that specific chat room.
			cs.cm.AddConnection(chatRoomID, conn)
			defer cs.cm.RemoveConnection(chatRoomID, conn)
		}

		// Loop to continually read messages from the WebSocket connection
		for {
			// Read messages from the WebSocket
			_, p, err := conn.ReadMessage()
			if err != nil {
				log.Printf("read error: %v", err)
				break
			}

			// Unmarshal the JSON into a TypingMessage struct
			var typingMessage TypingMessage
			err = json.Unmarshal(p, &typingMessage)
			if err != nil {
				log.Printf("Error unmarshaling typing message: %v", err)
				continue
			}

			// Broadcast the typing status to other users in the same chat room.
			cs.BroadcastTypingStatus(typingMessage)

			// Unmarshal the JSON into a Conversation struct
			var conv models.Conversation
			err = json.Unmarshal(p, &conv)
			if err != nil {
				log.Printf("Error unmarshaling message: %v, message: %s", err, string(p))
				continue
			}

			// Convert chatRoomID from string to int and handle any error
			conv.ChatRoomID, err = strconv.Atoi(chatRoomID)
			if err != nil {
				log.Printf("Invalid chatRoomID: %v", err)
				continue
			}

			// Save the conversation and broadcast it
			err = cs.SaveAndBroadcast(conv)
			if err != nil {
				log.Printf("error saving and broadcasting message: %v", err)
			} else {
				// After successfully saving and broadcasting the message,

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
		CompanyName:     conv.CompanyName, // Assuming you fetch this from your database or have it in your conversation model
		Name:            conv.Speaker,     // Assuming you fetch this from your database based on UserID
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
