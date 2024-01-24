package chat

import (
	"ato_chat/models"
	"database/sql"
	"fmt"
)

// ConversationRepository adalah antarmuka untuk menyimpan dan mengambil percakapan.
type ConversationRepository interface {
	SaveConversation(conversation *models.Conversation) error
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

// SaveConversation menyimpan percakapan ke database
func (cr *conversationRepository) SaveConversation(conversation *models.Conversation) error {
	query := "INSERT INTO conversations (japanese_text, english_text, user_id, speaker, company_id, chat_room_id, created_at, date) VALUES (?, ?, ?, ?, ?, ?, ?, ?)"
	result, err := cr.db.Exec(query, conversation.JapaneseText, conversation.EnglishText, conversation.UserID, conversation.Speaker, conversation.CompanyID, conversation.ChatRoomID, conversation.CreatedAt, conversation.Date)
	if err != nil {
		// Handle error saat penyimpanan ke dalam database
		return fmt.Errorf("error executing query: %v", err)
	}

	// Dapatkan ID yang dihasilkan
	lastInsertID, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("error getting last insert ID: %v", err)
	}

	// Atur nilai id pada conversation
	conversation.ID = int(lastInsertID)

	return nil
}

func (cr *conversationRepository) GetAllConversations() ([]*models.Conversation, error) {
	query := "SELECT * FROM conversations ORDER BY created_at ASC"
	rows, err := cr.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var conversations []*models.Conversation
	for rows.Next() {
		
		var dateInt64 sql.NullInt64
		var conv models.Conversation
		var get models.GetAllConversations
		err := rows.Scan(&conv.ID, &conv.JapaneseText, &conv.EnglishText, &conv.UserID, &conv.CompanyID, &conv.ChatRoomID, &get.CreatedAt, &dateInt64, &conv.Speaker)

		if err != nil {
			return nil, err
		}

	

		if dateInt64.Valid {
			conv.Date = dateInt64.Int64
		} else if dateInt64.Int64 == 0 {
			conv.Date = 0
		} else {
			conv.Date = -1
		}

		conversations = append(conversations, &conv)
	}

	return conversations, nil
}
