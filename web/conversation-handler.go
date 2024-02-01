package web

import (
	"ato_chat/chat"
	"ato_chat/jwt"
	"ato_chat/models"
	"ato_chat/translation"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"time"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	//"golang.org/x/text/language"
)

type Server struct {
	DB               *sql.DB
	Router           *mux.Router
	ConversationRepo chat.ConversationRepository
	GPT4Translator   translation.Translator
	ChatRoomHandler  *ChatRoomHandler
}

// Konstruktor untuk membuat instance baru dari Server dengan ChatRoomHandler
func NewServer(db *sql.DB, conversationRepo chat.ConversationRepository, gpt4Translator translation.Translator, chatRoomHandler *ChatRoomHandler) *Server {
	server := &Server{
		DB:               db,
		Router:           mux.NewRouter(),
		ConversationRepo: conversationRepo,
		GPT4Translator:   gpt4Translator,
		ChatRoomHandler:  chatRoomHandler,
	}
	server.initializeRoutes() // Inisialisasi rute setelah semua handler siap
	return server
}

// SaveConversationHandler menangani permintaan untuk menyimpan percakapan
func (s *Server) SaveConversationHandler(w http.ResponseWriter, r *http.Request) {
	// Ekstrak token dari header Authorization
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Authorization header is required", http.StatusUnauthorized)
		return
	}

	// Biasanya, token dikirim sebagai "Bearer <token>", jadi kita perlu memisahkan kata "Bearer"
	splitToken := strings.Split(authHeader, "Bearer ")
	if len(splitToken) != 2 {
		http.Error(w, "Invalid Authorization token format", http.StatusUnauthorized)
		return
	}
	tokenString := splitToken[1]

	// Validasi token dan ekstrak email dan companyID
	email, companyID, err := jwt.ValidateTokenOrSession(tokenString)
	if err != nil {
		http.Error(w, "Invalid token: "+err.Error(), http.StatusUnauthorized)
		return
	}

	// Jika email atau companyID digunakan untuk mengautentikasi percakapan
	// atau untuk menentukan izin pengguna, masukkan logika disini.
	// Contoh: mencatat informasi pengguna yang menyimpan percakapan
	log.Infof("User %s from company %d is saving a conversation", email, companyID)

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

	// Ambil nama pengguna dari database berdasarkan user_id
	var userName string
	err = s.DB.QueryRow("SELECT name FROM users WHERE id = ?", translationRequest.UserID).Scan(&userName)
	if err != nil {
		log.WithError(err).Error("Failed to retrieve user name")
		http.Error(w, "Failed to retrieve user name", http.StatusInternalServerError)
		return
	}

	// Jika tidak ada nama, gunakan default value atau handle sesuai kebutuhan
	if userName == "" {
		userName = "Unknown Speaker" // atau handle lainnya
	}

	// Tentukan bahasa berdasarkan speaker
	var japaneseText, englishText string
	if strings.EqualFold(translationRequest.Speaker, "ato") {
		japaneseText = translationRequest.OriginalMessage
		englishText = translatedMessage
	} else {
		englishText = translationRequest.OriginalMessage
		japaneseText = translatedMessage
	}

	var chatRoomID int
	if translationRequest.ChatRoomID == 0 {
		// Jika ChatRoomID tidak disediakan, coba tentukan atau buat chat room baru
		chatRoomID, err = s.determineOrCreateChatRoom(translationRequest.UserID, translationRequest.OtherUserID) // otherUserID harus ditentukan
		if err != nil {
			http.Error(w, "Failed to determine or create chat room", http.StatusInternalServerError)
			return
		}
	} else {
		chatRoomID = translationRequest.ChatRoomID
		// Opsi: Validasi apakah chat room yang diberikan valid
	}

	// Create Conversation object dengan speaker dari database
	t := models.Conversation{
		Speaker:           userName, // Gunakan userName sebagai Speaker
		JapaneseText:      japaneseText,
		EnglishText:       englishText,
		UserID:            translationRequest.UserID,
		CompanyID:         translationRequest.CompanyID,
		ChatRoomID:        chatRoomID,
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
				Speaker:           userName, // Gunakan userName di sini juga
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

func (s *Server) determineOrCreateChatRoom(user1ID, user2ID int) (int, error) {
	var chatRoomID int

	// Cek apakah chat room sudah ada
	query := `SELECT id FROM chat_room WHERE (user1_id = ? AND user2ID = ?) OR (user1_id = ? AND user2ID = ?)`
	err := s.DB.QueryRow(query, user1ID, user2ID, user2ID, user1ID).Scan(&chatRoomID)

	if err == sql.ErrNoRows {
		// Chat room tidak ada, buat chat room baru
		insertQuery := `INSERT INTO chat_room (user1_id, user2_id) VALUES (?, ?)`
		result, err := s.DB.Exec(insertQuery, user1ID, user2ID)
		if err != nil {
			return 0, err
		}

		// Dapatkan ID chat room yang baru dibuat
		newChatRoomID, err := result.LastInsertId()
		if err != nil {
			return 0, err
		}

		return int(newChatRoomID), nil
	} else if err != nil {
		// Terjadi error selain ErrNoRows
		return 0, err
	}

	// Chat room sudah ada, kembalikan ID-nya
	return chatRoomID, nil
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

// Start menjalankan server web
func (s *Server) Start(port string) {
	fmt.Printf("Server is running on port %s...\n", port)
	http.ListenAndServe(":"+port, s.Router)
}
