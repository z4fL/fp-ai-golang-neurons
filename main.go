package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
	"github.com/z4fL/fp-ai-golang-neurons/service"
)

// Init services
var fileService = &service.FileService{}
var aiService = &service.AIService{}

func main() {
	// Load the .env file
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Retrieve the Hugging Face token from the environment variables
	token := os.Getenv("HUGGINGFACE_TOKEN")
	if token == "" {
		log.Fatal("HUGGINGFACE_TOKEN isn't set in the .env file")
	}

	// Set up the router
	router := mux.NewRouter()

	// File upload endpoint
	router.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		// Parse form data
		err := r.ParseMultipartForm(1 << 20) // 1MB
		if err != nil {
			http.Error(w, "Unable to parse form", http.StatusBadRequest)
			return
		}

		// Get the uploaded file
		file, handler, err := r.FormFile("file") // "file" sesuai nama field di frontend
		if err != nil {
			http.Error(w, "Error retrieving file", http.StatusInternalServerError)
			return
		}
		defer file.Close()

		if !strings.HasSuffix(handler.Filename, ".csv") { // hanya boleh .csv
			http.Error(w, "Only .csv files are allowed", http.StatusInternalServerError)
			return
		}

		// Membaca file content
		var buf bytes.Buffer
		if _, err := io.Copy(&buf, file); err != nil {
			http.Error(w, "Failed to read file content", http.StatusInternalServerError)
			return
		}
		fileContent := buf.String()

		// process file
		parsedData, err := fileService.ProcessFile(fileContent)
		if err != nil {
			http.Error(w, "Error processing file", http.StatusInternalServerError)
			log.Println("Error processing file:", err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(parsedData)

		// queries := []string{
		// 	"Find the least electricity usage appliance.",
		// 	"Find the most electricity usage appliance.",
		// }

		// // analisi data
		// answer, err := aiService.AnalyzeFile(parsedData, queries, token)
		// if err != nil {
		// 	http.Error(w, "Failed to analyze data: "+err.Error(), http.StatusInternalServerError)
		// 	return
		// }

		// response := model.Response{Status: "success", Answer: answer}
		// w.Header().Set("Content-Type", "application/json")
		// w.WriteHeader(http.StatusOK)
		// if err := json.NewEncoder(w).Encode(&response); err != nil {
		// 	http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		// }

	}).Methods("POST")

	// Chat endpoint
	router.HandleFunc("/chat", func(w http.ResponseWriter, r *http.Request) {

	}).Methods("POST")

	// Enable CORS
	corsHandler := cors.New(cors.Options{
		AllowedOrigins: []string{"*"}, // Allow your React app's origin
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Content-Type", "Authorization"},
	}).Handler(router)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Start the server
	log.Printf("Server running on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, corsHandler))
}
