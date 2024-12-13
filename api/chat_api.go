package api

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"path/filepath"

	"github.com/gorilla/mux"
	"github.com/z4fL/fp-ai-golang-neurons/model"
	"github.com/z4fL/fp-ai-golang-neurons/utility"
	"github.com/z4fL/fp-ai-golang-neurons/utility/projectpath"
)

const dataFilePath = "upload/data-series.csv"

var dir = filepath.Join(projectpath.Root, dataFilePath)

func (h *API) ChatWithAI(w http.ResponseWriter, r *http.Request) {
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
		filePath := dir

		if !h.fileService.GetRepo().FileExists(filePath) {
			utility.JSONResponse(w, http.StatusNotFound, "failed", "Data file not found")
			log.Printf("File not found: %s", filePath)
			return
		}

		contentFile, err := h.fileService.GetRepo().ReadFile(filePath)
		if err != nil {
			utility.JSONResponse(w, http.StatusInternalServerError, "failed", "Failed to read data file")
			log.Printf("ReadFile error: %v", err)
			return
		}

		parsedData, err := h.fileService.ParseCSV(string(contentFile))
		if err != nil {
			utility.JSONResponse(w, http.StatusInternalServerError, "failed", "Failed to parse CSV data")
			log.Printf("ParseCSV error: %v", err)
			return
		}

		answer, err = h.aiService.AnalyzeData(parsedData, chatReq.Query, h.token)
		if err != nil {
			utility.JSONResponse(w, http.StatusInternalServerError, "failed", "Failed to analyze data with AI")
			log.Printf("AnalyzeData error: %v", err)
			return
		}
		log.Println("Chat request processed successfully with google/tapas-base-finetuned-wtq")

	case "phi":
		answer, err = h.aiService.ChatWithAI(chatReq.PreviousChat, chatReq.Query, h.token)
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

func (h *API) CreateChat(w http.ResponseWriter, r *http.Request) {
	var req struct {
		UserID      string           `json:"user_id"`
		ChatHistory []map[string]any `json:"chat_history"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if err := h.chatService.CreateChat(req.UserID, req.ChatHistory); err != nil {
		http.Error(w, "Failed to create chat", http.StatusInternalServerError)
		return
	}

	utility.JSONResponse(w, http.StatusCreated, "success", "Chat created successfully")
}

func (h *API) AddMessage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["userID"]

	var req map[string]any
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if err := h.chatService.AddMessage(userID, req); err != nil {
		http.Error(w, "Failed to add message", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Message added successfully"})
}

func (h *API) RemoveSession(w http.ResponseWriter, r *http.Request) {
	if !h.fileService.GetRepo().FileExists(dir) {
		utility.JSONResponse(w, http.StatusNotFound, "failed", "File not found")
		log.Printf("File not found: %s", dir)
		return
	}

	if err := h.fileService.GetRepo().RemoveFile(dir); err != nil {
		utility.JSONResponse(w, http.StatusInternalServerError, "failed", "Failed to delete file")
		log.Printf("Failed to delete file %s: %v", dir, err)
		return
	}

	utility.JSONResponse(w, http.StatusOK, "success", "File deleted successfully")
	log.Println("Session file deleted successfully")
}
