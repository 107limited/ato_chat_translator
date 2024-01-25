package web

import (
	"ato_chat/chat"
	"ato_chat/models"
	"ato_chat/translation"
	"encoding/json"
	"fmt"
	"net/http"

	"time"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	//"golang.org/x/text/language"
)

// Server adalah struktur data untuk server web
type Server struct {
	Router           *mux.Router
	ConversationRepo chat.ConversationRepository
	GPT4Translator   translation.Translator
}

// NewServer membuat instance baru dari Server
func NewServer(conversationRepo chat.ConversationRepository, gpt4Translator translation.Translator) *Server {
	router := mux.NewRouter()

	server := &Server{
		Router:           router,
		ConversationRepo: conversationRepo,
		GPT4Translator:   gpt4Translator,
	}

	server.initializeRoutes()

	return server
}

// SaveConversationHandler menangani permintaan untuk menyimpan percakapan
func (s *Server) SaveConversationHandler(w http.ResponseWriter, r *http.Request) {
	// Parse JSON request body
	var translationRequest models.TranslationRequest
	if err := json.NewDecoder(r.Body).Decode(&translationRequest); err != nil {
		log.WithError(err).Error("Failed to parse request body")
		http.Error(w, "Failed to parse request body", http.StatusBadRequest)
		return
	}

	// Konversi nilai "date" ke int64
    var dateInt64 int64
    if translationRequest.Date >= 0 {
        dateInt64 = int64(translationRequest.Date)
    } else {
        log.Warn("Invalid 'date' value")
        http.Error(w, "Invalid 'date' value", http.StatusBadRequest)
        return
    }

	// Validate required fields
    if err := translationRequest.Validate(); err != nil {
        log.WithError(err).Warn("Validation failed for translation request")
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    // Translate original message
    translatedMessage, err := s.GPT4Translator.TranslateMessage(translationRequest.OriginalMessage)
    if err != nil {
        log.WithError(err).Error("Failed to translate message")
        http.Error(w, fmt.Sprintf("Failed to translate message: %v", err), http.StatusInternalServerError)
        return
    }

	// Tentukan bahasa berdasarkan speaker
	var japaneseText, englishText string
	if translationRequest.Speaker == "Ato" {
		japaneseText = translationRequest.OriginalMessage
		englishText = translatedMessage
	} else {
		englishText = translationRequest.OriginalMessage
		japaneseText = translatedMessage
	}

	// Create Conversation object
	t := models.Conversation{
		Speaker:           translationRequest.Speaker,
		JapaneseText:      japaneseText,
		EnglishText:       englishText,
		UserID:            translationRequest.UserID,
		CompanyID:         translationRequest.CompanyID,
		ChatRoomID:        translationRequest.ChatRoomID,
		OriginalMessage:   translationRequest.OriginalMessage,
		TranslatedMessage: translatedMessage,
		CreatedAt:         time.Now(),
		Date:              dateInt64,
	}

	// Save conversation to repository
    err = s.ConversationRepo.SaveConversation(&t)
    if err != nil {
        log.WithError(err).Error("Failed to save conversation")
        http.Error(w, fmt.Sprintf("Failed to save conversation: %v", err), http.StatusInternalServerError)
        return
    }

    log.Info("Conversation saved successfully")

	// Create TranslationResponse
	translationResponse := models.TranslationResponse{
		Conversations: []struct {
			Speaker           string `json:"speaker"`
			OriginalMessage   string `json:"original_message"`
			TranslatedMessage string `json:"translated_message"`
		}{
			{
				Speaker:           translationRequest.Speaker,
				OriginalMessage:   translationRequest.OriginalMessage,
				TranslatedMessage: translatedMessage,
			},
		},
	}

	// Send success response
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(translationResponse)
}

// isJapanese checks if the given text is in Japanese
// func (s *Server) isJapanese(text string) bool {
// 	// Identifikasi bahasa menggunakan golang.org/x/text/language
// 	tag, err := language.Parse(text)
// 	if err != nil {
// 		// Handle error jika parsing gagal
// 		return false
// 	}

// 	// Bandingkan dengan tag bahasa Jepang
// 	return tag == language.Japanese
// }

// GetAllConversationsHandler menangani permintaan untuk mendapatkan semua percakapan
func (s *Server) GetAllConversationsHandler(w http.ResponseWriter, r *http.Request) {
	conversations, err := s.ConversationRepo.GetAllConversations()
	if err != nil {
		log.Printf("Error retrieving conversations: %v", err)
		http.Error(w, "Failed to retrieve conversations", http.StatusInternalServerError)
		return
	}

	// Mengkonversi data percakapan ke dalam format response yang diinginkan
	var response []map[string]interface{}
	for _, conv := range conversations {
		conversationMap := map[string]interface{}{
			"id":            conv.ID,
			"japanese_text": conv.JapaneseText,
			"english_text":  conv.EnglishText,
			"speaker":       conv.Speaker,
			"user_id":       conv.UserID,
			"company_id":    conv.CompanyID,
			"chat_room_id":  conv.ChatRoomID,
			"created_at":    conv.CreatedAt,
			"date":          conv.Date,
		}
		response = append(response, conversationMap)
	}

	// Mengirim response JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

}

// TranslateMessageHandler menangani permintaan untuk menerjemahkan pesan
func (s *Server) TranslateMessageHandler(w http.ResponseWriter, r *http.Request) {
	// Parse JSON request body
	var translationRequest models.TranslationRequest
	if err := json.NewDecoder(r.Body).Decode(&translationRequest); err != nil {
		http.Error(w, "Failed to parse request body", http.StatusBadRequest)
		return
	}

	// Mendapatkan terjemahan menggunakan GPT-3.5-turbo
	translatedMessage, err := s.GPT4Translator.TranslateMessage(translationRequest.OriginalMessage)
	if err != nil {
		http.Error(w, "Failed to translate message", http.StatusInternalServerError)
		return
	}

	// Membuat objek hasil terjemahan
	translationResponse := models.TranslationResponse{
		Conversations: []struct {
			Speaker           string `json:"speaker"`
			OriginalMessage   string `json:"original_message"`
			TranslatedMessage string `json:"translated_message"`
		}{
			{
				Speaker:           translationRequest.Speaker,
				OriginalMessage:   translationRequest.OriginalMessage,
				TranslatedMessage: translatedMessage,
			},
		},
	}

	// Kirim response sukses
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(translationResponse)

}

// initializeRoutes mengatur rute-rute untuk server
func (s *Server) initializeRoutes() {
	// Contoh rute
	s.Router.HandleFunc("/api/conversations", s.SaveConversationHandler).Methods("POST")
	s.Router.HandleFunc("/api/conversations", s.GetAllConversationsHandler).Methods("GET")
	s.Router.HandleFunc("/api/translate", s.TranslateMessageHandler).Methods("POST")
}

// Start menjalankan server web
func (s *Server) Start(port string) {
	fmt.Printf("Server is running on port %s...\n", port)
	http.ListenAndServe(":"+port, s.Router)
}
