package storage

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
)

type FileStorage interface {
	SaveFile(file *multipart.FileHeader, uploadPath string) (string, string, error)
	DeleteFile(filePath string) error
	GetFileURL(filePath, baseURL string) string
}

type fileStorage struct{}

func NewFileStorage() FileStorage {
	return &fileStorage{}
}

func (fs *fileStorage) SaveFile(file *multipart.FileHeader, uploadPath string) (string, string, error) {
	// Generate unique filename
	ext := filepath.Ext(file.Filename)
	uniqueID := uuid.New().String()
	timestamp := time.Now().Format("20060102-150405")
	fileName := fmt.Sprintf("%s-%s%s", timestamp, uniqueID, ext)
	
	// Create full path
	fullPath := filepath.Join(uploadPath, fileName)
	
	// Create directory if it doesn't exist
	if err := os.MkdirAll(uploadPath, 0755); err != nil {
		return "", "", err
	}
	
	// Open uploaded file
	src, err := file.Open()
	if err != nil {
		return "", "", err
	}
	defer src.Close()
	
	// Create destination file
	dst, err := os.Create(fullPath)
	if err != nil {
		return "", "", err
	}
	defer dst.Close()
	
	// Copy file content
	_, err = io.Copy(dst, src)
	if err != nil {
		return "", "", err
	}
	
	return fileName, fullPath, nil
}

func (fs *fileStorage) DeleteFile(filePath string) error {
	return os.Remove(filePath)
}

func (fs *fileStorage) GetFileURL(filePath, baseURL string) string {
	return fmt.Sprintf("%s/files/%s", baseURL, filepath.Base(filePath))
}