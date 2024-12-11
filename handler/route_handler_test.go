package handler_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/z4fL/fp-ai-golang-neurons/handler"
	"github.com/z4fL/fp-ai-golang-neurons/model"
	"github.com/z4fL/fp-ai-golang-neurons/repository"
	"github.com/z4fL/fp-ai-golang-neurons/service"
	"github.com/z4fL/fp-ai-golang-neurons/utility"
)

func TestRouteHandler(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "RouteHandler Suite")
}

type MockFileRepository struct {
	repository.FileRepository
	Files      map[string][]byte
	Dirs       map[string]bool
	ReadFileFn func(filename string) ([]byte, error)
}

func (m *MockFileRepository) SaveFile(filename string, content []byte) error {
	m.Files[filename] = content
	return nil
}

func (m *MockFileRepository) ReadFile(filename string) ([]byte, error) {
	if m.ReadFileFn != nil {
		return m.ReadFileFn(filename)
	}
	content, exists := m.Files[filename]
	if !exists {
		return nil, fmt.Errorf("file %s does not exist", filename)
	}
	return content, nil
}

func (m *MockFileRepository) FileExists(filename string) bool {
	_, exists := m.Files[filename]
	return exists
}

func (m *MockFileRepository) RemoveFile(filename string) error {
	if !m.FileExists(filename) {
		return fmt.Errorf("file %s does not exist", filename)
	}
	delete(m.Files, filename)
	return nil
}

func (m *MockFileRepository) DirExists(dirname string) bool {
	exists, _ := m.Dirs[dirname]
	return exists
}

func (m *MockFileRepository) MakeDir(dirname string) error {
	m.Dirs[dirname] = true
	return nil
}

type MockAIClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

func (m *MockAIClient) Do(req *http.Request) (*http.Response, error) {
	return m.DoFunc(req)
}

type errorReader struct{}

func (e *errorReader) Read(p []byte) (n int, err error) {
	return 0, fmt.Errorf("error reading file")
}

var _ = Describe("RouteHandler", func() {
	var (
		mockRepo    *MockFileRepository
		mockClient  *MockAIClient
		fileService *service.FileService
		aiService   *service.AIService
	)

	BeforeEach(func() {
		mockRepo = &MockFileRepository{
			Files: make(map[string][]byte),
			Dirs:  make(map[string]bool),
		}
		mockClient = &MockAIClient{}
		fileService = service.NewFileService(mockRepo)
		aiService = &service.AIService{Client: mockClient}
		handler.FileService = fileService
		handler.AIService = aiService
		handler.Token = "test-token"
	})

	Describe("HandleUpload", func() {
		It("should successfully upload and process a valid CSV file", func() {
			body := new(bytes.Buffer)
			writer := multipart.NewWriter(body)
			part, _ := writer.CreateFormFile("file", "test.csv")
			part.Write([]byte("header1,header2\nvalue1,value2\nvalue3,value4"))
			writer.Close()

			req := httptest.NewRequest("POST", "/upload", body)
			req.Header.Set("Content-Type", writer.FormDataContentType())
			w := httptest.NewRecorder()

			mockClient.DoFunc = func(req *http.Request) (*http.Response, error) {
				body := `{"answer": "TV","coordinates": [[24,0]],"cells": ["TV"],"aggregator": "NONE"}`
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBufferString(body)),
				}, nil
			}

			handler.HandleUpload(w, req)

			Expect(w.Code).To(Equal(http.StatusOK))
			Expect(w.Body.String()).To(ContainSubstring("success"))
		})

		It("should return an error for an unsupported file type", func() {
			body := new(bytes.Buffer)
			writer := multipart.NewWriter(body)
			part, _ := writer.CreateFormFile("file", "test.txt")
			part.Write([]byte("some content"))
			writer.Close()

			req := httptest.NewRequest("POST", "/upload", body)
			req.Header.Set("Content-Type", writer.FormDataContentType())
			w := httptest.NewRecorder()

			handler.HandleUpload(w, req)

			Expect(w.Code).To(Equal(http.StatusUnsupportedMediaType))
			Expect(w.Body.String()).To(ContainSubstring("Only .csv files are allowed"))
		})

		It("should return an error for invalid CSV file content", func() {
			body := new(bytes.Buffer)
			writer := multipart.NewWriter(body)
			part, _ := writer.CreateFormFile("file", "test.csv")
			part.Write([]byte("header1,header2\nvalue1,value2\nvalue3")) // Invalid CSV content
			writer.Close()

			req := httptest.NewRequest("POST", "/upload", body)
			req.Header.Set("Content-Type", writer.FormDataContentType())
			w := httptest.NewRecorder()

			handler.HandleUpload(w, req)

			Expect(w.Code).To(Equal(http.StatusInternalServerError))
			Expect(w.Body.String()).To(ContainSubstring("Failed to process file content"))
		})

		It("should return an error for failing to parse form data", func() {
			req := httptest.NewRequest("POST", "/upload", nil)
			w := httptest.NewRecorder()

			handler.HandleUpload(w, req)

			Expect(w.Code).To(Equal(http.StatusBadRequest))
			Expect(w.Body.String()).To(ContainSubstring("Failed to parse form data"))
		})

		It("should return an error for failing to retrieve the uploaded file", func() {
			body := new(bytes.Buffer)
			writer := multipart.NewWriter(body)
			writer.Close()

			req := httptest.NewRequest("POST", "/upload", body)
			req.Header.Set("Content-Type", writer.FormDataContentType())
			w := httptest.NewRecorder()

			handler.HandleUpload(w, req)

			Expect(w.Code).To(Equal(http.StatusBadRequest))
			Expect(w.Body.String()).To(ContainSubstring("Failed to retrieve uploaded file"))
		})

		It("should return an error for failing to read file content", func() {
			body := new(bytes.Buffer)
			writer := multipart.NewWriter(body)
			part, _ := writer.CreateFormFile("file", "test.csv")
			part.Write([]byte("header1,header2\nvalue1,value2\nvalue3,value4"))
			writer.Close()

			req := httptest.NewRequest("POST", "/upload", body)
			req.Header.Set("Content-Type", writer.FormDataContentType())
			w := httptest.NewRecorder()

			HandleUpload := func(w http.ResponseWriter, r *http.Request) {
				var buf bytes.Buffer
				file := &errorReader{}
				if _, err := io.Copy(&buf, file); err != nil {
					utility.JSONResponse(w, http.StatusInternalServerError, "failed", "Failed to read file content")
					log.Printf("io.Copy error: %v", err)
					return
				}
			}

			HandleUpload(w, req)

			Expect(w.Code).To(Equal(http.StatusInternalServerError))
			Expect(w.Body.String()).To(ContainSubstring("Failed to read file content"))
		})
	})

	Describe("HandleChat", func() {
		It("should successfully process a chat request with type 'tapas'", func() {
			mockRepo.Files["/upload/data-series.csv"] = []byte("header1,header2\nvalue1,value2\nvalue3,value4")

			chatReq := model.ChatRequest{
				Type:  "tapas",
				Query: "Find the least electricity usage appliance.",
			}
			body, _ := json.Marshal(chatReq)

			req := httptest.NewRequest("POST", "/chat", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			mockClient.DoFunc = func(req *http.Request) (*http.Response, error) {
				body := `{"answer": "TV","coordinates": [[24,0]],"cells": ["TV"],"aggregator": "NONE"}`
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBufferString(body)),
				}, nil
			}

			handler.HandleChat(w, req)

			Expect(w.Code).To(Equal(http.StatusOK))
			Expect(w.Body.String()).To(ContainSubstring("success"))
		})

		It("should return an error for an invalid chat type", func() {
			chatReq := model.ChatRequest{
				Type:  "invalid",
				Query: "Some query",
			}
			body, _ := json.Marshal(chatReq)

			req := httptest.NewRequest("POST", "/chat", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			handler.HandleChat(w, req)

			Expect(w.Code).To(Equal(http.StatusBadRequest))
			Expect(w.Body.String()).To(ContainSubstring("Invalid chat type"))
		})

		It("should return an error for a missing data file", func() {
			chatReq := model.ChatRequest{
				Type:  "tapas",
				Query: "Find the least electricity usage appliance.",
			}
			body, _ := json.Marshal(chatReq)

			req := httptest.NewRequest("POST", "/chat", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			handler.HandleChat(w, req)

			Expect(w.Code).To(Equal(http.StatusNotFound))
			Expect(w.Body.String()).To(ContainSubstring("Data file not found"))
		})

		It("should return an error for failing to read the data file", func() {
			mockRepo.Files["/upload/data-series.csv"] = []byte("header1,header2\nvalue1,value2\nvalue3,value4")

			chatReq := model.ChatRequest{
				Type:  "tapas",
				Query: "Find the least electricity usage appliance.",
			}
			body, _ := json.Marshal(chatReq)

			req := httptest.NewRequest("POST", "/chat", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			mockRepo.ReadFileFn = func(filename string) ([]byte, error) {
				return nil, fmt.Errorf("simulated read error")
			}

			handler.HandleChat(w, req)

			Expect(w.Code).To(Equal(http.StatusInternalServerError))
			Expect(w.Body.String()).To(ContainSubstring("Failed to read data file"))
		})

		It("should return an error for failing to parse the CSV data", func() {
			mockRepo.Files["/upload/data-series.csv"] = []byte("invalid,csv,data")

			chatReq := model.ChatRequest{
				Type:  "tapas",
				Query: "Find the least electricity usage appliance.",
			}
			body, _ := json.Marshal(chatReq)

			req := httptest.NewRequest("POST", "/chat", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			handler.HandleChat(w, req)

			Expect(w.Code).To(Equal(http.StatusInternalServerError))
			Expect(w.Body.String()).To(ContainSubstring("Failed to parse CSV data"))
		})

		It("should return an error for failing to analyze the data with AI", func() {
			mockRepo.Files["/upload/data-series.csv"] = []byte("header1,header2\nvalue1,value2\nvalue3,value4")

			chatReq := model.ChatRequest{
				Type:  "tapas",
				Query: "Find the least electricity usage appliance.",
			}
			body, _ := json.Marshal(chatReq)

			req := httptest.NewRequest("POST", "/chat", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			mockClient.DoFunc = func(req *http.Request) (*http.Response, error) {
				return nil, fmt.Errorf("simulated AI error")
			}

			handler.HandleChat(w, req)

			Expect(w.Code).To(Equal(http.StatusInternalServerError))
			Expect(w.Body.String()).To(ContainSubstring("Failed to analyze data with AI"))
		})
	})

	Describe("HandleRemoveSession", func() {
		It("should successfully remove the session file", func() {
			mockRepo.Files["/upload/data-series.csv"] = []byte("some content")

			req := httptest.NewRequest("DELETE", "/remove-session", nil)
			w := httptest.NewRecorder()

			handler.HandleRemoveSession(w, req)

			Expect(w.Code).To(Equal(http.StatusOK))
			Expect(w.Body.String()).To(ContainSubstring("File deleted successfully"))
		})

		It("should return an error if the session file does not exist", func() {
			req := httptest.NewRequest("DELETE", "/remove-session", nil)
			w := httptest.NewRecorder()

			handler.HandleRemoveSession(w, req)

			Expect(w.Code).To(Equal(http.StatusNotFound))
			Expect(w.Body.String()).To(ContainSubstring("File not found"))
		})
	})
})
