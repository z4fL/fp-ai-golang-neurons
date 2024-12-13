package service

import (
	"encoding/csv"
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/z4fL/fp-ai-golang-neurons/repository"
	"github.com/z4fL/fp-ai-golang-neurons/utility/projectpath"
)

type FileService interface {
	ProcessFile(fileContent string) (map[string][]string, error)
	ParseCSV(fileContent string) (map[string][]string, error)
	GetRepo() repository.FileRepository
}

type fileService struct {
	repo repository.FileRepository
}

func NewFileService(repo repository.FileRepository) FileService {
	return &fileService{repo}
}

func (s *fileService) ProcessFile(fileContent string) (map[string][]string, error) {
	if strings.TrimSpace(fileContent) == "" {
		return nil, errors.New("file content is empty")
	}

	dir := filepath.Join(projectpath.Root, "upload")
	if !s.repo.DirExists(dir) {
		if err := s.repo.MakeDir(dir); err != nil {
			return nil, err
		}
	}

	filename := "data-series.csv"
	filePath := dir + "/" + filename

	var contentFile string

	if s.repo.FileExists(filePath) {
		existingContent, err := s.repo.ReadFile(filePath)
		if err != nil {
			return nil, fmt.Errorf("error reading file: %v", err)
		}

		if strings.TrimSpace(fileContent) == strings.TrimSpace(string(existingContent)) {
			contentFile = string(existingContent)
		} else {
			contentFile = fileContent
			err := s.repo.SaveFile(filePath, []byte(fileContent))
			if err != nil {
				return nil, fmt.Errorf("failed to save file")
			}
		}
	} else {
		err := s.repo.SaveFile(filePath, []byte(fileContent))
		if err != nil {
			return nil, fmt.Errorf("failed to save file")
		}

		contentFile = fileContent
	}

	parsedData, err := s.ParseCSV(contentFile)
	if err != nil {
		return nil, fmt.Errorf("error parsing CSV: %v", err)
	}

	return parsedData, nil
}

func (s *fileService) ParseCSV(fileContent string) (map[string][]string, error) {
	reader := csv.NewReader(strings.NewReader(fileContent))

	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("invalid CSV data")
	}

	if len(records) <= 1 {
		return nil, errors.New("CSV does not contain data")
	}

	parsedData := make(map[string][]string)
	headers := records[0]

	for i := 1; i < len(records); i++ {
		if len(records[i]) != len(headers) {
			return nil, fmt.Errorf("invalid CSV data at line %d", i+1)
		}

		for ii, header := range headers {
			parsedData[header] = append(parsedData[header], records[i][ii])
		}
	}

	return parsedData, nil
}

func (s *fileService) GetRepo() repository.FileRepository {
	return s.repo
}
