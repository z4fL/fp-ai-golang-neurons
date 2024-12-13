package service

import (
	"fmt"
	"time"

	"github.com/z4fL/fp-ai-golang-neurons/model"
	"github.com/z4fL/fp-ai-golang-neurons/repository"
)

type SessionService interface {
	AddSession(session model.Session) error
	UpdateSession(session model.Session) error
	DeleteSession(sessionToken string) error
	SessionAvailID(id uint) error
	TokenExpired(session model.Session) bool
	TokenValidity(token string) (model.Session, error)
	GetUserIDByToken(token string) (uint, error)
}

type sessionService struct {
	sessionRepository repository.SessionsRepository
}

func NewSessionService(sessionRepository repository.SessionsRepository) SessionService {
	return &sessionService{sessionRepository}
}

func (s *sessionService) GetUserIDByToken(token string) (uint, error) {
	return s.sessionRepository.GetUserIDByToken(token)
}

func (s *sessionService) SessionAvailID(id uint) error {
	return s.sessionRepository.SessionAvailID(id)
}

func (s *sessionService) AddSession(session model.Session) error {
	return s.sessionRepository.AddSessions(session)
}

func (s *sessionService) UpdateSession(session model.Session) error {
	return s.sessionRepository.UpdateSessions(session)
}

func (s *sessionService) DeleteSession(sessionToken string) error {
	return s.sessionRepository.DeleteSession(sessionToken)
}

func (s *sessionService) TokenValidity(token string) (model.Session, error) {
	session, err := s.sessionRepository.SessionAvailToken(token)
	if err != nil {
		return model.Session{}, err
	}

	if s.TokenExpired(session) {
		err := s.sessionRepository.DeleteSession(token)
		if err != nil {
			return model.Session{}, err
		}
		return model.Session{}, fmt.Errorf("Token is Expired!")
	}

	return session, nil
}

func (s *sessionService) TokenExpired(session model.Session) bool {
	return session.Expiry.Before(time.Now())
}
