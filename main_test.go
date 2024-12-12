package main_test

import (
	"bytes"
	"io"
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/z4fL/fp-ai-golang-neurons/repository"
	"github.com/z4fL/fp-ai-golang-neurons/service"
)

var _ = Describe("FileService", func() {
	var fileService service.FileService

	BeforeEach(func() {
		fileRepo := repository.NewFileRepository()

		fileService = service.NewFileService(fileRepo)
	})

	Describe("ProcessFile", func() {
		It("should return the correct result for valid CSV data", func() {
			fileContent := "header1,header2\nvalue1,value2\nvalue3,value4"
			expected := map[string][]string{
				"header1": {"value1", "value3"},
				"header2": {"value2", "value4"},
			}

			result, err := fileService.ProcessFile(fileContent)
			Expect(err).ToNot(HaveOccurred())
			Expect(result).To(Equal(expected))
		})

		It("should return an error for empty CSV data", func() {
			fileContent := ``

			result, err := fileService.ProcessFile(fileContent)
			Expect(err).To(HaveOccurred())
			Expect(result).To(BeNil())
		})

		It("should return an error for invalid CSV data", func() {
			fileContent := "header1,header2\nvalue1,value2\nvalue3"

			result, err := fileService.ProcessFile(fileContent)
			Expect(err).To(HaveOccurred())
			Expect(result).To(BeNil())
		})
	})
})

type MockClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

func (m *MockClient) Do(req *http.Request) (*http.Response, error) {
	return m.DoFunc(req)
}

var _ = Describe("AIService", func() {
	var (
		mockClient *MockClient
		aiService  service.AIService
	)

	BeforeEach(func() {
		mockClient = &MockClient{}
		aiService = service.NewAIService(mockClient)
	})

	Describe("AnalyzeData", func() {
		It("should return the correct result for a valid response", func() {
			mockClient.DoFunc = func(req *http.Request) (*http.Response, error) {
				body := `{"answer": "TV","coordinates": [[24,0]],"cells": ["TV"],"aggregator": "NONE"}`
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBufferString(body)),
				}, nil
			}

			table := map[string][]string{
				"header1": {"value1", "value2"},
			}
			query := "query"
			token := "token"

			result, err := aiService.AnalyzeData(table, query, token)
			Expect(err).ToNot(HaveOccurred())
			Expect(result).To(Equal("TV"))
		})

		It("should return an error for an empty table", func() {
			table := map[string][]string{}
			query := "query"
			token := "token"

			result, err := aiService.AnalyzeData(table, query, token)
			Expect(err).To(HaveOccurred())
			Expect(result).To(BeEmpty())
		})

		It("should return an error for an error response", func() {
			mockClient.DoFunc = func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusInternalServerError,
					Body:       io.NopCloser(bytes.NewBufferString(`{"error":"internal error"}`)),
				}, nil
			}

			table := map[string][]string{
				"header1": {"value1", "value2"},
			}
			query := "query"
			token := "token"

			result, err := aiService.AnalyzeData(table, query, token)
			Expect(err).To(HaveOccurred())
			Expect(result).To(BeEmpty())
		})
	})

	Describe("ChatWithAI", func() {
		It("should return the correct response for a valid request", func() {
			mockClient.DoFunc = func(req *http.Request) (*http.Response, error) {
				body := `{"choices":[{"message":{"role":"assistant","content":"response"}}]}`
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBufferString(body)),
				}, nil
			}

			context := "context"
			query := "query"
			token := "token"

			result, err := aiService.ChatWithAI(context, query, token)
			Expect(err).ToNot(HaveOccurred())
			Expect(result).To(Equal("response"))
		})

		It("should return an error for an error response", func() {
			mockClient.DoFunc = func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusInternalServerError,
					Body:       io.NopCloser(bytes.NewBufferString(`{"error":"internal error"}`)),
				}, nil
			}

			context := "context"
			query := "query"
			token := "token"

			result, err := aiService.ChatWithAI(context, query, token)
			Expect(err).To(HaveOccurred())
			Expect(result).To(BeEmpty())
		})
	})
})
