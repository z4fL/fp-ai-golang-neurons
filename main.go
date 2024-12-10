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
	"github.com/z4fL/fp-ai-golang-neurons/utility"
)

// Init services
var fileService = &service.FileService{}
var aiService = &service.AIService{Client: &http.Client{}}

const defaultPort = "8080"
const dataFilePath = "upload/data-series.csv"

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
			utility.JSONResponse(w, http.StatusBadRequest, "failed", "Failed to parse form data")
			log.Printf("ParseMultipartForm error: %v", err)
			return
		}

		// Get the uploaded file
		file, handler, err := r.FormFile("file") // "file" matches the field name in the frontend
		if err != nil {
			utility.JSONResponse(w, http.StatusBadRequest, "failed", "Failed to retrieve uploaded file")
			log.Printf("FormFile error: %v", err)
			return
		}
		defer file.Close()

		if !strings.HasSuffix(handler.Filename, ".csv") { // only .csv files are allowed
			utility.JSONResponse(w, http.StatusUnsupportedMediaType, "failed", "Only .csv files are allowed")
			log.Printf("Unsupported file type: %s", handler.Filename)
			return
		}

		// Read file content
		var buf bytes.Buffer
		if _, err := io.Copy(&buf, file); err != nil {
			utility.JSONResponse(w, http.StatusInternalServerError, "failed", "Failed to read file content")
			log.Printf("io.Copy error: %v", err)
			return
		}
		fileContent := buf.String()

		// process file
		parsedData, err := fileService.ProcessFile(fileContent)
		if err != nil {
			utility.JSONResponse(w, http.StatusInternalServerError, "failed", "Failed to process file content")
			log.Printf("ProcessFile error: %v", err)
			return
		}

		queries := []string{
			"Find the least electricity usage appliance.",
			"Find the most electricity usage appliance.",
		}

		// analyze data
		answer, err := aiService.AnalyzeFile(parsedData, queries, token)
		if err != nil {
			utility.JSONResponse(w, http.StatusInternalServerError, "failed", "Failed to analyze data")
			log.Printf("AnalyzeFile error: %v", err)
			return
		}

		utility.JSONResponse(w, http.StatusOK, "success", answer)
		log.Println("Success to upload file")
	}).Methods("POST")

	// Chat endpoint
	router.HandleFunc("/chat", func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			utility.JSONResponse(w, http.StatusBadRequest, "failed", "Failed to read request body")
			log.Printf("ReadAll error: %v", err)
			return
		}
		defer r.Body.Close()

		var chatReq model.ChatRequest
		err = json.Unmarshal(body, &chatReq)
		if err != nil {
			utility.JSONResponse(w, http.StatusBadRequest, "failed", "Invalid JSON format in request body")
			log.Printf("Unmarshal error: %v", err)
			return
		}

		var answer string
		switch chatReq.Type {
		case "tapas":
			filePath := dataFilePath

			if !fileService.Repo.FileExists(filePath) {
				utility.JSONResponse(w, http.StatusNotFound, "failed", "Data file not found")
				log.Printf("File not found: %s", filePath)
				return
			}

			contentFile, err := fileService.Repo.ReadFile(filePath)
			if err != nil {
				utility.JSONResponse(w, http.StatusInternalServerError, "failed", "Failed to read data file")
				log.Printf("ReadFile error: %v", err)
				return
			}

			parsedData, err := fileService.ParseCSV(string(contentFile))
			if err != nil {
				utility.JSONResponse(w, http.StatusInternalServerError, "failed", "Failed to parse CSV data")
				log.Printf("ParseCSV error: %v", err)
				return
			}

			answer, err = aiService.AnalyzeData(parsedData, chatReq.Query, token)
			if err != nil {
				utility.JSONResponse(w, http.StatusInternalServerError, "failed", "Failed to analyze data with AI")
				log.Printf("AnalyzeData error: %v", err)
				return
			}
			log.Println("Chat request processed successfully with google/tapas-base-finetuned-wtq")

		case "phi":
			answer, err = aiService.ChatWithAI(chatReq.PreviousChat, chatReq.Query, token)
			if err != nil {
				utility.JSONResponse(w, http.StatusInternalServerError, "failed", "Failed to chat with AI Phi")
				log.Printf("ChatWithAI error: %v", err)
				return
			}
			log.Println("Chat request processed successfully with microsoft/Phi-3.5-mini-instruct")

		default:
			utility.JSONResponse(w, http.StatusBadRequest, "failed", "Invalid chat type: "+chatReq.Type)
			log.Printf("Invalid chat type: %s", chatReq.Type)
			return
		}

		utility.JSONResponse(w, http.StatusOK, "success", answer)
	}).Methods("POST")

	// Removes the session file (data-series.csv) from the server.
	router.HandleFunc("/removesession", func(w http.ResponseWriter, r *http.Request) {
		if err := fileService.Repo.RemoveFile(dataFilePath); err != nil {
			if os.IsNotExist(err) {
				utility.JSONResponse(w, http.StatusNotFound, "failed", "File not found")
				log.Printf("File not found: %s", dataFilePath)
			} else {
				utility.JSONResponse(w, http.StatusInternalServerError, "failed", "Failed to delete file")
				log.Printf("RemoveFile error: %v", err)
			}
			return
		}

		utility.JSONResponse(w, http.StatusOK, "success", "File deleted successfully")
		log.Println("Session file deleted successfully")
	}).Methods("POST")

	// Enable CORS
	corsHandler := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:5173"}, // Replace with your specific allowed origins
		// AllowedOrigins: []string{"*"}, // Allow your Vite app's origin
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Content-Type", "Authorization"},
	}).Handler(router)

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	// Start the server
	log.Printf("Server running on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, corsHandler))
}
