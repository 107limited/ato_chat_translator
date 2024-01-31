package web

// initializeRoutes mengatur rute-rute untuk server
func (s *Server) initializeRoutes() {
    // Rute yang sudah ada
    s.Router.HandleFunc("/api/conversations", s.SaveConversationHandler).Methods("POST")
    s.Router.HandleFunc("/api/conversations", s.GetAllConversationsHandler).Methods("GET")
    s.Router.HandleFunc("/api/translate", s.TranslateMessageHandler).Methods("POST")
    s.Router.HandleFunc("/api/register", s.RegisterUserHandler).Methods("POST")
    s.Router.HandleFunc("/api/login", s.LoginUserHandler).Methods("POST")

    // Tambahkan rute untuk PersonalDataHandler
    s.Router.HandleFunc("/api/personaldata", s.PersonalDataHandler).Methods("POST")
}
