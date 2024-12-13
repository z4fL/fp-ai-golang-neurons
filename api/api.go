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
	chatService    service.ChatService
}

func NewAPI(token string, userService service.UserService, sessionService service.SessionService, fileService service.FileService, aiService service.AIService, chatService service.ChatService) API {
	api := API{
		token,
		userService,
		sessionService,
		fileService,
		aiService,
		chatService,
	}

	return api
}

func RegisterRoutes(token string, router *mux.Router, userService service.UserService, sessionService service.SessionService, fileService service.FileService, aiService service.AIService, chatService service.ChatService) {
	api := NewAPI(token, userService, sessionService, fileService, aiService, chatService)

	router.HandleFunc("/register", api.Register).Methods("POST")
	router.HandleFunc("/login", api.Login).Methods("POST")
	router.HandleFunc("/logout/{id}", api.Logout).Methods("POST")

	router.HandleFunc("/upload", api.Upload).Methods("POST")
	router.HandleFunc("/chat-with-ai", api.ChatWithAI).Methods("POST")

	router.HandleFunc("/chats", api.CreateChat).Methods("POST")
	router.HandleFunc("/chats/{userID}", api.AddMessage).Methods("PUT")
	router.HandleFunc("/remove-session", api.RemoveSession).Methods("POST")
}
