package utility

import (
	"encoding/json"
	"net/http"

	"github.com/z4fL/fp-ai-golang-neurons/model"
)

func JSONResponse(w http.ResponseWriter, statusCode int, status string, answer string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := model.Response{
		Status: status,
		Answer: answer,
	}

	json.NewEncoder(w).Encode(&response)
}
