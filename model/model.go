package model

type Response struct {
	Status string `json:"status"`
	Answer string `json:"answer"`
}

type ChatRequest struct {
	Type         string `json:"type"`
	Query        string `json:"query"`
	PreviousChat string `json:"prevChat"`
}

type Inputs struct {
	Table map[string][]string `json:"table"`
	Query string              `json:"query"`
}

type TapasRequest struct {
	Inputs Inputs `json:"inputs"`
}

type TapasResponse struct {
	Answer      string   `json:"answer"`
	Coordinates [][]int  `json:"coordinates"`
	Cells       []string `json:"cells"`
	Aggregator  string   `json:"aggregator"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type PhiRequest struct {
	Model     string    `json:"model"`
	Messages  []Message `json:"messages"`
	MaxTokens int       `json:"max_tokens"`
	Stream    bool      `json:"stream"`
}

type Choice struct {
	Message Message `json:"message"`
}

type PhiResponse struct {
	Choices []Choice `json:"choices"`
}
