package service

import (
	"encoding/csv"
	"errors"
	"fmt"
	"log"
	"path/filepath"
	"strings"

	repository "github.com/z4fL/fp-ai-golang-neurons/repository/fileRepository"
	"github.com/z4fL/fp-ai-golang-neurons/utility/projectpath"
)

type FileService struct {
	Repo *repository.FileRepository
}

func (s *FileService) ProcessFile(fileContent string) (map[string][]string, error) {
	if strings.TrimSpace(fileContent) == "" {
		return nil, errors.New("file content is empty")
	}

	dir := filepath.Join(projectpath.Root, "upload")
	if !s.Repo.DirExists(dir) {
		if err := s.Repo.MakeDir(dir); err != nil {
			return nil, err
		}
	}

	log.Println(dir)

	filename := "data-series.csv"
	filePath := dir + "/" + filename

	var contentFile string

	if s.Repo.FileExists(filePath) {
		existingContent, err := s.Repo.ReadFile(filePath)
		if err != nil {
			return nil, fmt.Errorf("error reading file: %v", err)
		}

		if strings.TrimSpace(fileContent) == strings.TrimSpace(string(existingContent)) {
			contentFile = string(existingContent)
		} else {
			contentFile = fileContent
			err := s.Repo.SaveFile(filePath, []byte(fileContent))
			if err != nil {
				return nil, fmt.Errorf("error saving file: %v", err)
			}
		}
	} else {
		err := s.Repo.SaveFile(filePath, []byte(fileContent))
		if err != nil {
			return nil, fmt.Errorf("error saving file: %v", err)
		}

		contentFile = fileContent
	}

	parsedData, err := s.ParseCSV(contentFile)
	if err != nil {
		return nil, fmt.Errorf("error parsing CSV: %v", err)
	}

	return parsedData, nil
}

func (s *FileService) ParseCSV(fileContent string) (map[string][]string, error) {
	reader := csv.NewReader(strings.NewReader(fileContent))

	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("error reading CSV: %v", err)
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
