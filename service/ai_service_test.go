package service_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/z4fL/fp-ai-golang-neurons/model"
	"github.com/z4fL/fp-ai-golang-neurons/service"
)

type MockHTTPClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return m.DoFunc(req)
}

var _ = Describe("AIService", func() {
	var (
		mockClient *MockHTTPClient
		aiService  service.AIService
		token      string
	)

	BeforeEach(func() {
		mockClient = &MockHTTPClient{}
		aiService = service.NewAIService(mockClient)
		token = "test-token"
	})

	Describe("AnalyzeData", func() {
		It("should return an error if the table is empty", func() {
			result, err := aiService.AnalyzeData(map[string][]string{}, "query", token)
			Expect(err).To(HaveOccurred())
			Expect(result).To(BeEmpty())
		})

		It("should return a valid response for a valid request", func() {
			mockResponse := model.TapasResponse{
				Cells:      []string{"cell1", "cell2"},
				Answer:     "cell1, cell2",
				Aggregator: "NONE",
			}
			mockClient.DoFunc = func(req *http.Request) (*http.Response, error) {
				responseBody, _ := json.Marshal(mockResponse)
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBuffer(responseBody)),
				}, nil
			}

			table := map[string][]string{"column1": {"value1", "value2"}}
			result, err := aiService.AnalyzeData(table, "query", token)
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(Equal("Count: 2, List: cell1, cell2"))
		})

		It("should return an error if the API response is not OK", func() {
			mockClient.DoFunc = func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusInternalServerError,
					Body:       io.NopCloser(bytes.NewBufferString("")),
				}, nil
			}

			table := map[string][]string{"column1": {"value1", "value2"}}
			result, err := aiService.AnalyzeData(table, "query", token)
			Expect(err).To(HaveOccurred())
			Expect(result).To(BeEmpty())
		})
	})

	Describe("AnalyzeFile", func() {
		It("should return a valid response for multiple queries", func() {
			mockResponse := model.TapasResponse{
				Cells:      []string{"cell1", "cell2"},
				Answer:     "answer",
				Aggregator: "NONE",
			}
			mockClient.DoFunc = func(req *http.Request) (*http.Response, error) {
				responseBody, _ := json.Marshal(mockResponse)
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBuffer(responseBody)),
				}, nil
			}

			table := map[string][]string{"column1": {"value1", "value2"}}
			queries := []string{"query1", "query2"}
			result, err := aiService.AnalyzeFile(table, queries, token)
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(ContainSubstring("Least Electricity"))
			Expect(result).To(ContainSubstring("Most Electricity"))
		})
	})

	Describe("ChatWithAI", func() {
		It("should return a valid response for a chat request", func() {
			mockResponse := model.PhiResponse{
				Choices: []model.Choice{
					{Message: model.Message{Content: "response"}},
				},
			}
			mockClient.DoFunc = func(req *http.Request) (*http.Response, error) {
				responseBody, _ := json.Marshal(mockResponse)
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBuffer(responseBody)),
				}, nil
			}

			result, err := aiService.ChatWithAI("context", "query", token)
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(Equal("response"))
		})

		It("should return an error if the API response is not OK", func() {
			mockClient.DoFunc = func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusInternalServerError,
					Body:       io.NopCloser(bytes.NewBufferString("")),
				}, nil
			}

			result, err := aiService.ChatWithAI("context", "query", token)
			Expect(err).To(HaveOccurred())
			Expect(result).To(BeEmpty())
		})
	})
})
