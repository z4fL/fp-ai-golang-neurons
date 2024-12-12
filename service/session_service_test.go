package service_test

import (
	"errors"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/z4fL/fp-ai-golang-neurons/model"
	"github.com/z4fL/fp-ai-golang-neurons/service"
)

type MockSessionsRepository struct {
	SessionAvailNameFunc  func(username string) error
	AddSessionsFunc       func(session model.Session) error
	UpdateSessionsFunc    func(session model.Session) error
	DeleteSessionFunc     func(sessionToken string) error
	SessionAvailTokenFunc func(token string) (model.Session, error)
}

func (m *MockSessionsRepository) SessionAvailName(username string) error {
	return m.SessionAvailNameFunc(username)
}

func (m *MockSessionsRepository) AddSessions(session model.Session) error {
	return m.AddSessionsFunc(session)
}

func (m *MockSessionsRepository) UpdateSessions(session model.Session) error {
	return m.UpdateSessionsFunc(session)
}

func (m *MockSessionsRepository) DeleteSession(sessionToken string) error {
	return m.DeleteSessionFunc(sessionToken)
}

func (m *MockSessionsRepository) SessionAvailToken(token string) (model.Session, error) {
	return m.SessionAvailTokenFunc(token)
}

var _ = Describe("SessionService", func() {
	var (
		mockRepo       *MockSessionsRepository
		sessionService service.SessionService
	)

	BeforeEach(func() {
		mockRepo = &MockSessionsRepository{}
		sessionService = service.NewSessionService(mockRepo)
	})

	Describe("SessionAvailName", func() {
		It("should return an error if session name is not available", func() {
			mockRepo.SessionAvailNameFunc = func(username string) error {
				return errors.New("username not available")
			}

			err := sessionService.SessionAvailName("testuser")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("username not available"))
		})

		It("should return nil if session name is available", func() {
			mockRepo.SessionAvailNameFunc = func(username string) error {
				return nil
			}

			err := sessionService.SessionAvailName("testuser")
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Describe("AddSession", func() {
		It("should return an error if adding session fails", func() {
			mockRepo.AddSessionsFunc = func(session model.Session) error {
				return errors.New("failed to add session")
			}

			err := sessionService.AddSession(model.Session{Username: "testuser"})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("failed to add session"))
		})

		It("should add session successfully", func() {
			mockRepo.AddSessionsFunc = func(session model.Session) error {
				return nil
			}

			err := sessionService.AddSession(model.Session{Username: "testuser"})
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Describe("UpdateSession", func() {
		It("should return an error if updating session fails", func() {
			mockRepo.UpdateSessionsFunc = func(session model.Session) error {
				return errors.New("failed to update session")
			}

			err := sessionService.UpdateSession(model.Session{Username: "testuser"})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("failed to update session"))
		})

		It("should update session successfully", func() {
			mockRepo.UpdateSessionsFunc = func(session model.Session) error {
				return nil
			}

			err := sessionService.UpdateSession(model.Session{Username: "testuser"})
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Describe("DeleteSession", func() {
		It("should return an error if deleting session fails", func() {
			mockRepo.DeleteSessionFunc = func(sessionToken string) error {
				return errors.New("failed to delete session")
			}

			err := sessionService.DeleteSession("testtoken")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("failed to delete session"))
		})

		It("should delete session successfully", func() {
			mockRepo.DeleteSessionFunc = func(sessionToken string) error {
				return nil
			}

			err := sessionService.DeleteSession("testtoken")
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Describe("TokenValidity", func() {
		It("should return an error if session token is not valid", func() {
			mockRepo.SessionAvailTokenFunc = func(token string) (model.Session, error) {
				return model.Session{}, errors.New("invalid token")
			}

			_, err := sessionService.TokenValidity("invalidtoken")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("invalid token"))
		})

		It("should return an error if token is expired", func() {
			mockRepo.SessionAvailTokenFunc = func(token string) (model.Session, error) {
				return model.Session{Expiry: time.Now().Add(-time.Hour)}, nil
			}
			mockRepo.DeleteSessionFunc = func(sessionToken string) error {
				return nil
			}

			_, err := sessionService.TokenValidity("expiredtoken")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("Token is Expired!"))
		})

		It("should return session if token is valid", func() {
			mockRepo.SessionAvailTokenFunc = func(token string) (model.Session, error) {
				return model.Session{Expiry: time.Now().Add(time.Hour)}, nil
			}

			session, err := sessionService.TokenValidity("validtoken")
			Expect(err).NotTo(HaveOccurred())
			Expect(session.Expiry.After(time.Now())).To(BeTrue())
		})
	})

	Describe("TokenExpired", func() {
		It("should return true if token is expired", func() {
			expiredSession := model.Session{Expiry: time.Now().Add(-time.Hour)}
			Expect(sessionService.TokenExpired(expiredSession)).To(BeTrue())
		})

		It("should return false if token is not expired", func() {
			validSession := model.Session{Expiry: time.Now().Add(time.Hour)}
			Expect(sessionService.TokenExpired(validSession)).To(BeFalse())
		})
	})
})
