package api

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/z4fL/fp-ai-golang-neurons/utility"
)

func (api *API) Upload(w http.ResponseWriter, r *http.Request) {
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
	parsedData, err := api.fileService.ProcessFile(fileContent)
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
	answer, err := api.aiService.AnalyzeFile(parsedData, queries, api.token)
	if err != nil {
		utility.JSONResponse(w, http.StatusInternalServerError, "failed", "Failed to analyze data")
		log.Printf("AnalyzeFile error: %v", err)
		return
	}

	utility.JSONResponse(w, http.StatusOK, "success", answer)
	log.Println("Success to upload file")
}
