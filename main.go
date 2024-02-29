package main

import (
	"ato_chat/chat"
	"ato_chat/config"
	"ato_chat/translation"
	"ato_chat/web"
	"ato_chat/websocket"
	"database/sql"
	"fmt"

	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/handlers"
	log "github.com/sirupsen/logrus"
)

func init() {
	// Konfigurasi logrus
	log.SetFormatter(&log.TextFormatter{
		ForceColors:   true, // Mengaktifkan warna
		FullTimestamp: true, // Menampilkan timestamp lengkap
	})

	// Jika Anda ingin level log ditampilkan dalam huruf kapital
	log.SetLevel(log.InfoLevel)
	log.SetReportCaller(true) // Jika Anda ingin melihat di mana log dipanggil
}

func main() {
	// Load database configuration from .env
	dbConfig := config.LoadDBConfig()
	log.Printf("Database Config: %+v\n", dbConfig)

	// Get the database connection string
	dbConnectionString := dbConfig.GetDBConnectionString()

	// Attempt to open a connection to the database
	db, err := sql.Open("mysql", dbConnectionString)
	if err != nil {
		log.Fatalf("Error opening database: %v\n", err)
	}
	defer db.Close()

	// Attempt to ping the database to check the connection
	err = db.Ping()
	if err != nil {
		log.Fatalf("Error connecting to database: %v\n", err)
	}
	fmt.Println("Database connection successful!")

	// Load API key from environment variable or .env file
	apiKey := os.Getenv("OPENAI_API_KEY")

	// Create GPT-4 Translator
	gpt4Translator := &translation.GPT4Client{APIKey: apiKey, Model: "gpt-3.5-turbo-0613", DB: chat.NewConversationRepository(db)}

	// Create MySQL Conversation Repository with GPT-4 Translator
	conversationRepo := chat.ConversationRepository(gpt4Translator)

	chatRoomHandler := web.NewChatRoomHandler(db)

	// Assuming chat.NewConversationRepository and websocket.NewConnectionManager are correctly implemented
    repo := chat.NewConversationRepository(db) // Creates a new instance of the conversation repository
    cm := websocket.NewConnectionManager()     // Creates a new instance of the connection manager

    // Initialize ConversationService with the repo and connection manager
    //cs := websocket.NewConversationService(repo, cm)
	cs := websocket.NewConversationService(repo, cm)
	

	// Create HTTP server
	server := web.NewServer(db, conversationRepo, gpt4Translator, chatRoomHandler, cs)
	//server.ConnectionManager = websocket.NewConnectionManager() // Inisialisasi ConnectionManager di sini

	// Set log format as text formatter with full timestamp
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})

	
	
	// conversationRepo = chat.NewConversationRepository(db)
    // connectionManager := websocket.NewConnectionManager()
    // conversationService := websocket.NewConversationService(conversationRepo, cm)
    server.ConnectionManager = cm
    //wsHandler := websocket.HandleWebSocket(conversationService)
    
    // // Assuming `server.Router` is correctly set up elsewhere:
    server.Router.HandleFunc("/ws", func (w http.ResponseWriter, r *http.Request) {
	    websocket.HandleWSL(w,r)
    })

	// Get the server port from the environment or .env file
	port := os.Getenv("PORT_SERVER")
	if port == "" {
		port = "8080" // Port default jika tidak ditemukan
	}
	log.Infof("Server is running on port %s...", port)

	log.Fatal(http.ListenAndServe(":"+port, handlers.CORS(
		handlers.AllowedOrigins([]string{"http://localhost:5173", "https://ato-puce.vercel.app", "https://chat-ato.vercel.app", "https://t-chat.107.jp"}),
		handlers.AllowedMethods([]string{"GET", "POST", "DELETE", "PUT", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}),
	)(server.Router)))
}