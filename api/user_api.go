package api

import (
	"encoding/json"
	"net/http"
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
	session := model.Session{Token: sessionToken, UserID: credentials.ID, Expiry: expiresAt}

	err = api.sessionService.SessionAvailID(session.UserID)
	if err != nil {
		err = api.sessionService.AddSession(session)
	} else {
		err = api.sessionService.UpdateSession(session)
	}

	if err != nil {
		utility.JSONResponse(w, http.StatusInternalServerError, "failed", "Internal Server Error")
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Path:     "/",
		Value:    sessionToken,
		Expires:  expiresAt,
		HttpOnly: true,                  // Supaya cookie nggak bisa diakses via JavaScript
		SameSite: http.SameSiteNoneMode, // Cookie bisa dikirim cross-origin
		Secure:   false,                 // Set ke true kalau pakai HTTPS
	})

	utility.JSONResponse(w, http.StatusOK, "success", "user "+user.Username+" loggin successfully")
}

func (api *API) Logout(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("session_token")
	if err != nil {
		if err == http.ErrNoCookie {
			utility.JSONResponse(w, http.StatusUnauthorized, "failed", "Internal Server Error")
			return
		}

		utility.JSONResponse(w, http.StatusBadRequest, "failed", "Internal Server Error")
		return
	}
	sessionToken := c.Value

	api.sessionService.DeleteSession(sessionToken)

	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   "",
		Expires: time.Now(),
	})

	utility.JSONResponse(w, http.StatusOK, "success", "logout successfully")
}
