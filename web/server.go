package web

import (
	"ato_chat/chat"
	"ato_chat/models"
	"ato_chat/translation"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"time"

	"github.com/gorilla/mux"
	"golang.org/x/text/language"
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
		http.Error(w, "Failed to parse request body", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if err := translationRequest.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Log request body
	fmt.Printf("Received Request Body: %+v\n", translationRequest.OriginalMessage)

	// Translate original message
	translatedMessage, err := s.GPT4Translator.TranslateMessage(translationRequest.OriginalMessage)
	if err != nil {
		errorMessage := fmt.Sprintf("Failed to translate message: %v", err)
		http.Error(w, errorMessage, http.StatusInternalServerError)
		return
	}

	// Determine language of original message
	var japaneseText, englishText string
	if s.isJapanese(translationRequest.OriginalMessage) {
		japaneseText = translationRequest.OriginalMessage
		englishText = translatedMessage
	} else {
		japaneseText = translatedMessage
		englishText = translationRequest.OriginalMessage
	}

	// Create Conversation object
	// var conversation *models.Conversation
	t := models.Conversation{
		JapaneseText:      japaneseText,
		EnglishText:       englishText,
		Speaker:           translationRequest.Speaker,
		UserID:            translationRequest.UserID,
		CompanyID:         translationRequest.CompanyID,
		ChatRoomID:        translationRequest.ChatRoomID,
		OriginalMessage:   translationRequest.OriginalMessage,
		TranslatedMessage: translatedMessage,
		CreatedAt:         time.Now(),
	}
	// conversation = &t
	// Save conversation to repository
	err = s.ConversationRepo.SaveConversation(&t)
	if err != nil {
		errorMessage := fmt.Sprintf("Failed to save conversation: %v", err)
		http.Error(w, errorMessage, http.StatusInternalServerError)
		return
	}

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
func (s *Server) isJapanese(text string) bool {
	// Identifikasi bahasa menggunakan golang.org/x/text/language
	tag, err := language.Parse(text)
	if err != nil {
		// Handle error jika parsing gagal
		return false
	}

	// Bandingkan dengan tag bahasa Jepang
	return tag == language.Japanese
}

// GetAllConversationsHandler menangani permintaan untuk mendapatkan semua percakapan
func (s *Server) GetAllConversationsHandler(w http.ResponseWriter, r *http.Request) {
	conversations, err := s.ConversationRepo.GetAllConversations()
	if err != nil {
		log.Printf("Error retrieving conversations: %v", err)
		http.Error(w, "Failed to retrieve conversations", http.StatusInternalServerError)
		return
	}

	// Send successful response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(conversations)
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
