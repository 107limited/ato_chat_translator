package websocket

import (
	"ato_chat/models"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// chat_room_id: number;
// company_id: number;
// created_at: string;
// date: number;
// english_text: string;
// id: number;
// japanese_text: string;
// speaker: string;
// user_id: number;

type Response struct {
	Conversations *models.ConversationWebsocket `json:"conversations"`
	Sidebar       SidebarMessage                `json:"sidebar"`
}

func HandleWSL(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Error upgrading to WebSocket: %v", err)
		return
	}
	defer conn.Close()
	log.Println("WebSocket connection established.")

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading WebSocket message: %v", err)
			break
		}

		var message *models.ConversationWebsocket
		err = json.Unmarshal(msg, &message)
		if err != nil {
			log.Printf("Error unmarshaling message: %v", err)
			break
		}
		log.Println("Received and unmarshaled message successfully.")

		lastmessage := LastMessage{
			UserID:   message.UserID,
			English:  message.EnglishText,
			Japanese: message.JapaneseText,
			Date:     message.Date,
		}

		sidebar := SidebarMessage{
			UserID:      message.UserID2,
			CompanyName: message.CompanyName,
			Name:        message.UserName,
			ChatRoomID:  message.ChatRoomID,
			CreatedAt:   "", // consider using time.Now().Format(...) if you want to include the current time
			LastMessage: lastmessage,
		}

		responseMssg := Response{
			Sidebar:       sidebar,
			Conversations: message,
		}

		responseMsg, err := json.Marshal(responseMssg)
		if err != nil {
			log.Printf("Error marshaling response message: %v", err)
			break
		}

		err = conn.WriteMessage(websocket.TextMessage, responseMsg)
		if err != nil {
			log.Printf("Error sending response message: %v", err)
			break
		}
		log.Println("Response message sent successfully.")
	}
}
