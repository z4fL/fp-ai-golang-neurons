package service

import (
	"encoding/json"
	"errors"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"github.com/z4fL/fp-ai-golang-neurons/model"
)

type MockChatRepository struct {
	chats map[string]*model.Chat
}

func NewMockChatRepository() *MockChatRepository {
	return &MockChatRepository{chats: make(map[string]*model.Chat)}
}

func (m *MockChatRepository) AddChat(chat *model.Chat) error {
	if _, exists := m.chats[chat.UserID]; exists {
		return errors.New("chat already exists")
	}
	m.chats[chat.UserID] = chat
	return nil
}

func (m *MockChatRepository) GetChatByUserID(userID string) (*model.Chat, error) {
	chat, exists := m.chats[userID]
	if !exists {
		return nil, errors.New("chat not found")
	}
	return chat, nil
}

func (m *MockChatRepository) UpdateChat(chat *model.Chat) error {
	m.chats[chat.UserID] = chat
	return nil
}

var _ = ginkgo.Describe("ChatService", func() {
	var (
		mockRepo    *MockChatRepository
		chatService ChatService
		userID      string
		chatHistory []map[string]any
		newMessage  map[string]any
	)

	ginkgo.BeforeEach(func() {
		mockRepo = NewMockChatRepository()
		chatService = NewChatService(mockRepo)
		userID = "user123"
		chatHistory = []map[string]any{
			{"message": "Hello"},
		}
		newMessage = map[string]any{"message": "Hi"}
	})

	ginkgo.Describe("CreateChat", func() {
		ginkgo.It("should create a new chat", func() {
			err := chatService.CreateChat(userID, chatHistory)
			gomega.Expect(err).To(gomega.BeNil())

			chat, err := mockRepo.GetChatByUserID(userID)
			gomega.Expect(err).To(gomega.BeNil())
			gomega.Expect(chat).ToNot(gomega.BeNil())
			gomega.Expect(chat.UserID).To(gomega.Equal(userID))

			var storedChatHistory []map[string]any
			err = json.Unmarshal(chat.ChatHistory, &storedChatHistory)
			gomega.Expect(err).To(gomega.BeNil())
			gomega.Expect(storedChatHistory).To(gomega.Equal(chatHistory))
		})

		ginkgo.It("should return an error if chat history serialization fails", func() {
			invalidChatHistory := []map[string]any{
				{"message": make(chan int)},
			}

			err := chatService.CreateChat(userID, invalidChatHistory)
			gomega.Expect(err).ToNot(gomega.BeNil())
		})
	})

	ginkgo.Describe("AddMessage", func() {
		ginkgo.It("should add a new message to the chat", func() {
			err := chatService.CreateChat(userID, chatHistory)
			gomega.Expect(err).To(gomega.BeNil())

			err = chatService.AddMessage(userID, newMessage)
			gomega.Expect(err).To(gomega.BeNil())

			chat, err := mockRepo.GetChatByUserID(userID)
			gomega.Expect(err).To(gomega.BeNil())
			gomega.Expect(chat).ToNot(gomega.BeNil())

			var storedChatHistory []map[string]any
			err = json.Unmarshal(chat.ChatHistory, &storedChatHistory)
			gomega.Expect(err).To(gomega.BeNil())
			gomega.Expect(storedChatHistory).To(gomega.Equal(append(chatHistory, newMessage)))
		})

		ginkgo.It("should return an error if chat is not found", func() {
			err := chatService.AddMessage(userID, newMessage)
			gomega.Expect(err).To(gomega.MatchError("chat not found"))
		})

		ginkgo.It("should return an error if chat history deserialization fails", func() {
			invalidChat := &model.Chat{
				UserID:      userID,
				ChatHistory: []byte(`invalid json`),
			}
			mockRepo.chats[userID] = invalidChat

			err := chatService.AddMessage(userID, newMessage)
			gomega.Expect(err).ToNot(gomega.BeNil())
		})

		ginkgo.It("should return an error if chat history serialization fails", func() {
			err := chatService.CreateChat(userID, chatHistory)
			gomega.Expect(err).To(gomega.BeNil())

			invalidMessage := map[string]any{"message": make(chan int)}
			err = chatService.AddMessage(userID, invalidMessage)
			gomega.Expect(err).ToNot(gomega.BeNil())
		})
	})
})
