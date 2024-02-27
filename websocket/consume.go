package websocket

import (
	"ato_chat/models"
	"encoding/json"
	"net/http"
	"time"

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
	Conversations MessageFmt `json:"conversations"`
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
		var message *models.IMessage
		err = json.Unmarshal(msg, &message)
		if err != nil {
			// Handle error
			break
		}

		

		parseMessage := Messagefmt{
			Id:        message.Id,
			ToId:      message.ToId,
			UserId:    message.UserId,
			Message:   message.Message,
			Date:      message.Date,
			RoomId:    14,
			CompanyId: 2,
			CreatedAt: time.Now().UTC().Format(time.RFC3339),
			English:   "Text",
			Japanese:  "Con",
			Speaker:   "ATO",
			
		}

		lastmessage := LastMessage{
			UserID:   parseMessage.UserId,
			English:  "naoan",
			Japanese: "aihsd",
			Date:     int64(message.Date),
		}

		sidebar := SidebarMessage{
			UserID:      message.UserId,
			CompanyName: parseMessage.Speaker,
			Name:        "Test",
			ChatRoomID:  parseMessage.RoomId,
			CreatedAt:   parseMessage.CreatedAt,
			LastMessage: lastmessage,
		}

		responseMssg := Response{
			Sidebar:      sidebar,
			Conversations: parseMessage,
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