package web

import (
	// "ato_chat/websocket"
	// "net/http"
)

//import (
//"ato_chat/websocket"

//"net/http"
//"ato_chat/websocket"

//"net/http"
//)

// initializeRoutes mengatur rute-rute untuk server
func (s *Server) initializeRoutes() {

	// Pastikan chatRepo dan connManager sudah diinisialisasi
	//conversationService := websocket.NewConversationService(chatRepo, connManager)
	// Setup WebSocket route
	// s.Router.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
	// 	websocket.HandleWebSocket(s.ConversationService)(w, r)
	// }).Methods("GET")

	// // Rute yang sudah ada
	s.Router.HandleFunc("/api/conversations", s.SaveConversationHandler).Methods("POST")
	s.Router.HandleFunc("/api/conversations", s.GetAllConversationsHandler).Methods("GET")
	s.Router.HandleFunc("/api/translate", s.TranslateMessageHandler).Methods("POST")

	// Tambahkan endpoint untuk Register dan Personal data
	s.Router.HandleFunc("/api/register", s.RegisterUserHandler).Methods("POST")
	s.Router.HandleFunc("/api/personaldata", s.PersonalDataHandler).Methods("POST")
	// Tambahkan endpoint untuk login
	s.Router.HandleFunc("/api/login", s.LoginUserHandler).Methods("POST")
	// Tambahkan endpoint untuk mendapatkan semua pengguna
	s.Router.HandleFunc("/api/users", s.GetAllUsersHandler).Methods("GET")
	// Get User By Id
	s.Router.HandleFunc("/api/user/{id}", s.GetUserByIdHandler).Methods("GET")
	// Get User By Company ID
	s.Router.HandleFunc("/api/users/company/{companyId}", s.GetUsersByCompanyIdHandler).Methods("GET")
	s.Router.HandleFunc("/api/users/{companyIdentifier}", s.GetUsersByCompanyIdentifierHandler).Methods("GET")
	// Tambahkan rute untuk Logout
	s.Router.HandleFunc("/api/logout", s.LogoutHandler).Methods("POST")

	// // Tambahkan rute untuk GetAllRolesHandler
	s.Router.HandleFunc("/api/roles", s.CreateRoleHandler).Methods("POST")
	s.Router.HandleFunc("/api/roles", s.GetAllRolesHandler).Methods("GET")

	// // Tambahkan rute untuk ChatRoomHandler
	s.Router.HandleFunc("/api/chatrooms", s.ChatRoomHandler.CreateChatRoom).Methods("POST")
	// // Get All chat room
	s.Router.HandleFunc("/api/chatrooms", s.GetAllChatRoomsHandler).Methods("GET")
	// // Get chat room by User id
	s.Router.HandleFunc("/api/chatrooms/user/{user_id}", s.GetChatRoomsByUserIDHandler).Methods("GET")
	// // Get chat room by id
	s.Router.HandleFunc("/api/chatrooms/{id}", s.GetChatRoomByIdHandler).Methods("GET")
	// // Tambahkan rute untuk mendapatkan percakapan berdasarkan chat_room_id.
	s.Router.HandleFunc("/api/conversations-by-chat-room-id/{chat_room_id}", s.ChatRoomHandler.GetConversationsByChatRoomHandler).Methods("GET")
	// GetChatRoom retrieves an existing chat room between two users
	s.Router.HandleFunc("/api/chatrooms/{user1_id}/{user2_id}", s.ChatRoomHandler.GetChatRoom).Methods("GET")
	// Di dalam fungsi setup router Anda, tambahkan:
	s.Router.HandleFunc("/ws/chatrooms/{user_id}", s.ChatRoomsWebSocketHandler)

}
