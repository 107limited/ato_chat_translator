package main

import (
	"ato_chat/chat"
	"ato_chat/config"
	"ato_chat/translation"
	"ato_chat/web"
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	
	
)

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

	// Create CORS middleware
    corsMiddleware := CORSMiddleware(server.Router)

    // Attach the CORS middleware before your routes
    http.Handle("/", corsMiddleware)

	// Create a context that listens for the interrupt signal from the OS
	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Get the server port from the environment or .env file
	port := os.Getenv("PORT_SERVER")
	if port == "" {
		port = "8080" // Port default jika tidak ditemukan
	}

	// Start HTTP server
	fmt.Printf("Server is running on port %s...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, server.Router))
}

// CORSMiddleware menangani kebijakan CORS
func CORSMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		

        if r.Method == http.MethodOptions {
            // Handling preflight request
            w.WriteHeader(http.StatusOK)
            return
        }

        next.ServeHTTP(w, r)
    })
}
