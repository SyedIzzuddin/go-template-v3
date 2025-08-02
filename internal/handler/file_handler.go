package handler

import (
	"path/filepath"
	"strconv"

	"go-template/internal/dto"
	"go-template/internal/logger"
	"go-template/internal/service"
	"go-template/pkg/response"
	"go-template/pkg/validator"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type FileHandler struct {
	fileService service.FileService
	validator   *validator.Validator
}

func NewFileHandler(fileService service.FileService, validator *validator.Validator) *FileHandler {
	return &FileHandler{
		fileService: fileService,
		validator:   validator,
	}
}

func (h *FileHandler) UploadFile(c echo.Context) error {
	requestID := c.Response().Header().Get(echo.HeaderXRequestID)
	logger.Info("UploadFile request started", zap.String("request_id", requestID))

	// TODO: Make sure to edit this part and use the correct method to get the user id
	userID := 1

	// Get file from form
	file, err := c.FormFile("file")
	if err != nil {
		logger.Error("Failed to get file from form", zap.Error(err), zap.String("request_id", requestID))
		return response.BadRequest(c, "File is required", err.Error())
	}

	// Bind additional form data
	var req dto.UploadFileRequest
	if err := c.Bind(&req); err != nil {
		logger.Error("Failed to bind upload request", zap.Error(err), zap.String("request_id", requestID))
		return response.BadRequest(c, "Invalid request data", err.Error())
	}

	// Validate request
	if validationErrors := h.validator.ValidateStruct(req); validationErrors != nil {
		logger.Warn("Upload validation failed", zap.Any("errors", validationErrors), zap.String("request_id", requestID))
		return response.ValidationError(c, "Validation failed", validationErrors)
	}

	fileResponse, err := h.fileService.UploadFile(c.Request().Context(), file, req, userID)
	if err != nil {
		logger.Error("Failed to upload file", zap.Error(err), zap.String("request_id", requestID))
		return response.BadRequest(c, "Failed to upload file", err.Error())
	}

	logger.Info("UploadFile request completed", zap.String("request_id", requestID))
	return response.Created(c, "File uploaded successfully", fileResponse)
}

func (h *FileHandler) GetFile(c echo.Context) error {
	requestID := c.Response().Header().Get(echo.HeaderXRequestID)
	logger.Info("GetFile request started", zap.String("request_id", requestID))

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logger.Error("Invalid file ID", zap.Error(err), zap.String("request_id", requestID))
		return response.BadRequest(c, "Invalid file ID", err.Error())
	}

	file, err := h.fileService.GetFileByID(c.Request().Context(), id)
	if err != nil {
		logger.Error("Failed to get file", zap.Error(err), zap.String("request_id", requestID))
		return response.NotFound(c, "File not found")
	}

	logger.Info("GetFile request completed", zap.String("request_id", requestID))
	return response.Success(c, "File retrieved successfully", file)
}

func (h *FileHandler) GetMyFiles(c echo.Context) error {
	requestID := c.Response().Header().Get(echo.HeaderXRequestID)
	userID := 1 // Hardcoded for demo purposes
	logger.Info("GetMyFiles request started", zap.String("request_id", requestID), zap.Int("user_id", userID))

	files, err := h.fileService.GetFilesByUserID(c.Request().Context(), userID)
	if err != nil {
		logger.Error("Failed to get user files", zap.Error(err), zap.String("request_id", requestID))
		return response.InternalServerError(c, "Failed to get files", err.Error())
	}

	logger.Info("GetMyFiles request completed", zap.String("request_id", requestID))
	return response.Success(c, "Files retrieved successfully", files)
}

func (h *FileHandler) GetAllFiles(c echo.Context) error {
	requestID := c.Response().Header().Get(echo.HeaderXRequestID)
	logger.Info("GetAllFiles request started", zap.String("request_id", requestID))

	files, err := h.fileService.GetAllFiles(c.Request().Context())
	if err != nil {
		logger.Error("Failed to get all files", zap.Error(err), zap.String("request_id", requestID))
		return response.InternalServerError(c, "Failed to get files", err.Error())
	}

	logger.Info("GetAllFiles request completed", zap.String("request_id", requestID))
	return response.Success(c, "Files retrieved successfully", files)
}

func (h *FileHandler) UpdateFile(c echo.Context) error {
	requestID := c.Response().Header().Get(echo.HeaderXRequestID)
	logger.Info("UpdateFile request started", zap.String("request_id", requestID))

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logger.Error("Invalid file ID", zap.Error(err), zap.String("request_id", requestID))
		return response.BadRequest(c, "Invalid file ID", err.Error())
	}

	var req dto.UpdateFileRequest
	if err := c.Bind(&req); err != nil {
		logger.Error("Failed to bind update request", zap.Error(err), zap.String("request_id", requestID))
		return response.BadRequest(c, "Invalid request body", err.Error())
	}

	// Validate request
	if validationErrors := h.validator.ValidateStruct(req); validationErrors != nil {
		logger.Warn("Update validation failed", zap.Any("errors", validationErrors), zap.String("request_id", requestID))
		return response.ValidationError(c, "Validation failed", validationErrors)
	}

	file, err := h.fileService.UpdateFile(c.Request().Context(), id, req)
	if err != nil {
		logger.Error("Failed to update file", zap.Error(err), zap.String("request_id", requestID))
		return response.BadRequest(c, "Failed to update file", err.Error())
	}

	logger.Info("UpdateFile request completed", zap.String("request_id", requestID))
	return response.Success(c, "File updated successfully", file)
}

func (h *FileHandler) DeleteFile(c echo.Context) error {
	requestID := c.Response().Header().Get(echo.HeaderXRequestID)
	logger.Info("DeleteFile request started", zap.String("request_id", requestID))

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logger.Error("Invalid file ID", zap.Error(err), zap.String("request_id", requestID))
		return response.BadRequest(c, "Invalid file ID", err.Error())
	}

	err = h.fileService.DeleteFile(c.Request().Context(), id)
	if err != nil {
		logger.Error("Failed to delete file", zap.Error(err), zap.String("request_id", requestID))
		return response.BadRequest(c, "Failed to delete file", err.Error())
	}

	logger.Info("DeleteFile request completed", zap.String("request_id", requestID))
	return response.Success(c, "File deleted successfully", nil)
}

func (h *FileHandler) DownloadFile(c echo.Context) error {
	requestID := c.Response().Header().Get(echo.HeaderXRequestID)
	logger.Info("DownloadFile request started", zap.String("request_id", requestID))

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logger.Error("Invalid file ID", zap.Error(err), zap.String("request_id", requestID))
		return response.BadRequest(c, "Invalid file ID", err.Error())
	}

	// Get file entity directly for download
	file, err := h.fileService.GetFileEntity(c.Request().Context(), id)
	if err != nil {
		logger.Error("Failed to get file entity", zap.Error(err), zap.String("request_id", requestID))
		return response.NotFound(c, "File not found")
	}

	// Set headers for file download
	c.Response().Header().Set("Content-Disposition", "attachment; filename=\""+file.OriginalName+"\"")
	c.Response().Header().Set("Content-Type", file.MimeType)

	logger.Info("DownloadFile request completed", zap.String("request_id", requestID))
	return c.File(file.FilePath)
}

func (h *FileHandler) ServeFile(c echo.Context) error {
	requestID := c.Response().Header().Get(echo.HeaderXRequestID)
	logger.Info("ServeFile request started", zap.String("request_id", requestID))

	filename := c.Param("filename")
	if filename == "" {
		logger.Error("Missing filename", zap.String("request_id", requestID))
		return response.BadRequest(c, "Filename is required", nil)
	}

	// Construct file path
	filePath := filepath.Join("uploads", filename)

	// Check if file exists
	if _, err := filepath.Abs(filePath); err != nil {
		logger.Error("Invalid file path", zap.Error(err), zap.String("request_id", requestID))
		return response.NotFound(c, "File not found")
	}

	logger.Info("ServeFile request completed", zap.String("request_id", requestID))
	return c.File(filePath)
}
