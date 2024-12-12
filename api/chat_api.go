package api

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"path/filepath"

	"github.com/z4fL/fp-ai-golang-neurons/model"
	"github.com/z4fL/fp-ai-golang-neurons/utility"
	"github.com/z4fL/fp-ai-golang-neurons/utility/projectpath"
)

const dataFilePath = "upload/data-series.csv"

var dir = filepath.Join(projectpath.Root, dataFilePath)

func (api *API) Chat(w http.ResponseWriter, r *http.Request) {
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

		if !api.fileService.GetRepo().FileExists(filePath) {
			utility.JSONResponse(w, http.StatusNotFound, "failed", "Data file not found")
			log.Printf("File not found: %s", filePath)
			return
		}

		contentFile, err := api.fileService.GetRepo().ReadFile(filePath)
		if err != nil {
			utility.JSONResponse(w, http.StatusInternalServerError, "failed", "Failed to read data file")
			log.Printf("ReadFile error: %v", err)
			return
		}

		parsedData, err := api.fileService.ParseCSV(string(contentFile))
		if err != nil {
			utility.JSONResponse(w, http.StatusInternalServerError, "failed", "Failed to parse CSV data")
			log.Printf("ParseCSV error: %v", err)
			return
		}

		answer, err = api.aiService.AnalyzeData(parsedData, chatReq.Query, api.token)
		if err != nil {
			utility.JSONResponse(w, http.StatusInternalServerError, "failed", "Failed to analyze data with AI")
			log.Printf("AnalyzeData error: %v", err)
			return
		}
		log.Println("Chat request processed successfully with google/tapas-base-finetuned-wtq")

	case "phi":
		answer, err = api.aiService.ChatWithAI(chatReq.PreviousChat, chatReq.Query, api.token)
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

func (api *API) RemoveSession(w http.ResponseWriter, r *http.Request) {
	if !api.fileService.GetRepo().FileExists(dir) {
		utility.JSONResponse(w, http.StatusNotFound, "failed", "File not found")
		log.Printf("File not found: %s", dir)
		return
	}

	if err := api.fileService.GetRepo().RemoveFile(dir); err != nil {
		utility.JSONResponse(w, http.StatusInternalServerError, "failed", "Failed to delete file")
		log.Printf("Failed to delete file %s: %v", dir, err)
		return
	}

	utility.JSONResponse(w, http.StatusOK, "success", "File deleted successfully")
	log.Println("Session file deleted successfully")
}
