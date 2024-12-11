package fileRepository

import (
	"fmt"
	"log"
	"os"
)

type FileRepository struct{}

// SaveFile saves the uploaded file content to the server's file system
func (r *FileRepository) SaveFile(filename string, content []byte) error {
	return os.WriteFile(filename, content, 0644)
}

// ReadFile reads the content of a file from the server's file system
func (r *FileRepository) ReadFile(filename string) ([]byte, error) {
	return os.ReadFile(filename)
}

// FileExists checks if a file already exists
func (r *FileRepository) FileExists(filename string) bool {
	_, err := os.Stat(filename)
	if err != nil && !os.IsNotExist(err) {
		log.Printf("Error checking file: %v\n", err)
	}
	return !os.IsNotExist(err)
}

func (r *FileRepository) RemoveFile(filename string) error {
	if !r.FileExists(filename) {
		return fmt.Errorf("file %s does not exist", filename)
	}
	return os.Remove(filename)
}

// DirExists checks if a directory already exists
func (r *FileRepository) DirExists(dirname string) bool {
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
func (r *FileRepository) MakeDir(dirname string) error {
	return os.MkdirAll(dirname, os.ModePerm)
}
