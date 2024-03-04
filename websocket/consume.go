package websocket

import (
	"log"
	"net/http"
	"strconv"
	"sync" // Import sync package for mutex

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
	Id   string
}

var (
	clients   = make(map[*websocket.Conn]string)
	clientsMu sync.Mutex // Mutex to synchronize access to clients map
)

var broadcast = make(chan Response)

func CreateClient(conn *websocket.Conn, id string) *Client {
	return &Client{Conn: conn, Id: id}
}

// Setting up the WebSocket upgrader with read and write buffer sizes
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// CheckOrigin returns true to allow all connections regardless of the origin
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func HandleWSL(w http.ResponseWriter, r *http.Request) {

	conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Printf("Upgrade failed: %v", err)
        return
    }
    defer conn.Close() // Letakkan di awal fungsi untuk memastikan koneksi ditutup ketika fungsi selesai dijalankan

    userId := r.URL.Query().Get("userId")
    clientsMu.Lock()
    clients[conn] = userId
    clientsMu.Unlock()

    log.Println("UserId:", userId, "Connected")

	for {
		var msg RequestConversation
		err := conn.ReadJSON(&msg)
		if err != nil {
			clientsMu.Lock()
			delete(clients, conn)
			clientsMu.Unlock()
			break
		}

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
		msg := <-broadcast

		// Konversi UserID dan UserID2 ke string
		userIDStr := strconv.Itoa(msg.UserID)
		userID2Str := strconv.Itoa(msg.UserID2)

		clientsMu.Lock()
		// Iterasi melalui semua klien untuk menemukan dan mengirim pesan ke UserID dan UserID2
		for client, uId := range clients {
			// Cek apakah client adalah UserID atau UserID2
			if uId == userIDStr || uId == userID2Str {
				if err := client.WriteJSON(msg); err != nil {
					log.Printf("Error: %v", err)
					client.Close()
					delete(clients, client)
				}
			}
		}
		clientsMu.Unlock()
	}
}

