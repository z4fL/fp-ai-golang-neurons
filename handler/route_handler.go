package handler

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/z4fL/fp-ai-golang-neurons/model"
	"github.com/z4fL/fp-ai-golang-neurons/repository"
	"github.com/z4fL/fp-ai-golang-neurons/service"
	"github.com/z4fL/fp-ai-golang-neurons/utility"
)

// Init services
var FileService = service.NewFileService(&repository.FileRepository{})
var AIService = &service.AIService{Client: &http.Client{}}

const dataFilePath = "/upload/data-series.csv"

var Token string

func HandleUpload(w http.ResponseWriter, r *http.Request) {
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
	parsedData, err := FileService.ProcessFile(fileContent)
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
	answer, err := AIService.AnalyzeFile(parsedData, queries, Token)
	if err != nil {
		utility.JSONResponse(w, http.StatusInternalServerError, "failed", "Failed to analyze data")
		log.Printf("AnalyzeFile error: %v", err)
		return
	}

	utility.JSONResponse(w, http.StatusOK, "success", answer)
	log.Println("Success to upload file")
}

func HandleChat(w http.ResponseWriter, r *http.Request) {
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

		if !FileService.Repo.FileExists(filePath) {
			utility.JSONResponse(w, http.StatusNotFound, "failed", "Data file not found")
			log.Printf("File not found: %s", filePath)
			return
		}

		contentFile, err := FileService.Repo.ReadFile(filePath)
		if err != nil {
			utility.JSONResponse(w, http.StatusInternalServerError, "failed", "Failed to read data file")
			log.Printf("ReadFile error: %v", err)
			return
		}

		parsedData, err := FileService.ParseCSV(string(contentFile))
		if err != nil {
			utility.JSONResponse(w, http.StatusInternalServerError, "failed", "Failed to parse CSV data")
			log.Printf("ParseCSV error: %v", err)
			return
		}

		answer, err = AIService.AnalyzeData(parsedData, chatReq.Query, Token)
		if err != nil {
			utility.JSONResponse(w, http.StatusInternalServerError, "failed", "Failed to analyze data with AI")
			log.Printf("AnalyzeData error: %v", err)
			return
		}
		log.Println("Chat request processed successfully with google/tapas-base-finetuned-wtq")

	case "phi":
		answer, err = AIService.ChatWithAI(chatReq.PreviousChat, chatReq.Query, Token)
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
}

func HandleRemoveSession(w http.ResponseWriter, r *http.Request) {
	if !FileService.Repo.FileExists(dataFilePath) {
		utility.JSONResponse(w, http.StatusNotFound, "failed", "File not found")
		log.Printf("File not found: %s", dataFilePath)
		return
	}

	if err := FileService.Repo.RemoveFile(dataFilePath); err != nil {
		utility.JSONResponse(w, http.StatusInternalServerError, "failed", "Failed to delete file")
		log.Printf("Failed to delete file %s: %v", dataFilePath, err)
		return
	}

	utility.JSONResponse(w, http.StatusOK, "success", "File deleted successfully")
	log.Println("Session file deleted successfully")
}
