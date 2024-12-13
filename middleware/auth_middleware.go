package middleware

import (
	"context"
	"net/http"

	"github.com/z4fL/fp-ai-golang-neurons/service"
)

type ContextKey string

const UserIDKey ContextKey = "userID"

func AuthMiddleware(sessionService service.SessionService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Ambil session_token dari cookie
			cookie, err := r.Cookie("session_token")
			if err != nil {
				http.Error(w, "Unauthorized: Missing session token", http.StatusUnauthorized)
				return
			}

			// Validasi session_token dan ambil userID
			userID, err := sessionService.GetUserIDByToken(cookie.Value)
			if err != nil {
				http.Error(w, "Unauthorized: Invalid session token", http.StatusUnauthorized)
				return
			}

			// Masukkan userID ke context
			ctx := context.WithValue(r.Context(), UserIDKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
