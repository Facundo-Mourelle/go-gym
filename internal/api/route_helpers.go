package api

import (
	"net/http"

	"github.com/Facundo-Mourelle/go-gym/internal/api/middleware"
	"github.com/Facundo-Mourelle/go-gym/internal/service"
)

// Protected wraps a handler with auth middleware
func Protected(authService *service.AuthService, handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		middleware.AuthMiddleware(authService)(handler).ServeHTTP(w, r)
	}
}
