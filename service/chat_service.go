package service

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/z4fL/fp-ai-golang-neurons/model"
	"github.com/z4fL/fp-ai-golang-neurons/repository"
)

type ChatService interface {
	CreateChat(userID string, chatHistory []map[string]any) (*model.Chat, error)
	AddMessage(userID, chatID string, newMessage []map[string]any) error
	GetChatUser(userID, chatID string) ([]map[string]any, error)
	ListUserChats(userID string) ([]map[string]any, error)
}

type chatService struct {
	repo repository.ChatRepository
}

func NewChatService(repo repository.ChatRepository) ChatService {
	return &chatService{repo: repo}
}

func (s *chatService) ListUserChats(userID string) ([]map[string]any, error) {
	chats, err := s.repo.ListUserChats(userID)
	if err != nil {
		return nil, err
	}

	var result []map[string]any
	for _, chat := range chats {
		var chatHistory []model.ChatHistoryEntry
		if err := json.Unmarshal([]byte(chat.ChatHistory), &chatHistory); err != nil {
			return nil, err
		}

		for _, entry := range chatHistory {
			if entry.ID == 3 {
				if contentStr, ok := entry.Content.(string); ok {
					words := strings.Fields(contentStr)
					if len(words) > 4 {
						contentStr = strings.Join(words[:4], " ") + "..."
					}
					contentMap := map[string]any{
						"chatID":  chat.ID,
						"content": contentStr,
					}
					result = append(result, contentMap)
				}
				break
			}
		}
	}

	if len(result) == 0 {
		return nil, errors.New("no content with id 3 found")
	}

	return result, nil
}

func (s *chatService) GetChatUser(userID, chatID string) ([]map[string]any, error) {
	chat, err := s.repo.GetChatUser(userID, chatID)
	if err != nil {
		return nil, errors.New("chat not found")
	}

	var chatHistory []map[string]any
	if err := json.Unmarshal(chat.ChatHistory, &chatHistory); err != nil {
		return nil, err
	}

	return chatHistory, nil
}

func (s *chatService) CreateChat(userID string, chatHistory []map[string]any) (*model.Chat, error) {
	// Serialize chatHistory to JSON
	chatHistoryJSON, err := json.Marshal(chatHistory)
	if err != nil {
		return nil, err
	}

	chat := &model.Chat{
		UserID:      userID,
		ChatHistory: chatHistoryJSON,
	}

	createdChat, err := s.repo.AddChat(chat)
	if err != nil {
		return nil, err
	}

	return createdChat, nil
}

func (s *chatService) AddMessage(userID, chatID string, newMessage []map[string]any) error {
	// Get existing chat
	chat, err := s.repo.GetChatUser(userID, chatID)
	if err != nil {
		return errors.New("chat not found")
	}

	// Deserialize chat history
	var chatHistory []map[string]any
	if err := json.Unmarshal(chat.ChatHistory, &chatHistory); err != nil {
		return err
	}

	// Append new message
	chatHistory = append(chatHistory, newMessage...)

	// Serialize updated chat history
	updatedChatHistory, err := json.Marshal(chatHistory)
	if err != nil {
		return err
	}

	// Update chat
	chat.ChatHistory = updatedChatHistory
	return s.repo.UpdateChat(chat)
}
