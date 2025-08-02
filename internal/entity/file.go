package entity

import (
	"time"
)

type File struct {
	ID           int       `json:"id"`
	FileName     string    `json:"file_name"`
	OriginalName string    `json:"original_name"`
	FilePath     string    `json:"file_path"`
	FileSize     int64     `json:"file_size"`
	MimeType     string    `json:"mime_type"`
	Description  string    `json:"description"`
	Category     string    `json:"category"`
	UploadedBy   int       `json:"uploaded_by"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}