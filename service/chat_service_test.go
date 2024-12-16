package service_test

import (
	"encoding/json"
	"errors"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/z4fL/fp-ai-golang-neurons/model"
	"github.com/z4fL/fp-ai-golang-neurons/service"
)

type MockChatRepository struct {
	AddChatFunc       func(chat *model.Chat) (*model.Chat, error)
	GetChatUserFunc   func(userID, chatID string) (*model.Chat, error)
	UpdateChatFunc    func(chat *model.Chat) error
	ListUserChatsFunc func(userID string) ([]model.Chat, error)
}

func (m *MockChatRepository) AddChat(chat *model.Chat) (*model.Chat, error) {
	return m.AddChatFunc(chat)
}

func (m *MockChatRepository) GetChatUser(userID, chatID string) (*model.Chat, error) {
	return m.GetChatUserFunc(userID, chatID)
}

func (m *MockChatRepository) UpdateChat(chat *model.Chat) error {
	return m.UpdateChatFunc(chat)
}

func (m *MockChatRepository) ListUserChats(userID string) ([]model.Chat, error) {
	return m.ListUserChatsFunc(userID)
}

var _ = Describe("ChatService", func() {
	var (
		mockRepo    *MockChatRepository
		chatService service.ChatService
	)

	BeforeEach(func() {
		mockRepo = &MockChatRepository{}
		chatService = service.NewChatService(mockRepo)
	})

	Describe("CreateChat", func() {
		It("should create a new chat successfully", func() {
			mockRepo.AddChatFunc = func(chat *model.Chat) (*model.Chat, error) {
				return chat, nil
			}

			chatHistory := []map[string]any{{"message": "Hello"}}
			chat, err := chatService.CreateChat("user1", chatHistory)
			Expect(err).NotTo(HaveOccurred())
			Expect(chat.UserID).To(Equal("user1"))
		})

		It("should return an error if chat creation fails", func() {
			mockRepo.AddChatFunc = func(chat *model.Chat) (*model.Chat, error) {
				return nil, errors.New("failed to create chat")
			}

			chatHistory := []map[string]any{{"message": "Hello"}}
			_, err := chatService.CreateChat("user1", chatHistory)
			Expect(err).To(HaveOccurred())
		})
	})

	Describe("AddMessage", func() {
		It("should add a new message to the chat", func() {
			mockRepo.GetChatUserFunc = func(userID, chatID string) (*model.Chat, error) {
				chatHistory, _ := json.Marshal([]map[string]any{{"message": "Hello"}})
				return &model.Chat{UserID: userID, ChatHistory: chatHistory}, nil
			}
			mockRepo.UpdateChatFunc = func(chat *model.Chat) error {
				return nil
			}

			newMessage := []map[string]any{{"message": "Hi"}}
			err := chatService.AddMessage("user1", "chat1", newMessage)
			Expect(err).NotTo(HaveOccurred())
		})

		It("should return an error if chat is not found", func() {
			mockRepo.GetChatUserFunc = func(userID, chatID string) (*model.Chat, error) {
				return nil, errors.New("chat not found")
			}

			newMessage := []map[string]any{{"message": "Hi"}}
			err := chatService.AddMessage("user1", "chat1", newMessage)
			Expect(err).To(HaveOccurred())
		})
	})

	Describe("GetChatUser", func() {
		It("should return chat history for a user", func() {
			mockRepo.GetChatUserFunc = func(userID, chatID string) (*model.Chat, error) {
				chatHistory, _ := json.Marshal([]map[string]any{{"message": "Hello"}})
				return &model.Chat{UserID: userID, ChatHistory: chatHistory}, nil
			}

			chatHistory, err := chatService.GetChatUser("user1", "chat1")
			Expect(err).NotTo(HaveOccurred())
			Expect(chatHistory).To(HaveLen(1))
		})

		It("should return an error if chat is not found", func() {
			mockRepo.GetChatUserFunc = func(userID, chatID string) (*model.Chat, error) {
				return nil, errors.New("chat not found")
			}

			_, err := chatService.GetChatUser("user1", "chat1")
			Expect(err).To(HaveOccurred())
		})
	})

	Describe("ListUserChats", func() {
		It("should return a list of user chats", func() {
			mockRepo.ListUserChatsFunc = func(userID string) ([]model.Chat, error) {
				chatHistory, _ := json.Marshal([]model.ChatHistoryEntry{{ID: 3, Content: "Hello"}})
				return []model.Chat{{UserID: userID, ChatHistory: chatHistory}}, nil
			}

			chats, err := chatService.ListUserChats("user1")
			Expect(err).NotTo(HaveOccurred())
			Expect(chats).To(HaveLen(1))
		})

		It("should return an error if no chats are found", func() {
			mockRepo.ListUserChatsFunc = func(userID string) ([]model.Chat, error) {
				return nil, errors.New("no chats found")
			}

			_, err := chatService.ListUserChats("user1")
			Expect(err).To(HaveOccurred())
		})
	})
})
