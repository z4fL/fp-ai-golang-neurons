package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/z4fL/fp-ai-golang-neurons/model"
	"github.com/z4fL/fp-ai-golang-neurons/utility"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type AIService struct {
	Client HTTPClient
}

func (s *AIService) AnalyzeData(table map[string][]string, query, token string) (string, error) {
	url := "https://api-inference.huggingface.co/models/google/tapas-base-finetuned-wtq"
	requestData := &model.TapasRequest{
		Inputs: model.Inputs{
			Table: table,
			Query: query,
		},
	}

	body, err := json.Marshal(*requestData)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-wait-for-model", "true")

	res, err := s.Client.Do(req)
	if err != nil {
		return "", err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return "", errors.New("failed to get a valid response from the AI model")
	}

	var tapasRes model.TapasResponse
	if err := json.NewDecoder(res.Body).Decode(&tapasRes); err != nil {
		return "", err
	}

	processor := utility.TapasProcessor{
		Cells: &tapasRes.Cells,
	}

	var answer string

	if tapasRes.Aggregator == "NONE" {
		answer = tapasRes.Answer
	} else if tapasRes.Aggregator == "COUNT" {
		count, list := processor.Count()
		answer = fmt.Sprintf("Count: %d, List: %v", count, list)
	} else if tapasRes.Aggregator == "SUM" {
		sum := processor.Sum()
		answer = fmt.Sprintf("Sum: %f", sum)
	} else if tapasRes.Aggregator == "AVERAGE" {
		avg := processor.Average()
		answer = fmt.Sprintf("Sum: %f", avg)
	} else if tapasRes.Aggregator == "MIN" {
		min, _ := processor.Min()
		answer = fmt.Sprintf("Min: %f", min)
	} else if tapasRes.Aggregator == "MIN" {
		max, _ := processor.Max()
		answer = fmt.Sprintf("Max: %f", max)
	}

	return answer, nil
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

func (s *AIService) ChatWithAI(context, query, token string) (string, error) {
	url := "https://api-inference.huggingface.co/models/microsoft/Phi-3.5-mini-instruct/v1/chat/completions"

	requestData := &model.PhiRequest{
		Model: "microsoft/Phi-3.5-mini-instruct",
		Messages: []model.Message{
			{Role: "system", Content: "You are an intelligent assistant designed to help users optimize energy consumption in their smart homes. You must respond clearly, concisely, and in a user-friendly manner. If the user asks for recommendations, base your advice on energy-saving strategies while considering the data insights."},
			{Role: "assistant", Content: context},
			{Role: "user", Content: query},
		},
		MaxTokens: 500,
		Stream:    false,
	}

	body, err := json.Marshal(*requestData)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-wait-for-model", "true")

	res, err := s.Client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return "", errors.New("failed to get a valid response from the AI model")
	}

	var phiResponse model.PhiResponse
	if err := json.NewDecoder(res.Body).Decode(&phiResponse); err != nil {
		return "", err
	}

	answer := phiResponse.Choices[0].Message.Content

	return answer, nil
}
