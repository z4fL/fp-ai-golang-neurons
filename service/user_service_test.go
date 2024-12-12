package service_test

import (
	"errors"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/z4fL/fp-ai-golang-neurons/model"
	"github.com/z4fL/fp-ai-golang-neurons/service"
)

type MockUserRepository struct {
	AddFunc          func(user model.User) error
	AuthenticateFunc func(username, password string) (model.User, error)
}

func (m *MockUserRepository) Add(user model.User) error {
	return m.AddFunc(user)
}

func (m *MockUserRepository) Authenticate(username, password string) (model.User, error) {
	return m.AuthenticateFunc(username, password)
}

var _ = Describe("UserService", func() {
	var (
		mockRepo    *MockUserRepository
		userService service.UserService
	)

	BeforeEach(func() {
		mockRepo = &MockUserRepository{}
		userService = service.NewUserService(mockRepo)
	})

	Describe("Register", func() {
		It("should return an error if user registration fails", func() {
			mockRepo.AddFunc = func(user model.User) error {
				return errors.New("registration failed")
			}

			err := userService.Register(model.User{Username: "testuser", Password: "password"})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("registration failed"))
		})

		It("should register a user successfully", func() {
			mockRepo.AddFunc = func(user model.User) error {
				return nil
			}

			err := userService.Register(model.User{Username: "testuser", Password: "password"})
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Describe("Login", func() {
		It("should return an error if authentication fails", func() {
			mockRepo.AuthenticateFunc = func(username, password string) (model.User, error) {
				return model.User{}, errors.New("authentication failed")
			}

			_, err := userService.Login("testuser", "wrongpassword")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("authentication failed"))
		})

		It("should login a user successfully", func() {
			mockRepo.AuthenticateFunc = func(username, password string) (model.User, error) {
				return model.User{Username: "testuser", Password: "password"}, nil
			}

			user, err := userService.Login("testuser", "password")
			Expect(err).NotTo(HaveOccurred())
			Expect(user.Username).To(Equal("testuser"))
		})
	})
})
