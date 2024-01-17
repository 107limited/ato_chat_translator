package chat

import (
	"ato_chat/models"
	"encoding/json"
	"net/http"
	"time"
)

// ConversationHandler adalah handler untuk percakapan
type ConversationHandler struct {
	Repository ConversationRepository
}

// NewConversationHandler membuat instance baru dari ConversationHandler
func NewConversationHandler(repo ConversationRepository) *ConversationHandler {
	return &ConversationHandler{
		Repository: repo,
	}
}

// SaveConversationHandler menangani permintaan untuk menyimpan percakapan
func (ch *ConversationHandler) SaveConversationHandler(w http.ResponseWriter, r *http.Request) {
	var conv models.Conversation
	err := json.NewDecoder(r.Body).Decode(&conv)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	conv.CreatedAt = time.Now()

	// err = ch.Repository.SaveConversation(&conv)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }

	w.WriteHeader(http.StatusCreated)
}

// GetAllConversationsHandler menangani permintaan untuk mendapatkan semua percakapan
func (ch *ConversationHandler) GetAllConversationsHandler(w http.ResponseWriter, r *http.Request) {
	conversations, err := ch.Repository.GetAllConversations()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(conversations)
}
