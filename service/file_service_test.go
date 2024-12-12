package service_test

import (
	"errors"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/z4fL/fp-ai-golang-neurons/service"
)

type MockFileRepository struct {
	DirExistsFunc  func(path string) bool
	MakeDirFunc    func(path string) error
	FileExistsFunc func(path string) bool
	ReadFileFunc   func(path string) ([]byte, error)
	SaveFileFunc   func(path string, content []byte) error
	RemoveFileFunc func(filename string) error
}

func (m *MockFileRepository) RemoveFile(filename string) error {
	return m.RemoveFileFunc(filename)
}

func (m *MockFileRepository) DirExists(path string) bool {
	return m.DirExistsFunc(path)
}

func (m *MockFileRepository) MakeDir(path string) error {
	return m.MakeDirFunc(path)
}

func (m *MockFileRepository) FileExists(path string) bool {
	return m.FileExistsFunc(path)
}

func (m *MockFileRepository) ReadFile(path string) ([]byte, error) {
	return m.ReadFileFunc(path)
}

func (m *MockFileRepository) SaveFile(path string, content []byte) error {
	return m.SaveFileFunc(path, content)
}

var _ = Describe("FileService", func() {
	var (
		mockRepo    *MockFileRepository
		fileService service.FileService
	)

	BeforeEach(func() {
		mockRepo = &MockFileRepository{}
		fileService = service.NewFileService(mockRepo)
	})

	Describe("ProcessFile", func() {
		It("should return an error if the file content is empty", func() {
			_, err := fileService.ProcessFile("")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("file content is empty"))
		})

		It("should create directory if it does not exist", func() {
			mockRepo.DirExistsFunc = func(path string) bool { return false }
			mockRepo.MakeDirFunc = func(path string) error { return nil }
			mockRepo.FileExistsFunc = func(path string) bool { return false }
			mockRepo.SaveFileFunc = func(path string, content []byte) error { return nil }
			mockRepo.ReadFileFunc = func(path string) ([]byte, error) { return nil, nil }

			_, err := fileService.ProcessFile("header1,header2\nvalue1,value2")
			Expect(err).NotTo(HaveOccurred())
		})

		It("should return an error if directory creation fails", func() {
			mockRepo.DirExistsFunc = func(path string) bool { return false }
			mockRepo.MakeDirFunc = func(path string) error { return errors.New("failed to create directory") }

			_, err := fileService.ProcessFile("header1,header2\nvalue1,value2")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("failed to create directory"))
		})

		It("should save file if it does not exist", func() {
			mockRepo.DirExistsFunc = func(path string) bool { return true }
			mockRepo.FileExistsFunc = func(path string) bool { return false }
			mockRepo.SaveFileFunc = func(path string, content []byte) error { return nil }

			_, err := fileService.ProcessFile("header1,header2\nvalue1,value2")
			Expect(err).NotTo(HaveOccurred())
		})

		It("should return an error if file saving fails", func() {
			mockRepo.DirExistsFunc = func(path string) bool { return true }
			mockRepo.FileExistsFunc = func(path string) bool { return false }
			mockRepo.SaveFileFunc = func(path string, content []byte) error { return errors.New("failed to save file") }

			_, err := fileService.ProcessFile("header1,header2\nvalue1,value2")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("failed to save file"))
		})

		It("should parse CSV content correctly", func() {
			mockRepo.DirExistsFunc = func(path string) bool { return true }
			mockRepo.FileExistsFunc = func(path string) bool { return false }
			mockRepo.SaveFileFunc = func(path string, content []byte) error { return nil }

			parsedData, err := fileService.ProcessFile("header1,header2\nvalue1,value2")
			Expect(err).NotTo(HaveOccurred())
			Expect(parsedData).To(HaveKeyWithValue("header1", []string{"value1"}))
			Expect(parsedData).To(HaveKeyWithValue("header2", []string{"value2"}))
		})
	})

	Describe("ParseCSV", func() {
		It("should return an error if CSV content is invalid", func() {
			_, err := fileService.ParseCSV("header1,header2\nvalue1")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("invalid CSV data"))
		})

		It("should return an error if CSV does not contain data", func() {
			_, err := fileService.ParseCSV("header1,header2")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("CSV does not contain data"))
		})

		It("should parse valid CSV content correctly", func() {
			parsedData, err := fileService.ParseCSV("header1,header2\nvalue1,value2")
			Expect(err).NotTo(HaveOccurred())
			Expect(parsedData).To(HaveKeyWithValue("header1", []string{"value1"}))
			Expect(parsedData).To(HaveKeyWithValue("header2", []string{"value2"}))
		})
	})
})
