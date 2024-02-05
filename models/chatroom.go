package models

import "time"


type LastMessage struct {
    English  string `json:"english"`
    Japanese string `json:"japanese"`
}

type ChatRoomDetail struct {
	UserID      int       `json:"user_id"`
	CompanyName string    `json:"company_name"`
	Name        string    `json:"name"`
	ChatRoomID  int       `json:"chat_room_id"`
	CreatedAt   time.Time `json:"created_at"`
	LastMessage LastMessage `json:"last_message,omitempty"`
	
}

// ChatRoom represents a chat room structure with basic user IDs and creation timestamp.
type ChatRoom struct {
	ID        int       `json:"id"`
	User1ID   int       `json:"user1_id"`
	User2ID   int       `json:"user2_id"`
	CreatedAt time.Time `json:"created_at"`
}

type ChatRoomResponse struct {
	Message    string `json:"message"`
	ChatRoomID int64  `json:"chat_room_id,omitempty"` // Omitempty akan menyembunyikan field ini jika nilainya adalah 0
}
