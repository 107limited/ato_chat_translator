package web

// initializeRoutes mengatur rute-rute untuk server
func (s *Server) initializeRoutes() {
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
	
}
