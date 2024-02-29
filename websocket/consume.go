package websocket

import (
	"encoding/json"
	"fmt"
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

type RequestConversation struct {
	ID           int    `json:"id"`
	JapaneseText string `json:"japanese_text"`
	EnglishText  string `json:"english_text"`
	Speaker      string `json:"speaker"`
	UserID       int    `json:"user_id"`
	UserID2      int    `json:"user2_id"`
	CompanyID    int    `json:"company_id"`
	ChatRoomID   int    `json:"chat_room_id"`
	CreatedAt    string `json:"created_at"`
	Date         int64  `json:"date"`
	UserName     string `json:"user_name"`
	CompanyName  string `json:"company_name"`
	Sidebars	[]SidebarMessage `json:"sidebars"`
}

type ResponseConversation struct{
	ID           int    `json:"id"`
	JapaneseText string `json:"japanese_text"`
	EnglishText  string `json:"english_text"`
	Speaker      string `json:"speaker"`
	UserID       int    `json:"user_id"`
	CompanyID    int    `json:"company_id"`
	ChatRoomID   int    `json:"chat_room_id"`
	CreatedAt    string `json:"created_at"`
	Date         int64  `json:"date"`

	
}


type Response struct {
	Conversations *ResponseConversation `json:"conversations"`
	Sidebar *[]SidebarMessage `json:"sidebars"`
}

var rooms = make(map[string]map[*websocket.Conn]bool)

func HandleWSL(w http.ResponseWriter, r *http.Request) {

	roomId := r.URL.Query().Get("roomId")

	
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		// Handle error
		return
	}
	defer conn.Close()



	if _,ok := rooms[roomId]; !ok{
		rooms[roomId]= make(map[*websocket.Conn]bool)
	}

	rooms[roomId][conn] = true



	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Request not found")
			break
		}
	

		// Unmarshal JSON from the incoming message
		var message *RequestConversation
		err = json.Unmarshal(msg, &message)
		if err != nil {
			log.Println("Request")
			break
		}
		


			sidebars := message.Sidebars
			resMessage := ResponseConversation{
				ID: message.ID,
				JapaneseText: message.JapaneseText,
				EnglishText: message.EnglishText,
				Speaker: message.Speaker,
				UserID: message.UserID,
				CompanyID: message.CompanyID,
				ChatRoomID: message.ChatRoomID,
				CreatedAt: "",
				Date: message.Date,
			}
		
			
		

		

		// Marshal the modified object back to JSON
	

	
			log.Printf("[HandleMessages] Looking for messages.")
			responseMssg := Response{
				Sidebar : &sidebars,
				Conversations: &resMessage,
			}
			log.Printf("[HandleMessages] Message getted.")
	
			responseMsg, err := json.Marshal(responseMssg)
		if err != nil {
			// Handle error
			break
		}

			for c := range rooms[roomId]{
				if err := c.WriteMessage(websocket.TextMessage,responseMsg); err != nil{
					fmt.Println("Error writing message",err)
					delete(rooms[roomId],c)
				}
	
			}

	}

		
}
