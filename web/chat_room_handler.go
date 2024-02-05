package web

import (
	"ato_chat/dbAto"
	"ato_chat/models"
	"database/sql"
	"encoding/json"
	"fmt"

	"time"

	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	//"github.com/labstack/gommon/log"
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetReportCaller(true)
	log.SetLevel(log.DebugLevel) // or whatever level you need
}

// ChatRoomHandler adalah struct yang menangani request terkait chat room.
type ChatRoomHandler struct {
	DB *sql.DB
}

// NewChatRoomHandler adalah konstruktor untuk membuat instance baru dari ChatRoomHandler.
func NewChatRoomHandler(db *sql.DB) *ChatRoomHandler {
	return &ChatRoomHandler{DB: db}
}
func (h *ChatRoomHandler) GetExistingChatRoomId(user1ID, user2ID int) (int64, bool) {
	var chatRoomID int64
	query := `SELECT id FROM chat_room WHERE (user1_id = ? AND user2_id = ?) OR (user1_id = ? AND user2_id = ?) LIMIT 1`
	err := h.DB.QueryRow(query, user1ID, user2ID, user2ID, user1ID).Scan(&chatRoomID)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, false // No existing chat room found
		}
		log.Error("Failed to check if chat room exists: ", err)
		return 0, false
	}
	return chatRoomID, true
}

func (h *ChatRoomHandler) isUserExists(userID int) bool {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE id = ?)`
	err := h.DB.QueryRow(query, userID).Scan(&exists)
	if err != nil {
		log.Printf("Failed to check if user exists: %v", err)
		return false
	}
	return exists
}

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

	// Cek apakah chat room antara dua user sudah ada dan dapatkan ID-nya jika ada
	chatRoomID, exists := h.GetExistingChatRoomId(req.User1ID, req.User2ID)
	if exists {
		resp := models.ChatRoomResponse{
			Message:    "Chat room between these users already exists",
			ChatRoomID: chatRoomID,
		}
		w.WriteHeader(http.StatusBadRequest) // Atau gunakan http.StatusConflict jika lebih sesuai
		json.NewEncoder(w).Encode(resp)
		return
	}

	// Jika tidak ada, buat chat room baru
	query := `INSERT INTO chat_room (user1_id, user2_id) VALUES (?, ?)`
	result, err := h.DB.Exec(query, req.User1ID, req.User2ID)
	if err != nil {
		http.Error(w, "Failed to create chat room", http.StatusInternalServerError)
		log.Error("Failed to execute query: ", err)
		return
	}

	chatRoomID, err = result.LastInsertId()
	if err != nil {
		http.Error(w, "Failed to retrieve chat room ID", http.StatusInternalServerError)
		log.Error("Failed to retrieve last insert ID: ", err)
		return
	}

	// Respond dengan detail chat room yang baru dibuat
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int64{"chat_room_id": chatRoomID})
}

// GetAllChatRoomsHandler handles the request to get all chat rooms.
func (s *Server) GetAllChatRoomsHandler(w http.ResponseWriter, r *http.Request) {
	chatRooms, err := dbAto.GetAllChatRooms(s.DB)
	if err != nil {
		http.Error(w, "Failed to get chat rooms", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(chatRooms)
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

// GetConversationsByChatRoom mengambil percakapan dari database berdasarkan chat_room_id.
func GetConversationsByChatRoom(db *sql.DB, chatRoomID int) ([]models.Conversation, error) {
	var conversations []models.Conversation

	query := `SELECT id, japanese_text, english_text, user_id, company_id, chat_room_id, created_at, date, speaker FROM conversations WHERE chat_room_id = ?`
	rows, err := db.Query(query, chatRoomID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var conv models.Conversation
		var createdAt string // Menggunakan string untuk menampung waktu
		if err := rows.Scan(&conv.ID, &conv.JapaneseText, &conv.EnglishText, &conv.UserID, &conv.CompanyID, &conv.ChatRoomID, &createdAt, &conv.Date, &conv.Speaker); err != nil {
			return nil, err
		}

		// Menggunakan format yang sesuai dengan output dari database Anda
		conv.CreatedAt, err = time.Parse("2006-01-02 15:04:05", createdAt)
		if err != nil {
			return nil, err
		}
		conversations = append(conversations, conv)
	}

	return conversations, nil
}

// GetConversationsByChatRoomHandler menangani permintaan HTTP untuk mengambil percakapan berdasarkan chat_room_id.
func (h *ChatRoomHandler) GetConversationsByChatRoomHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	chatRoomIDStr := vars["chat_room_id"]

	chatRoomID, err := strconv.Atoi(chatRoomIDStr)
	if err != nil {
		http.Error(w, "Invalid chat room ID", http.StatusBadRequest)
		return
	}

	conversations, err := GetConversationsByChatRoom(h.DB, chatRoomID)
	if err != nil {
		http.Error(w, "Failed to get conversations: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(conversations)
}

// GetChatRoomByIdHandler handles requests to get a chat room by ID.
func (s *Server) GetChatRoomByIdHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	chatRoomIdStr, exists := vars["id"]
	if !exists {
		http.Error(w, "Chat room ID is required", http.StatusBadRequest)
		return
	}

	chatRoomId, err := strconv.Atoi(chatRoomIdStr)
	if err != nil {
		http.Error(w, "Invalid chat room ID", http.StatusBadRequest)
		return
	}

	chatRoom, err := dbAto.GetChatRoomById(s.DB, chatRoomId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound) // Use http.StatusNotFound if chat room does not exist.
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(chatRoom)
}

// GetChatRoomsByUserIDHandler handles the HTTP request to get chat rooms by user ID.
func (s *Server) GetChatRoomsByUserIDHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userIDStr, ok := vars["user_id"]
	if !ok {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	fmt.Printf("Received user ID string: %s\n", userIDStr) // Tambahkan untuk debug

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	chatRooms, err := dbAto.GetChatRoomsByUserID(s.DB, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(chatRooms)
}
