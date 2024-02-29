package models

type ConversationWebsocket struct {
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
}

type LastMessageWebSocket struct {
	English  string `json:"english"`
	Japanese string `json:"japanese"`
	UserID   int    `json:"user_id"`
	Date     int64  `json:"date"`
}