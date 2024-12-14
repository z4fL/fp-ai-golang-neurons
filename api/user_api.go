package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/z4fL/fp-ai-golang-neurons/model"
	"github.com/z4fL/fp-ai-golang-neurons/utility"
)

func (api *API) Register(w http.ResponseWriter, r *http.Request) {
	var user model.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		utility.JSONResponse(w, http.StatusBadRequest, "failed", "Invalid request payload")
		return
	}

	if err := api.userService.Register(user); err != nil {
		utility.JSONResponse(w, http.StatusInternalServerError, "failed", "Failed to register user")
		return
	}

	utility.JSONResponse(w, http.StatusCreated, "success", "User registered successfully")
}

func (api *API) Login(w http.ResponseWriter, r *http.Request) {
	var credentials model.User

	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		utility.JSONResponse(w, http.StatusBadRequest, "failed", "Invalid request payload")
		return
	}

	user, err := api.userService.Login(credentials.Username, credentials.Password)
	if err != nil {
		utility.JSONResponse(w, http.StatusUnauthorized, "failed", "Invalid username or password")
		return
	}

	sessionToken := uuid.NewString()
	expiresAt := time.Now().Add(5 * time.Hour)
	session := model.Session{Token: sessionToken, UserID: user.ID, Expiry: expiresAt}

	err = api.sessionService.SessionAvailID(session.UserID)
	if err != nil {
		log.Print("Add session")
		err = api.sessionService.AddSession(session)
	} else {
		log.Print("Update session")
		err = api.sessionService.UpdateSession(session)
	}

	if err != nil {
		utility.JSONResponse(w, http.StatusInternalServerError, "failed", "Internal Server Error")
		return
	}

	utility.JSONResponse(w, http.StatusOK, "success", sessionToken)
}

func (api *API) Logout(w http.ResponseWriter, r *http.Request) {
	// Ambil token dari header Authorization
	authHeader := r.Header.Get("Authorization")
	if !strings.HasPrefix(authHeader, "Bearer ") {
		utility.JSONResponse(w, http.StatusUnauthorized, "failed", "Missing or invalid token")
		return
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")

	// Validasi token
	_, err := api.sessionService.SessionAvailToken(token)
	if err != nil {
		utility.JSONResponse(w, http.StatusUnauthorized, "failed", "Invalid token")
		return
	}

	err = api.sessionService.DeleteSession(token)
	if err != nil {
		utility.JSONResponse(w, http.StatusInternalServerError, "failed", "Failed to logout")
		return
	}

	utility.JSONResponse(w, http.StatusOK, "success", "logout successfully")
}

func (api *API) ValidateSession(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if !strings.HasPrefix(authHeader, "Bearer ") {
		http.Error(w, "Unauthorized: Missing or invalid token", http.StatusUnauthorized)
		return
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")

	session, err := api.sessionService.SessionAvailToken(token)
	if err != nil || api.sessionService.TokenExpired(session) {
		if err == nil {
			api.sessionService.DeleteSession(token) // Hapus token expired
			log.Print("successfully delete token")
		}
		http.Error(w, "Unauthorized: Invalid or expired token", http.StatusUnauthorized)
		return
	}

	utility.JSONResponse(w, http.StatusOK, "success", "Token is valid")
}
