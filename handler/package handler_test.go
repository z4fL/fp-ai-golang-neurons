package handler_test

import (
	"bytes"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/z4fL/fp-ai-golang-neurons/handler"
	repository "github.com/z4fL/fp-ai-golang-neurons/repository/fileRepository"
	"github.com/z4fL/fp-ai-golang-neurons/service"
)

func TestRouteHandler(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "RouteHandler Suite")
}

type MockClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

func (m *MockClient) Do(req *http.Request) (*http.Response, error) {
	return m.DoFunc(req)
}

var _ = Describe("RouteHandler", func() {
	var (
		mockClient *MockClient
		recorder   *httptest.ResponseRecorder
		request    *http.Request
	)

	BeforeEach(func() {
		mockClient = &MockClient{}
		handler.AIService = &service.AIService{Client: mockClient}
		recorder = httptest.NewRecorder()
	})

	Describe("HandleUpload", func() {
		Context("when the file upload is successful", func() {
			BeforeEach(func() {
				body := &bytes.Buffer{}
				writer := multipart.NewWriter(body)
				part, _ := writer.CreateFormFile("file", "test.csv")
				part.Write([]byte("header1,header2\nvalue1,value2\nvalue3,value4"))
				writer.Close()

				request = httptest.NewRequest("POST", "/upload", body)
				request.Header.Set("Content-Type", writer.FormDataContentType())

				handler.FileService = &service.FileService{
					Repo: &repository.FileRepository{},
				}

				mockClient.DoFunc = func(req *http.Request) (*http.Response, error) {
					response := &http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(strings.NewReader(`{"answer": "TV"}`)),
					}
					return response, nil
				}
			})

			It("should return status success", func() {
				handler.HandleUpload(recorder, request)
				Expect(recorder.Code).To(Equal(http.StatusOK))
				Expect(recorder.Body.String()).To(ContainSubstring(`"status":"success"`))
			})
		})

		Context("when the file type is unsupported", func() {
			BeforeEach(func() {
				body := &bytes.Buffer{}
				writer := multipart.NewWriter(body)
				part, _ := writer.CreateFormFile("file", "test.txt")
				part.Write([]byte("header1,header2\nvalue1,value2\nvalue3,value4"))
				writer.Close()

				request = httptest.NewRequest("POST", "/upload", body)
				request.Header.Set("Content-Type", writer.FormDataContentType())
			})

			It("should return status unsupported media type", func() {
				handler.HandleUpload(recorder, request)
				Expect(recorder.Code).To(Equal(http.StatusUnsupportedMediaType))
				Expect(recorder.Body.String()).To(ContainSubstring(`"status":"failed"`))
			})
		})

		Context("when the form data is invalid", func() {
			BeforeEach(func() {
				body := &bytes.Buffer{}
				writer := multipart.NewWriter(body)
				writer.Close()

				request = httptest.NewRequest("POST", "/upload", body)
				request.Header.Set("Content-Type", writer.FormDataContentType())
			})

			It("should return status bad request", func() {
				handler.HandleUpload(recorder, request)
				Expect(recorder.Code).To(Equal(http.StatusBadRequest))
				Expect(recorder.Body.String()).To(ContainSubstring(`"status":"failed"`))
			})
		})

		Context("when there is an error reading the file", func() {
			BeforeEach(func() {
				body := &bytes.Buffer{}
				writer := multipart.NewWriter(body)
				part, _ := writer.CreateFormFile("file", "test.csv")
				part.Write([]byte("header1,header2\nvalue1,value2\nvalue3,value4"))
				writer.Close()

				request = httptest.NewRequest("POST", "/upload", body)
				request.Header.Set("Content-Type", writer.FormDataContentType())

				handler.FileService = &service.FileService{
					Repo: &repository.FileRepository{},
				}

				mockClient.DoFunc = func(req *http.Request) (*http.Response, error) {
					return nil, errors.New("failed to reading data")
				}
			})

			It("should return status internal server error", func() {
				handler.HandleUpload(recorder, request)
				Expect(recorder.Code).To(Equal(http.StatusInternalServerError))
				Expect(recorder.Body.String()).To(ContainSubstring(`"status":"failed"`))
			})
		})

		Context("when there is an error processing the file", func() {
			BeforeEach(func() {
				body := &bytes.Buffer{}
				writer := multipart.NewWriter(body)
				part, _ := writer.CreateFormFile("file", "test.csv")
				part.Write([]byte("header1,header2\nvalue1,value2\nvalue3,value4\nvalue5,value6"))
				writer.Close()

				request = httptest.NewRequest("POST", "/upload", body)
				request.Header.Set("Content-Type", writer.FormDataContentType())

				handler.FileService = &service.FileService{
					Repo: &repository.FileRepository{},
				}

				mockClient.DoFunc = func(req *http.Request) (*http.Response, error) {
					return nil, errors.New("failed to analyze data")
				}
			})

			It("should return status internal server error", func() {
				handler.HandleUpload(recorder, request)
				Expect(recorder.Code).To(Equal(http.StatusInternalServerError))
				Expect(recorder.Body.String()).To(ContainSubstring(`"status":"failed"`))
			})
		})
	})
})
