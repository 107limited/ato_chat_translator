package chat

import (
	"ato_chat/models"
	"database/sql"
	"fmt"
	"log"
	"time"
)

// ConversationRepository adalah antarmuka untuk menyimpan dan mengambil percakapan.
type ConversationRepository interface {
	SaveConversation(conversation *models.Conversation) (sql.Result, error)
	GetAllConversations() ([]*models.Conversation, error)
}

// NewConversationRepository membuat instance baru dari ConversationRepository
func NewConversationRepository(db *sql.DB) *conversationRepository {
	return &conversationRepository{
		db: db,
	}
}

// conversationRepository adalah implementasi ConversationRepository
type conversationRepository struct {
	db *sql.DB
}

func (cr *conversationRepository) SaveConversation(conversation *models.Conversation) (sql.Result, error) {
	// Query untuk menyimpan percakapan baru
	query := `INSERT INTO conversations (japanese_text, english_text, user_id, speaker, company_id, chat_room_id, created_at, date) VALUES (?, ?, ?, ?, ?, ?, NOW(), ?)`
	// Gunakan conversation.Date untuk kolom date, yang sudah dalam format milidetik sejak epoch Unix
	res, err := cr.db.Exec(query, conversation.JapaneseText, conversation.EnglishText, conversation.UserID, conversation.Speaker, conversation.CompanyID, conversation.ChatRoomID, conversation.Date)
	if err != nil {
		return nil, fmt.Errorf("error executing query: %v", err)
	}

	// Konversi timestamp milidetik ke detik
	timestampInSeconds := conversation.Date / 1000

	// Memperbarui last_message_date di chat_room dengan timestamp dalam detik
	updateQuery := `UPDATE chat_room SET last_message_date = FROM_UNIXTIME(?) WHERE id = ?`
	_, err = cr.db.Exec(updateQuery, timestampInSeconds, conversation.ChatRoomID)
	if err != nil {
		return nil, fmt.Errorf("error updating last message date in chat_room: %v", err)
	}

	return res, nil
}

func (cr *conversationRepository) GetAllConversations() ([]*models.Conversation, error) {
	query := "SELECT id, japanese_text, english_text, user_id, company_id, chat_room_id, created_at, date, speaker, read_message FROM conversations ORDER BY created_at ASC"
	rows, err := cr.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var conversations []*models.Conversation
	for rows.Next() {
		var dateInt64 sql.NullInt64
		var createdAtString string // Menggunakan string untuk menampung created_at sementara
		var conv models.Conversation

		err := rows.Scan(&conv.ID, &conv.JapaneseText, &conv.EnglishText, &conv.UserID, &conv.CompanyID, &conv.ChatRoomID, &createdAtString, &dateInt64, &conv.Speaker, &conv.ReadMessage)

		if err != nil {
			return nil, err
		}

		// Konversi createdAtString ke time.Time
		conv.CreatedAt, err = time.Parse("2006-01-02 15:04:05", createdAtString)
		if err != nil {
			log.Printf("Error parsing created_at: %v", err)
			// Handle error sesuai kebutuhan
		}

		if dateInt64.Valid {
			conv.Date = dateInt64.Int64
		} else {
			conv.Date = -1
		}

		conversations = append(conversations, &conv)
	}

	return conversations, nil
}
