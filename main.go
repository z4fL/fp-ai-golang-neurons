package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
	"github.com/z4fL/fp-ai-golang-neurons/handler"
)

const defaultPort = "8080"

func main() {
	// Load the .env file
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Retrieve the Hugging Face token from the environment variables
	handler.Token = os.Getenv("HUGGINGFACE_TOKEN")
	if handler.Token == "" {
		log.Fatal("HUGGINGFACE_TOKEN isn't set in the .env file")
	}

	// Set up the router
	router := mux.NewRouter()

	// File upload endpoint
	router.HandleFunc("/upload", handler.HandleUpload).Methods("POST")

	// Chat endpoint
	router.HandleFunc("/chat", handler.HandleChat).Methods("POST")

	// Endpoint to remove the session file (data-series.csv) from the server.
	router.HandleFunc("/remove-session", handler.HandleRemoveSession).Methods("POST")

	// Enable CORS
	corsHandler := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:5173"},
		// AllowedOrigins: []string{"*"}, // Allow all origins
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Content-Type", "Authorization"},
	}).Handler(router)

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	// Start the server
	log.Printf("Server running on port %s\n", port)
	if err := http.ListenAndServe(":"+port, corsHandler); err != nil {
		log.Printf("Server failed to start: %v", err)
	}
}
