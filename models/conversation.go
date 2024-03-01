package models

import (
	"errors"
	"time"
)

type Conversation struct {
	ID                int       `json:"id"`
	JapaneseText      string    `json:"japanese_text"`
	EnglishText       string    `json:"english_text"`
	Speaker           string    `json:"speaker"`
	UserID            int       `json:"user_id"`
	CompanyID         int       `json:"company_id"`
	CompanyName       string    `json:"company_name"`
	ChatRoomID        int       `json:"chat_room_id"`
	OriginalMessage   string    `json:"original_message"`
	TranslatedMessage string    `json:"translated_message"`
	CreatedAt         time.Time `json:"created_at"`
	Date              int64     `json:"date"`
	ReadMessage       bool      `json:"read_message"`
}

// TranslationRequest represents the JSON structure for translation request
type TranslationRequest struct {
	ID              int    `json:"id"`
	User1ID         int    `json:"user1_id"` // Pastikan tag JSON sesuai dengan key di request body
	User2ID         int    `json:"user2_id"`
	Speaker         string `json:"speaker"`
	CompanyID       int    `json:"company_id"`
	ChatRoomID      int    `json:"chat_room_id"`
	OriginalMessage string `json:"original_message"`
	JapaneseText    string `json:"japanese_text"`
	EnglishText     string `json:"english_text"`
	Date            int64  `json:"date"`
}

// You can add a method to validate the struct
func (req *TranslationRequest) Validate() error {
	if req.User1ID == 0 || req.Speaker == "" || req.CompanyID == 0 || req.OriginalMessage == "" || req.Date == 0 {
		return errors.New("missing required fields")
	}
	return nil
}

type GetAllConversations struct {
	ID                int    `json:"id"`
	UserID            string `json:"user_id"`
	Speaker           string `json:"speaker"`
	CompanyID         int    `json:"company_id"`
	ChatRoomID        string `json:"chat_room_id"`
	OriginalMessage   string `json:"original_message"`
	TranslatedMessage string `json:"translated_message"`
	CreatedAt         string `json:"created_at"`
}

// TranslationResponse represents the JSON structure for translation response
type TranslationResponse struct {
	Conversations []struct {
		Speaker           string `json:"speaker"`
		OriginalMessage   string `json:"original_message"`
		TranslatedMessage string `json:"translated_message"`
		CompanyName       string `json:"company_name"`
		ChatRoomID        int    `json:"chat_room_id"` // Pastikan field ini ada dalam definisi
		UserID            int    `json:"user_id"`      // Tambahkan field ini
	} `json:"conversations"`
}

type TranslationResponseTranslateHandler struct {
	Conversations []ConversationDetail `json:"conversations"`
}

type ConversationDetail struct {
	Speaker           string `json:"speaker"`
	OriginalMessage   string `json:"original_message"`
	TranslatedMessage string `json:"translated_message"`
	CompanyName       string `json:"company_name"`
	// Exclude ChatRoomID if it's not needed for this response
}
