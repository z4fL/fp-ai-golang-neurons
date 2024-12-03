package service

import (
	repository "github.com/z4fL/fp-ai-golang-neurons/repository/fileRepository"
)

type FileService struct {
	Repo *repository.FileRepository
}

func (s *FileService) ProcessFile(fileContent string) (map[string][]string, error) {
	return nil, nil
}
