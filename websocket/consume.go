package websocket

import (
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/gorilla/websocket"
)

type RequestConversation struct {
	ID           int               `json:"id"`
	JapaneseText string            `json:"japanese_text"`
	EnglishText  string            `json:"english_text"`
	Speaker      string            `json:"speaker"`
	UserID       int               `json:"user_id"`
	UserID2      int               `json:"user2_id"`
	CompanyID    int               `json:"company_id"`
	ChatRoomID   int               `json:"chat_room_id"`
	CreatedAt    string            `json:"created_at"`
	Date         int64             `json:"date"`
	UserName     string            `json:"user_name"`
	CompanyName  string            `json:"company_name"`
	Sidebars     *ResponseSidebars `json:"sidebars"`
}

type ResponseConversation struct {
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
	UserID        int                   `json:"user_id"`
	UserID2       int                   `json:"user2_id"`
	Conversations *ResponseConversation `json:"conversations"`
	Sidebar       *ResponseSidebars     `json:"sidebars"`
}

type Messages struct {
	UserId       string `json:"userId"`
	TargetUserId string `json:"TargetuserId"`
	Contents     string `json:"contents"`
}

type Client struct {
	Conn *websocket.Conn
	Id   int
}

var clients = make(map[*websocket.Conn]int)

// var broadcast = make(chan Messages)
var broadcast = make(chan Response)
var mutex = &sync.Mutex{}

func CreateClient(conn *websocket.Conn, id int) *Client {
	mutex.Lock()
	defer mutex.Unlock()
	clients[conn] = id
	return &Client{Conn: conn, Id: id}
}

func HandleWSL(w http.ResponseWriter, r *http.Request) {

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Upgraded failed")
		return
	}
	defer conn.Close()

	userId := r.URL.Query().Get("userId")

	i, err := strconv.Atoi(userId)
	if err != nil {
		log.Println("Error parse string")
	}

	clients[conn] = i

	log.Println("UserId :", userId, "Connected")

	for {
		var msg RequestConversation
		err := conn.ReadJSON(&msg)
		if err != nil {
			delete(clients, conn)
			break
		}

		// sidebars := msg.Sidebars
		resMessage := ResponseConversation{
			ID:           msg.ID,
			JapaneseText: msg.JapaneseText,
			EnglishText:  msg.EnglishText,
			Speaker:      msg.Speaker,
			UserID:       msg.UserID,
			CompanyID:    msg.CompanyID,
			ChatRoomID:   msg.ChatRoomID,
			CreatedAt:    "",
			Date:         msg.Date,
		}

		responeBroadcast := Response{
			UserID:        msg.UserID,
			UserID2:       msg.UserID2,
			Conversations: &resMessage,
			Sidebar:       msg.Sidebars,
		}

		broadcast <- responeBroadcast
	}

}

func HandleMessages() {
	for {
		log.Println("Waiting for message")
		msg := <-broadcast
		log.Println("Getting message from :", msg.UserID, "to :", msg.UserID2)

		// senderIdStr := strconv.Itoa(msg.UserID)
		// receiverIdStr := strconv.Itoa(msg.UserID2)

		mutex.Lock()
		for client, uId := range clients {
			if uId == msg.UserID || uId == msg.UserID2 {
				if err := client.WriteJSON(msg); err != nil {
					log.Printf("Error: %v", err)
					client.Close()
					delete(clients, client)
				}
			}
		}

	}
}
