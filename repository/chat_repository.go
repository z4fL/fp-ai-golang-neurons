package repository

import (
	"github.com/z4fL/fp-ai-golang-neurons/model"
	"gorm.io/gorm"
)

type ChatRepository interface {
	AddChat(chat *model.Chat) error
	GetChatUser(userID, chatID string) (*model.Chat, error)
	UpdateChat(chat *model.Chat) error
	ListUserChats(userID string) ([]model.Chat, error)
}

type chatRepository struct {
	db *gorm.DB
}

func NewChatRepository(db *gorm.DB) ChatRepository {
	return &chatRepository{db: db}
}

func (r *chatRepository) AddChat(chat *model.Chat) error {
	return r.db.Create(chat).Error
}

func (r *chatRepository) ListUserChats(userID string) ([]model.Chat, error) {
	var chats []model.Chat
	err := r.db.Where("user_id = ?", userID).Find(&chats).Error
	if err != nil {
		return nil, err
	}
	return chats, nil
}

func (r *chatRepository) GetChatUser(userID, chatID string) (*model.Chat, error) {
	var chat model.Chat
	if err := r.db.Where("user_id = ? AND id = ?", userID, chatID).First(&chat).Error; err != nil {
		return nil, err
	}
	return &chat, nil
}

func (r *chatRepository) UpdateChat(chat *model.Chat) error {
	return r.db.Save(chat).Error
}
