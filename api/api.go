package api

import (
	"github.com/gorilla/mux"
	"github.com/z4fL/fp-ai-golang-neurons/middleware"
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

	authMiddleware := middleware.AuthMiddleware(sessionService)
	securedRoutes := router.PathPrefix("/").Subrouter()
	securedRoutes.Use(authMiddleware)

	router.HandleFunc("/register", api.Register).Methods("POST")
	router.HandleFunc("/login", api.Login).Methods("POST")

	securedRoutes.HandleFunc("/logout", api.Logout).Methods("POST")

	securedRoutes.HandleFunc("/upload", api.Upload).Methods("POST")
	securedRoutes.HandleFunc("/chat-with-ai", api.ChatWithAI).Methods("POST")

	securedRoutes.HandleFunc("/chats", api.CreateChat).Methods("POST")
	securedRoutes.HandleFunc("/chats", api.AddMessage).Methods("PATCH")
	securedRoutes.HandleFunc("/remove-session", api.RemoveSession).Methods("POST")
}
