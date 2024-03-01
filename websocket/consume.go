package websocket

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)



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
	Sidebars	*ResponseSidebars `json:"sidebars"`
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

type ResponseSidebars struct {
	ATO *[]SidebarMessage `json:"ato_sidebars"`
	SNT *[]SidebarMessage `json:"snt_sidebars"`
}



type Response struct {
	Conversations *ResponseConversation `json:"conversations"`
	Sidebar *ResponseSidebars `json:"sidebars"`
}

var client = make(map[string]map[*websocket.Conn]bool)
var clientMutex = &sync.Mutex{} 

func HandleWSL(w http.ResponseWriter, r *http.Request) {

	userId := r.URL.Query().Get("userId")

	
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Upgraded failed")
		return
	}
	defer conn.Close()

	clientMutex.Lock() // Lock mutex sebelum mengakses atau memodifikasi map client
	if _, ok := client[userId]; !ok {
		client[userId] = make(map[*websocket.Conn]bool)
	}
	client[userId][conn] = true
	clientMutex.Unlock() // Unlock mutex setelah selesai mengakses atau memodifikasi map client

	defer func() {
		clientMutex.Lock() // Pastikan mengunci mutex sebelum menghapus koneksi
		delete(client[userId], conn)
		clientMutex.Unlock() // Unlock mutex setelah menghapus koneksi
	}()



	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Request not found")
			break
		}
	
		var message *RequestConversation
		err = json.Unmarshal(msg, &message)
		if err != nil {
			log.Println("Request")
			break
		}
		


			sidebars := ResponseSidebars{
				ATO: message.Sidebars.ATO,
				SNT: message.Sidebars.SNT,
			}
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

			for c := range client[userId]{
				if err := c.WriteMessage(websocket.TextMessage,responseMsg); err != nil{
					fmt.Println("Error writing message",err)
					delete(client[userId],c)
				}
	
			}

	}

		
}