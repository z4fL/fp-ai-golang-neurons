package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/z4fL/fp-ai-golang-neurons/service"
	"github.com/z4fL/fp-ai-golang-neurons/utility"
)

type ContextKey string

const UserIDKey ContextKey = "userID"

func AuthMiddleware(sessionService service.SessionService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Ambil token dari header Authorization
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
				utility.JSONResponse(w, http.StatusUnauthorized, "failed", "Missing or invalid token")
				return
			}

			// Ambil token tanpa prefix "Bearer "
			token := strings.TrimPrefix(authHeader, "Bearer ")

			// cek session berdasarkan token
			session, err := sessionService.TokenValidity(token)
			if err != nil {
				utility.JSONResponse(w, http.StatusUnauthorized, "failed", "Token is Expired")
				return
			}

			// Masukkan userID ke context
			ctx := context.WithValue(r.Context(), UserIDKey, session.UserID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
