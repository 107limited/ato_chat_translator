package websocket

import (
	"ato_chat/models"
	"encoding/json"
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
	Sidebar SidebarMessage `json:"sidebar"`
}


func HandleWSL(w http.ResponseWriter, r *http.Request) {
	
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		// Handle error
		return
	}
	defer conn.Close()

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			// Handle error
			break
		}

		// Unmarshal JSON from the incoming message
		var message *models.ConversationWebsocket
		err = json.Unmarshal(msg, &message)
		if err != nil {
			// Handle error
			break
		}

		

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
			CreatedAt:   "", 
			LastMessage: lastmessage,
		}

		responseMssg := Response{
			Sidebar:      sidebar,
			Conversations: message,
		}

		// Marshal the modified object back to JSON
		responseMsg, err := json.Marshal(responseMssg)
		if err != nil {
			// Handle error
			break
		}

		// Send the JSON response back to the client
		err = conn.WriteMessage(websocket.TextMessage, responseMsg)
		if err != nil {
			// Handle error
			break
		}
	}
}