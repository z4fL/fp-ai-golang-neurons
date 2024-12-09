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
	"github.com/z4fL/fp-ai-golang-neurons/model"
	"github.com/z4fL/fp-ai-golang-neurons/service"
)

// Init services
var fileService = &service.FileService{}
var aiService = &service.AIService{Client: &http.Client{}}

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
			log.Println("Unable to parse form: ", err.Error())
			return
		}

		// Get the uploaded file
		file, handler, err := r.FormFile("file") // "file" sesuai nama field di frontend
		if err != nil {
			http.Error(w, "Error retrieving file", http.StatusInternalServerError)
			log.Println("Error retrieving file: ", err.Error())
			return
		}
		defer file.Close()

		if !strings.HasSuffix(handler.Filename, ".csv") { // hanya boleh .csv
			http.Error(w, "Only .csv files are allowed", http.StatusInternalServerError)
			log.Println("Only .csv file")
			return
		}

		// Membaca file content
		var buf bytes.Buffer
		if _, err := io.Copy(&buf, file); err != nil {
			http.Error(w, "Failed to read file content", http.StatusInternalServerError)
			log.Println("Failed to read file content: ", err.Error())
			return
		}
		fileContent := buf.String()

		// process file
		parsedData, err := fileService.ProcessFile(fileContent)
		if err != nil {
			http.Error(w, "Error processing file", http.StatusInternalServerError)
			log.Println("Error processing file: ", err)
			return
		}

		queries := []string{
			"Find the least electricity usage appliance.",
			"Find the most electricity usage appliance.",
		}

		// analisi data
		answer, err := aiService.AnalyzeFile(parsedData, queries, token)
		if err != nil {
			http.Error(w, "Failed to analyze data", http.StatusInternalServerError)
			log.Println("Failed to analyze data: ", err.Error())
			return
		}

		response := model.Response{Status: "success", Answer: answer}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(&response); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			log.Println("Failed to encode response: ", err.Error())
		}

		log.Println("Success to upload file")
	}).Methods("POST")

	// Chat endpoint
	router.HandleFunc("/chat", func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error read Body", http.StatusInternalServerError)
			log.Println("Error read Body: ", err.Error())
		}
		defer r.Body.Close()

		var chatReq model.ChatRequest
		err = json.Unmarshal(body, &chatReq)
		if err != nil {
			http.Error(w, "Error unmarshalling request body", http.StatusBadRequest)
			log.Println("Error unmarshalling request body", err.Error())
			return
		}

		var answer string

		if chatReq.Type == "tapas" {
			filePath := "upload/data-series.csv"
			var contentFile string

			if fileService.Repo.FileExists(filePath) {
				existingContent, err := fileService.Repo.ReadFile(filePath)
				if err != nil {
					http.Error(w, "Error ", http.StatusBadRequest)
					log.Println("Error ", err.Error())
					return
				}
				contentFile = string(existingContent)
			} else {
				http.Error(w, "Error, file not found", http.StatusBadRequest)
				log.Println("Error not found ", filePath)
				return
			}

			parsedData, err := fileService.ParseCSV(contentFile)
			if err != nil {
				http.Error(w, "Error ", http.StatusBadRequest)
				log.Println("Error ", err.Error())
				return
			}

			answer, err = aiService.AnalyzeData(parsedData, chatReq.Query, token)
			if err != nil {
				http.Error(w, "Failed to chat with ai", http.StatusInternalServerError)
				log.Println("Failed to chat with ai: ", err.Error())
				return
			}
		} else if chatReq.Type == "phi" {
			answer, err = aiService.ChatWithAI(chatReq.PreviousChat, chatReq.Query, token)
			if err != nil {
				http.Error(w, "Failed to chat with ai", http.StatusInternalServerError)
				log.Println("Failed to chat with ai: ", err.Error())
				return
			}
		}

		response := model.Response{Status: "success", Answer: answer}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(&response); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			log.Println("Failed to encode response: ", err.Error())
		}

		log.Println("Success to chat with AI")
	}).Methods("POST")

	router.HandleFunc("/removesession", func(w http.ResponseWriter, r *http.Request) {
		filePath := "upload/data-series.csv"
		err := fileService.Repo.RemoveFile(filePath)
		if err != nil {
			http.Error(w, "Error deleting file", http.StatusInternalServerError)
			log.Printf("Error deleting file %s: %v", filePath, err)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("File deleted successfully"))
	}).Methods("POST")

	// Enable CORS
	corsHandler := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		// AllowedOrigins: []string{"http://localhost:5173"}, // Allow your Vite app's origin
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
