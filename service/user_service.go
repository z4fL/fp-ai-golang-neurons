package service

import (
	"errors"

	"github.com/z4fL/fp-ai-golang-neurons/model"
	"github.com/z4fL/fp-ai-golang-neurons/repository"
)

type UserService interface {
	Register(user model.User) error
	Login(username, password string) (model.User, error)
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{userRepo}
}

func (s *userService) Register(user model.User) error {
	return s.userRepo.Add(user)
}

func (s *userService) Login(username, password string) (model.User, error) {
	user, err := s.userRepo.Authenticate(username, password)
	if err != nil {
		return model.User{}, errors.New("authentication failed")
	}
	return user, nil
}
