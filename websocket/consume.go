package websocket

import (
	"ato_chat/models"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// Response struct as provided.
type Response struct {
	Conversations *models.ConversationWebsocket `json:"conversations"`
	Sidebar       SidebarMessage                `json:"sidebar"`
}


func HandleWSL(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade failed: %v", err)
		return
	}
	defer conn.Close()
	log.Println("WebSocket connection successfully upgraded.")

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading message: %v", err)
			break
		}

		var message *models.ConversationWebsocket
		err = json.Unmarshal(msg, &message)
		if err != nil {
			log.Printf("Error unmarshaling message: %v", err)
			break
		}
		log.Println("Message successfully unmarshaled.")

		

	// 	parseMessage := *models.ConversationsWebsocket{
	// 		ID :message.Id,        
	// User1ID :message.UserId,         
	// User2ID :message.ToId,         
	// Speaker :,         
	// CompanyID,       
	// ChatRoomID,      
	// OriginalMessage, 
	// JapaneseText,    
	// EnglishText,     
	// Date,        
	// 	}

	lastMessage := LastMessage{
		UserID:   message.UserID,
		English:  message.EnglishText,
		Japanese: message.JapaneseText,
		Date:     message.Date,
	}

	var company string
	if message.CompanyID == 1 {
		company = "ATO"
	} else {
		company = "107"
	}

	sidebar := SidebarMessage{
		UserID:      message.UserID2,
		CompanyName: company,
		Name:        message.UserName,
		ChatRoomID:  message.ChatRoomID,
		CreatedAt:   "",
		LastMessage: lastMessage,
	}

	responseMsg := Response{
		Sidebar:       sidebar,
		Conversations: message,
	}

	responseJSON, err := json.Marshal(responseMsg)
	if err != nil {
		log.Printf("Error marshaling response: %v", err)
		break
	}

	err = conn.WriteMessage(websocket.TextMessage, responseJSON)
	if err != nil {
		log.Printf("Error sending message: %v", err)
		break
	}
	log.Println("Message successfully sent to the client.")
}
}