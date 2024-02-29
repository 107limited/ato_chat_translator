package websocket

import (
	"ato_chat/models"
	"encoding/json"
	
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type Response struct {
	Conversations *models.ConversationWebsocket `json:"conversations"`
	Sidebar       SidebarMessage                `json:"sidebar"`
}

var rooms = make(map[string]map[*websocket.Conn]bool)
var roomsMutex = &sync.Mutex{} // Mutex untuk mengelola akses konkuren ke map rooms

func HandleWSL(w http.ResponseWriter, r *http.Request) {

	roomId := r.URL.Query().Get("roomId")
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to upgrade to websocket: %v", err)
		http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
		return
	}
	defer conn.Close()

	roomsMutex.Lock() // Lock mutex sebelum mengakses atau memodifikasi map rooms
    if _, ok := rooms[roomId]; !ok {
        rooms[roomId] = make(map[*websocket.Conn]bool)
    }
    rooms[roomId][conn] = true
    roomsMutex.Unlock() // Unlock mutex setelah selesai mengakses atau memodifikasi map rooms

    defer func() {
        roomsMutex.Lock() // Pastikan mengunci mutex sebelum menghapus koneksi
        delete(rooms[roomId], conn)
        roomsMutex.Unlock() // Unlock mutex setelah menghapus koneksi
    }()

	
	

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break // Exit the loop if there's an error (e.g., connection closed)
		}

		var message *models.ConversationWebsocket
		if err := json.Unmarshal(msg, &message); err != nil {
			log.Printf("Error unmarshalling message: %v", err)
			continue // Skip processing this message but continue listening
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
			Sidebar:       sidebar,
			Conversations: message,
		}
		log.Printf("[HandleMessages] Message getted.")

		responseJSON, err := json.Marshal(responseMssg)
		if err != nil {
			log.Printf("Error marshaling response: %v", err)
			continue // Skip sending this message but continue listening
		}

		broadcastToRoom(roomId, responseJSON)


		
		

	}

}
// Broadcast the message to all clients in the room
func broadcastToRoom(roomId string, message []byte) {
	for conn := range rooms[roomId] {
		if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
			log.Printf("Error broadcasting to room %s, err: %v", roomId, err)
			delete(rooms[roomId], conn)
			conn.Close()
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
