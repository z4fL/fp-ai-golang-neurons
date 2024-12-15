package model

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `gorm:"type:varchar(100);unique" json:"username"`
	Password string `gorm:"type:varchar(100)" json:"password"`
}

type Session struct {
	gorm.Model
	Token  string    `json:"token"`
	UserID uint      `json:"user_id"`
	Expiry time.Time `json:"expiry"`
}

type Chat struct {
	gorm.Model
	UserID      string         `gorm:"index;not null"`
	ChatHistory datatypes.JSON `gorm:"type:jsonb"` // Simpan history sebagai JSONB
}

type ChatHistoryEntry struct {
	ID      int    `json:"id"`
	Role    string `json:"role"`
	Type    string `json:"type"`
	Content any    `json:"content"`
}

type DBCredential struct {
	Host         string
	Username     string
	Password     string
	DatabaseName string
	Port         int
	Schema       string
}

type Response struct {
	Status string `json:"status"`
	Answer any    `json:"answer"`
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
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	MaxTokens   int       `json:"max_tokens"`
	Stream      bool      `json:"stream"`
	Temperature float64   `json:"temperature"`
}

type Choice struct {
	Message Message `json:"message"`
}

type PhiResponse struct {
	Choices []Choice `json:"choices"`
}
