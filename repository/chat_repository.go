package repository

import (
	"github.com/z4fL/fp-ai-golang-neurons/model"
	"gorm.io/gorm"
)

type ChatRepository interface {
	AddChat(chat *model.Chat) error
	GetChatByUserID(userID string) (*model.Chat, error)
	UpdateChat(chat *model.Chat) error
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

func (r *chatRepository) GetChatByUserID(userID string) (*model.Chat, error) {
	var chat model.Chat
	if err := r.db.Where("user_id = ?", userID).First(&chat).Error; err != nil {
		return nil, err
	}
	return &chat, nil
}

func (r *chatRepository) UpdateChat(chat *model.Chat) error {
	return r.db.Save(chat).Error
}
