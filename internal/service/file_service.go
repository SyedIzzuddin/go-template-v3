package service

import (
	"context"
	"database/sql"
	"errors"
	"go-template/internal/config"
	"go-template/internal/dto"
	"go-template/internal/entity"
	"go-template/internal/logger"
	"go-template/internal/repository"
	"go-template/pkg/storage"
	"mime/multipart"
	"slices"

	"go.uber.org/zap"
)

type FileService interface {
	UploadFile(ctx context.Context, file *multipart.FileHeader, req dto.UploadFileRequest, userID int) (*dto.FileResponse, error)
	GetFileByID(ctx context.Context, id int) (*dto.FileResponse, error)
	GetFilesByUserID(ctx context.Context, userID int) ([]dto.FileResponse, error)
	GetAllFiles(ctx context.Context) ([]dto.FileResponse, error)
	UpdateFile(ctx context.Context, id int, req dto.UpdateFileRequest) (*dto.FileResponse, error)
	DeleteFile(ctx context.Context, id int) error
	GetFileEntity(ctx context.Context, id int) (*entity.File, error)
}

type fileService struct {
	fileRepo    repository.FileRepository
	fileStorage storage.FileStorage
	config      *config.Config
}

func NewFileService(fileRepo repository.FileRepository, fileStorage storage.FileStorage, config *config.Config) FileService {
	return &fileService{
		fileRepo:    fileRepo,
		fileStorage: fileStorage,
		config:      config,
	}
}

func (s *fileService) UploadFile(ctx context.Context, file *multipart.FileHeader, req dto.UploadFileRequest, userID int) (*dto.FileResponse, error) {
	logger.Info("Uploading file", zap.String("original_name", file.Filename), zap.Int("user_id", userID))
	
	// Validate file size
	if file.Size > s.config.Upload.MaxFileSize {
		logger.Warn("File size exceeds limit", zap.Int64("size", file.Size), zap.Int64("max_size", s.config.Upload.MaxFileSize))
		return nil, errors.New("file size exceeds maximum allowed size")
	}
	
	// Validate file type
	if !slices.Contains(s.config.Upload.AllowedTypes, file.Header.Get("Content-Type")) {
		logger.Warn("Invalid file type", zap.String("mime_type", file.Header.Get("Content-Type")))
		return nil, errors.New("file type not allowed")
	}
	
	// Save file to storage
	fileName, filePath, err := s.fileStorage.SaveFile(file, s.config.Upload.UploadPath)
	if err != nil {
		logger.Error("Failed to save file", zap.Error(err))
		return nil, err
	}
	
	// Save to database
	fileEntity, err := s.fileRepo.Create(ctx, fileName, file.Filename, filePath, file.Size, file.Header.Get("Content-Type"), req.Description, req.Category, userID)
	if err != nil {
		// Delete file if database save fails
		s.fileStorage.DeleteFile(filePath)
		logger.Error("Failed to save file to database", zap.Error(err))
		return nil, err
	}
	
	logger.Info("File uploaded successfully", zap.Int("file_id", fileEntity.ID))
	
	return s.mapFileToResponse(fileEntity), nil
}

func (s *fileService) GetFileByID(ctx context.Context, id int) (*dto.FileResponse, error) {
	logger.Debug("Getting file by ID", zap.Int("file_id", id))
	
	file, err := s.fileRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.Warn("File not found", zap.Int("file_id", id))
			return nil, errors.New("file not found")
		}
		logger.Error("Failed to get file", zap.Error(err))
		return nil, err
	}
	
	return s.mapFileToResponse(file), nil
}

func (s *fileService) GetFilesByUserID(ctx context.Context, userID int) ([]dto.FileResponse, error) {
	logger.Debug("Getting files by user ID", zap.Int("user_id", userID))
	
	files, err := s.fileRepo.GetByUserID(ctx, userID)
	if err != nil {
		logger.Error("Failed to get files by user ID", zap.Error(err))
		return nil, err
	}
	
	var fileResponses []dto.FileResponse
	for _, file := range files {
		fileResponses = append(fileResponses, *s.mapFileToResponse(&file))
	}
	
	return fileResponses, nil
}

func (s *fileService) GetAllFiles(ctx context.Context) ([]dto.FileResponse, error) {
	logger.Debug("Getting all files")
	
	files, err := s.fileRepo.GetAll(ctx)
	if err != nil {
		logger.Error("Failed to get all files", zap.Error(err))
		return nil, err
	}
	
	var fileResponses []dto.FileResponse
	for _, file := range files {
		fileResponses = append(fileResponses, *s.mapFileToResponse(&file))
	}
	
	return fileResponses, nil
}

func (s *fileService) UpdateFile(ctx context.Context, id int, req dto.UpdateFileRequest) (*dto.FileResponse, error) {
	logger.Info("Updating file", zap.Int("file_id", id))
	
	// Check if file exists
	_, err := s.fileRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.Warn("File not found for update", zap.Int("file_id", id))
			return nil, errors.New("file not found")
		}
		logger.Error("Failed to get file for update", zap.Error(err))
		return nil, err
	}
	
	// Update file
	file, err := s.fileRepo.Update(ctx, id, req.Description, req.Category)
	if err != nil {
		logger.Error("Failed to update file", zap.Error(err))
		return nil, err
	}
	
	logger.Info("File updated successfully", zap.Int("file_id", id))
	
	return s.mapFileToResponse(file), nil
}

func (s *fileService) DeleteFile(ctx context.Context, id int) error {
	logger.Info("Deleting file", zap.Int("file_id", id))
	
	file, err := s.fileRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.Warn("File not found for deletion", zap.Int("file_id", id))
			return errors.New("file not found")
		}
		logger.Error("Failed to get file for deletion", zap.Error(err))
		return err
	}
	
	// Delete from database
	if err := s.fileRepo.Delete(ctx, id); err != nil {
		logger.Error("Failed to delete file from database", zap.Error(err))
		return err
	}
	
	// Delete from storage
	if err := s.fileStorage.DeleteFile(file.FilePath); err != nil {
		logger.Error("Failed to delete file from storage", zap.Error(err))
		// Note: File is already deleted from database, but physical file remains
		// In production, you might want to have a cleanup job for orphaned files
	}
	
	logger.Info("File deleted successfully", zap.Int("file_id", id))
	
	return nil
}

func (s *fileService) GetFileEntity(ctx context.Context, id int) (*entity.File, error) {
	return s.fileRepo.GetByID(ctx, id)
}

func (s *fileService) mapFileToResponse(file *entity.File) *dto.FileResponse {
	return &dto.FileResponse{
		ID:           file.ID,
		FileName:     file.FileName,
		OriginalName: file.OriginalName,
		FilePath:     s.fileStorage.GetFileURL(file.FilePath, s.config.Upload.BaseURL),
		FileSize:     file.FileSize,
		MimeType:     file.MimeType,
		Description:  file.Description,
		Category:     file.Category,
		UploadedBy:   file.UploadedBy,
		CreatedAt:    file.CreatedAt,
		UpdatedAt:    file.UpdatedAt,
	}
}