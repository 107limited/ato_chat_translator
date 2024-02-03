package models

import "time"

type ChatRoomDetail struct {
	UserID      int       `json:"user_id"`
	CompanyName string    `json:"company_name"`
	Name        string    `json:"name"`
	ChatRoomID  int       `json:"chat_room_id"`
	CreatedAt   time.Time `json:"created_at"`
}

// ChatRoom represents a chat room structure with basic user IDs and creation timestamp.
type ChatRoom struct {
	ID        int       `json:"id"`
	User1ID   int       `json:"user1_id"`
	User2ID   int       `json:"user2_id"`
	CreatedAt time.Time `json:"created_at"`
}
