package web

// initializeRoutes mengatur rute-rute untuk server
func (s *Server) initializeRoutes() {
    // Rute yang sudah ada
    s.Router.HandleFunc("/api/conversations", s.SaveConversationHandler).Methods("POST")
    s.Router.HandleFunc("/api/conversations", s.GetAllConversationsHandler).Methods("GET")
    s.Router.HandleFunc("/api/translate", s.TranslateMessageHandler).Methods("POST")
    s.Router.HandleFunc("/api/register", s.RegisterUserHandler).Methods("POST")

    // Tambahkan endpoint untuk Register dan Personal data
    s.Router.HandleFunc("/api/login", s.LoginUserHandler).Methods("POST")
    s.Router.HandleFunc("/api/personaldata", s.PersonalDataHandler).Methods("POST")

    // Tambahkan endpoint untuk login
    s.Router.HandleFunc("/api/login", s.LoginUserHandler).Methods("POST")

    // Tambahkan rute untuk ChatRoomHandler
    s.Router.HandleFunc("/api/chatrooms", s.ChatRoomHandler.CreateChatRoom).Methods("POST")
    s.Router.HandleFunc("/api/chatrooms/{user1_id}/{user2_id}", s.ChatRoomHandler.GetChatRoom).Methods("GET")

    // Tambahkan endpoint untuk mendapatkan semua pengguna
    s.Router.HandleFunc("/api/users", s.GetAllUsersHandler).Methods("GET")

    // Tambahkan rute untuk GetAllRolesHandler
    s.Router.HandleFunc("/api/roles", s.GetAllRolesHandler).Methods("GET")
}
