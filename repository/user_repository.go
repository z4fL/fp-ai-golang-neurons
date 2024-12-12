package repository

import (
	"errors"

	"github.com/z4fL/fp-ai-golang-neurons/model"
	"gorm.io/gorm"
)

type UserRepository interface {
	Add(user model.User) error
	Authenticate(username, password string) (model.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db}
}

func (r *userRepository) Add(user model.User) error {
	return r.db.Create(&user).Error
}

func (r *userRepository) Authenticate(username, password string) (model.User, error) {
	var user model.User
	if err := r.db.Where("username = ?", username).First(&user).Error; err != nil {
		return model.User{}, err
	}

	if user.Password != password {
		return model.User{}, errors.New("invalid username or password")
	}

	return user, nil
}
