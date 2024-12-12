package api

import (
	"github.com/gorilla/mux"
	"github.com/z4fL/fp-ai-golang-neurons/service"
)

type API struct {
	token          string
	userService    service.UserService
	sessionService service.SessionService
	fileService    service.FileService
	aiService      service.AIService
}

func NewAPI(token string, userService service.UserService, sessionService service.SessionService, fileService service.FileService, aiService service.AIService) API {
	api := API{
		token,
		userService,
		sessionService,
		fileService,
		aiService,
	}

	return api
}

func RegisterRoutes(token string, router *mux.Router, userService service.UserService, sessionService service.SessionService, fileService service.FileService, aiService service.AIService) {
	api := NewAPI(token, userService, sessionService, fileService, aiService)
	router.HandleFunc("/register", api.Register).Methods("POST")
	router.HandleFunc("/login", api.Login).Methods("POST")
	router.HandleFunc("/logout/{id}", api.Logout).Methods("POST")
	router.HandleFunc("/upload", api.Upload).Methods("POST")
	router.HandleFunc("/chat", api.Chat).Methods("POST")
	router.HandleFunc("/remove-session", api.RemoveSession).Methods("POST")
}
