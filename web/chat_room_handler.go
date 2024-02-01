package web

import (
	"database/sql"
	"encoding/json"

	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/labstack/gommon/log"
)

// ChatRoomHandler adalah struct yang menangani request terkait chat room.
type ChatRoomHandler struct {
    DB *sql.DB
}

// NewChatRoomHandler adalah konstruktor untuk membuat instance baru dari ChatRoomHandler.
func NewChatRoomHandler(db *sql.DB) *ChatRoomHandler {
    return &ChatRoomHandler{DB: db}
}

// CreateChatRoom handles the creation of a new chat room
func (h *ChatRoomHandler) CreateChatRoom(w http.ResponseWriter, r *http.Request) {
	// Extract user IDs from the request
	var req struct {
		User1ID int `json:"user1_id"`
		User2ID int `json:"user2_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validasi keberadaan user1_id dan user2_id di tabel users
	if !h.isUserExists(req.User1ID) || !h.isUserExists(req.User2ID) {
		http.Error(w, "One or both users not found", http.StatusBadRequest)
		return
	}

	// Create a new chat room in the database
	var chatRoomID int64
	query := `INSERT INTO chat_room (user1_id, user2_id) VALUES (?, ?)`
	result, err := h.DB.Exec(query, req.User1ID, req.User2ID)
	if err != nil {
		http.Error(w, "Failed to create chat room", http.StatusInternalServerError)
		log.Error("Failed to execute query: ", err)
		return
	}

	// Get the ID of the newly created chat room
	chatRoomID, err = result.LastInsertId()
	if err != nil {
		http.Error(w, "Failed to retrieve chat room ID", http.StatusInternalServerError)
		log.Error("Failed to retrieve last insert ID: ", err)
		return
	}

	// Respond with the created chat room details
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int64{"chat_room_id": chatRoomID})
}

// isUserExists checks if a user exists in the users table
func (h *ChatRoomHandler) isUserExists(userID int) bool {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE id = ?)`
	err := h.DB.QueryRow(query, userID).Scan(&exists)
	if err != nil {
		log.Error("Failed to execute query: ", err)
		return false
	}
	return exists
}

// GetChatRoom retrieves an existing chat room between two users
func (h *ChatRoomHandler) GetChatRoom(w http.ResponseWriter, r *http.Request) {
	user1ID := mux.Vars(r)["user1_id"]
	user2ID := mux.Vars(r)["user2_id"]

	// Convert user IDs to int
	u1ID, err := strconv.Atoi(user1ID)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}
	u2ID, err := strconv.Atoi(user2ID)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Query the database for the chat room
	var chatRoomID int
	query := `SELECT id FROM chat_room WHERE (user1_id = ? AND user2_id = ?) OR (user1_id = ? AND user2_id = ?)`
	err = h.DB.QueryRow(query, u1ID, u2ID, u2ID, u1ID).Scan(&chatRoomID)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "No chat room found", http.StatusNotFound)
		} else {
			http.Error(w, "Error querying chat room", http.StatusInternalServerError)
			log.Error("Failed to execute query: ", err)
		}
		return
	}

	// Respond with the chat room details
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int{"chat_room_id": chatRoomID})
}

// SaveConversation handles saving a conversation to a chat room
// func (h *ChatRoomHandler) SaveConversation(w http.ResponseWriter, r *http.Request) {
// 	// Parse JSON request body untuk mendapatkan detail percakapan
// 	var conversationRequest models.ConversationRequest
// 	if err := json.NewDecoder(r.Body).Decode(&conversationRequest); err != nil {
// 		http.Error(w, "Failed to parse request body", http.StatusBadRequest)
// 		return
// 	}

// 	// Validasi chat_room_id (opsional)
// 	if !h.isChatRoomExists(conversationRequest.ChatRoomID) {
// 		http.Error(w, "Chat room does not exist", http.StatusBadRequest)
// 		return
// 	}

// 	// Translate original message jika diperlukan
// 	translatedMessage, err := h.GPT4Translator.TranslateMessage(conversationRequest.OriginalMessage)
// 	if err != nil {
// 		http.Error(w, fmt.Sprintf("Failed to translate message: %v", err), http.StatusInternalServerError)
// 		return
// 	}

// 	// Buat objek Conversation
// 	conversation := models.Conversation{
// 		Speaker:           conversationRequest.Speaker,
// 		JapaneseText:      "", // Atur sesuai dengan kebutuhan
// 		EnglishText:       translatedMessage,
// 		UserID:            conversationRequest.UserID,
// 		CompanyID:         conversationRequest.CompanyID,
// 		ChatRoomID:        conversationRequest.ChatRoomID,
// 		OriginalMessage:   conversationRequest.OriginalMessage,
// 		TranslatedMessage: translatedMessage,
// 		CreatedAt:         time.Now(),
// 	}

// 	// Simpan percakapan ke dalam database
// 	err = h.ConversationRepo.SaveConversation(&conversation)
// 	if err != nil {
// 		http.Error(w, fmt.Sprintf("Failed to save conversation: %v", err), http.StatusInternalServerError)
// 		return
// 	}

// 	// Kirim respons sukses
// 	w.WriteHeader(http.StatusCreated)
// 	json.NewEncoder(w).Encode(map[string]string{"message": "Conversation saved successfully"})
// }

// // isChatRoomExists checks if a chat room exists in the database
// func (h *ChatRoomHandler) isChatRoomExists(chatRoomID int) bool {
// 	var exists bool
// 	query := `SELECT EXISTS(SELECT 1 FROM chat_room WHERE id = ?)`
// 	err := h.DB.QueryRow(query, chatRoomID).Scan(&exists)
// 	if err != nil {
// 		log.Error("Failed to execute query: ", err)
// 		return false
// 	}
// 	return exists
// }
