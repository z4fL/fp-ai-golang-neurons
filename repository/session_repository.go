package repository

import (
	"errors"

	"github.com/z4fL/fp-ai-golang-neurons/model"
	"gorm.io/gorm"
)

type SessionsRepository interface {
	AddSessions(session model.Session) error
	DeleteSession(token string) error
	UpdateSessions(session model.Session) error
	SessionAvailID(id uint) error
	SessionAvailToken(token string) (model.Session, error)
	GetUserIDByToken(token string) (uint, error)
}

type sessionsRepoImpl struct {
	db *gorm.DB
}

func NewSessionRepo(db *gorm.DB) SessionsRepository {
	return &sessionsRepoImpl{db}
}

func (s *sessionsRepoImpl) GetUserIDByToken(token string) (uint, error) {
	var session model.Session
	if err := s.db.Where("token = ?", token).First(&session).Error; err != nil {
		return 0, errors.New("invalid session token")
	}
	return session.UserID, nil
}

func (s *sessionsRepoImpl) AddSessions(session model.Session) error {
	result := s.db.Create(&session)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (s *sessionsRepoImpl) DeleteSession(token string) error {
	result := s.db.Where("token=?", token).Delete(&model.Session{})
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (s *sessionsRepoImpl) UpdateSessions(session model.Session) error {
	result := s.db.Where("user_id=?", session.UserID).Updates(model.Session{Token: session.Token, Expiry: session.Expiry})
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (s *sessionsRepoImpl) SessionAvailID(id uint) error {
	var session model.Session
	result := s.db.Where(&model.Session{UserID: id}).First(&session)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (s *sessionsRepoImpl) SessionAvailToken(token string) (model.Session, error) {
	var session model.Session
	result := s.db.Where(&model.Session{Token: token}).First(&session)

	if result.Error != nil {
		return model.Session{}, result.Error
	}

	return session, nil
}
