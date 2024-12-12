package repository

import (
	"fmt"
	"log"
	"os"
)

type FileRepository interface {
	SaveFile(filename string, content []byte) error
	ReadFile(filename string) ([]byte, error)
	FileExists(filename string) bool
	RemoveFile(filename string) error
	DirExists(dirname string) bool
	MakeDir(dirname string) error
}

type fileRepository struct{}

func NewFileRepository() FileRepository {
	return &fileRepository{}
}

// SaveFile saves the uploaded file content to the server's file system
func (r *fileRepository) SaveFile(filename string, content []byte) error {
	return os.WriteFile(filename, content, 0644)
}

// ReadFile reads the content of a file from the server's file system
func (r *fileRepository) ReadFile(filename string) ([]byte, error) {
	return os.ReadFile(filename)
}

// FileExists checks if a file already exists
func (r *fileRepository) FileExists(filename string) bool {
	_, err := os.Stat(filename)
	if err != nil && !os.IsNotExist(err) {
		log.Printf("Error checking file: %v\n", err)
	}
	return !os.IsNotExist(err)
}

func (r *fileRepository) RemoveFile(filename string) error {
	if !r.FileExists(filename) {
		return fmt.Errorf("file %s does not exist", filename)
	}
	return os.Remove(filename)
}

// DirExists checks if a directory already exists
func (r *fileRepository) DirExists(dirname string) bool {
	info, err := os.Stat(dirname)
	if err != nil {
		if !os.IsNotExist(err) {
			log.Printf("Error checking directory: %v\n", err)
		}
		return false
	}
	return info.IsDir()
}

// MakeDir creates a new directory with the specified name
func (r *fileRepository) MakeDir(dirname string) error {
	return os.MkdirAll(dirname, os.ModePerm)
}
