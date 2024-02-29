package websocket

import (
	"ato_chat/models"
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

type Response struct {
	Conversations *models.ConversationWebsocket `json:"conversations"`
	Sidebar SidebarMessage `json:"sidebar"`
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

		

		// Marshal the modified object back to JSON
	

	
			log.Printf("[HandleMessages] Looking for messages.")
			responseMssg := Response{
				Sidebar:      sidebar,
				Conversations: message,
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

		// for c := range rooms[roomId]{
		// 	if err := c.WriteMessage(websocket.TextMessage,responseMsg); err != nil{
		// 		fmt.Println("Error writing message",err)
		// 		delete(rooms[roomId],c)
		// 	}

		// }

		

		// Send the JSON response back to the client
		// err = conn.WriteMessage(websocket.TextMessage, responseMsg)
		// if err != nil {
		// 	// Handle error
		// 	break
		// }
}

// func HandleMessages(){
// 	for {
// 		log.Printf("[HandleMessages] Looking for messages.")
// 		messages := <-broadcast
// 		log.Printf("[HandleMessages] Message getted.")

// 		for c := range rooms[roomId]{
// 			if err := c.WriteMessage(websocket.TextMessage,responseMsg); err != nil{
// 				fmt.Println("Error writing message",err)
// 				delete(rooms[roomId],c)
// 			}

// 		}

// 	}
// }