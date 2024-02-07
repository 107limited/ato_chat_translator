package models

import (
	
	"time"
)

type LastMessage struct {
    English  string `json:"english"`
    Japanese string `json:"japanese"`
    UserID   int64  `json:"user_id,omitempty"` // Gunakan int64 jika sesuai dengan skema data Anda
    Date     int64 `json:"date"`
}


type ChatRoomDetail struct {
	UserID          int           `json:"user_id"`
	CompanyName     string        `json:"company_name"`
	Name            string        `json:"name"`
	ChatRoomID      int           `json:"chat_room_id"`
	CreatedAt       time.Time     `json:"created_at"`
	LastMessage     LastMessage   `json:"last_message,omitempty"`
	//LastMessageUser sql.NullInt64 `json:"last_message_user_id,omitempty"`
	LastMessageUser int64 `json:"last_message_user,omitempty"`
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

type CreateChatRoomAndMessageRequest struct {
	User1ID           int       `json:"user1_id"`
	User2ID           int       `json:"user2_id"`
	Speaker           string    `json:"speaker"`
	CompanyID         int       `json:"company_id"`
	OriginalMessage   string    `json:"original_message"`
	Date              int64     `json:"date"`
	ChatRoomID        int       `json:"chat_room_id"`
	TranslatedMessage string    `json:"translated_message"`
	CreatedAt         time.Time `json:"created_at"`
}
