package dto

import "time"

type UploadFileRequest struct {
	Description string `form:"description" validate:"omitempty,max=500"`
	Category    string `form:"category" validate:"omitempty,max=50"`
}

type FileResponse struct {
	ID          int       `json:"id"`
	FileName    string    `json:"file_name"`
	OriginalName string   `json:"original_name"`
	FilePath    string    `json:"file_path"`
	FileSize    int64     `json:"file_size"`
	MimeType    string    `json:"mime_type"`
	Description string    `json:"description"`
	Category    string    `json:"category"`
	UploadedBy  int       `json:"uploaded_by"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type UpdateFileRequest struct {
	Description string `json:"description,omitempty" validate:"omitempty,max=500"`
	Category    string `json:"category,omitempty" validate:"omitempty,max=50"`
}