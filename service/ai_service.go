package service

import (
	"fmt"
	"net/http"

	"github.com/z4fL/fp-ai-golang-neurons/model"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type AIService struct {
	Client HTTPClient
}

func (s *AIService) AnalyzeData(table map[string][]string, query, token string) (string, error) {
	return "", nil
}

func (s *AIService) AnalyzeFile(table map[string][]string, queries []string, token string) (string, error) {
	results := make([]string, 0, len(queries))

	for _, query := range queries {
		result, err := s.AnalyzeData(table, query, token)
		if err != nil {
			return "", err
		}
		results = append(results, result)
	}

	answer := fmt.Sprintf("From the provided data, here are the Least Electricity: %s and the Most Electricity: %s.", results[0], results[1])

	return answer, nil
}

func (s *AIService) ChatWithAI(context, query, token string) (model.PhiResponse, error) {
	return model.PhiResponse{}, nil
}
