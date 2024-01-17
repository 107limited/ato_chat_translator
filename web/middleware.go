package web

import (
	"net/http"
	"github.com/gorilla/handlers"
)

// CORSMiddleware menangani kebijakan CORS
func CORSMiddleware(next http.Handler) http.Handler {
	return handlers.CORS(
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"}),
		handlers.AllowedOrigins([]string{"*"}),
	)(next)
}
