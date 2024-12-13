package service

import (
	"encoding/json"
	"errors"

	"github.com/z4fL/fp-ai-golang-neurons/model"
	"github.com/z4fL/fp-ai-golang-neurons/repository"
)

type ChatService interface {
	CreateChat(userID string, chatHistory []map[string]any) error
	AddMessage(userID string, newMessage map[string]any) error
}

type chatService struct {
	repo repository.ChatRepository
}

func NewChatService(repo repository.ChatRepository) ChatService {
	return &chatService{repo: repo}
}

func (s *chatService) CreateChat(userID string, chatHistory []map[string]any) error {
	// Serialize chatHistory to JSON
	chatHistoryJSON, err := json.Marshal(chatHistory)
	if err != nil {
		return err
	}

	chat := &model.Chat{
		UserID:      userID,
		ChatHistory: chatHistoryJSON,
	}

	return s.repo.AddChat(chat)
}

func (s *chatService) AddMessage(userID string, newMessage map[string]any) error {
	// Get existing chat
	chat, err := s.repo.GetChatByUserID(userID)
	if err != nil {
		return errors.New("chat not found")
	}

	// Deserialize chat history
	var chatHistory []map[string]any
	if err := json.Unmarshal(chat.ChatHistory, &chatHistory); err != nil {
		return err
	}

	// Append new message
	chatHistory = append(chatHistory, newMessage)

	// Serialize updated chat history
	updatedChatHistory, err := json.Marshal(chatHistory)
	if err != nil {
		return err
	}

	// Update chat
	chat.ChatHistory = updatedChatHistory
	return s.repo.UpdateChat(chat)
}
