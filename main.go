package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
	"github.com/z4fL/fp-ai-golang-neurons/api"
	"github.com/z4fL/fp-ai-golang-neurons/db"
	"github.com/z4fL/fp-ai-golang-neurons/model"
	"github.com/z4fL/fp-ai-golang-neurons/repository"
	"github.com/z4fL/fp-ai-golang-neurons/service"
	"github.com/z4fL/fp-ai-golang-neurons/utility"
)

const defaultPort = "8080"

func main() {
	// Load the .env file
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	db := db.NewDB()
	dbCredential, err := utility.GetDBCredential()
	if err != nil {
		log.Fatalf("Error getting DB credentials: %v", err)
	}

	conn, err := db.Connect(dbCredential)
	if err != nil {
		panic(err)
	}

	conn.AutoMigrate(&model.User{}, &model.Session{})

	// Retrieve the Hugging Face token from the environment variables
	token := os.Getenv("HUGGINGFACE_TOKEN")
	if token == "" {
		log.Fatal("Environment variable HUGGINGFACE_TOKEN isn't set in the .env file")
	}

	userRepo := repository.NewUserRepository(conn)
	sessionRepo := repository.NewSessionRepo(conn)
	fileRepo := repository.NewFileRepository()

	userService := service.NewUserService(userRepo)
	sessionService := service.NewSessionService(sessionRepo)
	fileService := service.NewFileService(fileRepo)
	aiService := service.NewAIService(&http.Client{})

	// Set up the router
	router := mux.NewRouter()
	api.RegisterRoutes(token, router, userService, sessionService, fileService, aiService)

	// Enable CORS
	corsHandler := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:5173"},
		// AllowedOrigins: []string{"*"}, // Allow all origins
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Content-Type", "Authorization"},
	}).Handler(router)

	// List all routes
	router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		pathTemplate, err := route.GetPathTemplate()
		if err != nil {
			return err
		}
		methods, err := route.GetMethods()
		if err != nil {
			return err
		}
		log.Printf("Route: %s %s", methods, pathTemplate)
		return nil
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	// Start the server
	log.Printf("Server running on port %s\n", port)
	if err := http.ListenAndServe(":"+port, corsHandler); err != nil {
		log.Fatalf("Server failed to start on port %s: %v", port, err)
	}
}
