// middleware.go

package web

import (
	"github.com/rs/cors"
	"net/http"
)

// CORSMiddleware menangani kebijakan CORS
func CORSMiddleware(next http.Handler) http.Handler {
	corsHandler := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:5173/"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders: []string{"Content-Type", "Authorization"},
		ExposedHeaders: []string{"Access-Control-Allow-Origin"},
	}).Handler

	return corsHandler(next)
}
