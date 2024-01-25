package main

import (
	"ato_chat/chat"
	"ato_chat/config"
	"ato_chat/translation"
	"ato_chat/web"
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
        ForceColors:   true,       // Mengaktifkan warna
        FullTimestamp: true,       // Menampilkan timestamp lengkap
    })

    // Jika Anda ingin level log ditampilkan dalam huruf kapital
    log.SetLevel(log.InfoLevel)
    log.SetReportCaller(true)    // Jika Anda ingin melihat di mana log dipanggil
}

func main() {
	// Load database configuration from .env
	dbConfig := config.LoadDBConfig()

	// Get database connection string
	dbConnectionString := dbConfig.GetDBConnectionString()

	// Attempt to open a connection to the database
	db, err := sql.Open("mysql", dbConnectionString)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Attempt to ping the database to check the connection
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Database connection successful!")

	// Load API key from environment variable or .env file
	apiKey := os.Getenv("OPENAI_API_KEY")

	// Create GPT-4 Translator
	gpt4Translator := &translation.GPT4Client{APIKey: apiKey, Model: "gpt-3.5-turbo-0613", DB: chat.NewConversationRepository(db)}

	// Create MySQL Conversation Repository with GPT-4 Translator
	conversationRepo := chat.ConversationRepository(gpt4Translator)

	// Create HTTP server
	server := web.NewServer(conversationRepo, gpt4Translator)

	
	
	// Set log format as text formatter with full timestamp
    log.SetFormatter(&log.TextFormatter{
        FullTimestamp: true,
    })

    // Get the server port from the environment or .env file
    port := os.Getenv("PORT_SERVER")
    if port == "" {
        port = "8080" // Port default jika tidak ditemukan
    }
    log.Infof("Server is running on port %s...", port)

	log.Fatal(http.ListenAndServe(":"+port,handlers.CORS(
		handlers.AllowedOrigins([]string{"http://localhost:5173", "https://ato-puce.vercel.app"}),
		handlers.AllowedMethods([]string{"GET","POST","DELETE","PUT","OPTIONS"}),
		handlers.AllowedHeaders([]string{"X-Requested-With","Content-Type","Authorization"}),
	)(server.Router)))
}
