package handler

import (
	"strconv"

	"go-template/internal/dto"
	"go-template/internal/logger"
	"go-template/internal/service"
	"go-template/pkg/response"
	"go-template/pkg/validator"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type UserHandler struct {
	userService service.UserService
	validator   *validator.Validator
}

func NewUserHandler(userService service.UserService, validator *validator.Validator) *UserHandler {
	return &UserHandler{
		userService: userService,
		validator:   validator,
	}
}

func (h *UserHandler) CreateUser(c echo.Context) error {
	requestID := c.Response().Header().Get(echo.HeaderXRequestID)
	logger.Info("CreateUser request started", zap.String("request_id", requestID))
	
	var req dto.CreateUserRequest
	if err := c.Bind(&req); err != nil {
		logger.Error("Failed to bind request", zap.Error(err), zap.String("request_id", requestID))
		return response.BadRequest(c, "Invalid request body", err.Error())
	}
	
	// Validate request
	if validationErrors := h.validator.ValidateStruct(req); validationErrors != nil {
		logger.Warn("Validation failed", zap.Any("errors", validationErrors), zap.String("request_id", requestID))
		return response.ValidationError(c, "Validation failed", validationErrors)
	}
	
	user, err := h.userService.CreateUser(c.Request().Context(), req)
	if err != nil {
		logger.Error("Failed to create user", zap.Error(err), zap.String("request_id", requestID))
		return response.BadRequest(c, "Failed to create user", err.Error())
	}
	
	logger.Info("CreateUser request completed", zap.String("request_id", requestID))
	return response.Created(c, "User created successfully", user)
}

func (h *UserHandler) GetUser(c echo.Context) error {
	requestID := c.Response().Header().Get(echo.HeaderXRequestID)
	logger.Info("GetUser request started", zap.String("request_id", requestID))
	
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logger.Error("Invalid user ID", zap.Error(err), zap.String("request_id", requestID))
		return response.BadRequest(c, "Invalid user ID", err.Error())
	}
	
	user, err := h.userService.GetUserByID(c.Request().Context(), id)
	if err != nil {
		logger.Error("Failed to get user", zap.Error(err), zap.String("request_id", requestID))
		return response.NotFound(c, "User not found")
	}
	
	logger.Info("GetUser request completed", zap.String("request_id", requestID))
	return response.Success(c, "User retrieved successfully", user)
}

func (h *UserHandler) UpdateUser(c echo.Context) error {
	requestID := c.Response().Header().Get(echo.HeaderXRequestID)
	logger.Info("UpdateUser request started", zap.String("request_id", requestID))
	
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logger.Error("Invalid user ID", zap.Error(err), zap.String("request_id", requestID))
		return response.BadRequest(c, "Invalid user ID", err.Error())
	}
	
	var req dto.UpdateUserRequest
	if err := c.Bind(&req); err != nil {
		logger.Error("Failed to bind request", zap.Error(err), zap.String("request_id", requestID))
		return response.BadRequest(c, "Invalid request body", err.Error())
	}
	
	// Validate request
	if validationErrors := h.validator.ValidateStruct(req); validationErrors != nil {
		logger.Warn("Validation failed", zap.Any("errors", validationErrors), zap.String("request_id", requestID))
		return response.ValidationError(c, "Validation failed", validationErrors)
	}
	
	user, err := h.userService.UpdateUser(c.Request().Context(), id, req)
	if err != nil {
		logger.Error("Failed to update user", zap.Error(err), zap.String("request_id", requestID))
		return response.BadRequest(c, "Failed to update user", err.Error())
	}
	
	logger.Info("UpdateUser request completed", zap.String("request_id", requestID))
	return response.Success(c, "User updated successfully", user)
}

func (h *UserHandler) DeleteUser(c echo.Context) error {
	requestID := c.Response().Header().Get(echo.HeaderXRequestID)
	logger.Info("DeleteUser request started", zap.String("request_id", requestID))
	
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logger.Error("Invalid user ID", zap.Error(err), zap.String("request_id", requestID))
		return response.BadRequest(c, "Invalid user ID", err.Error())
	}
	
	err = h.userService.DeleteUser(c.Request().Context(), id)
	if err != nil {
		logger.Error("Failed to delete user", zap.Error(err), zap.String("request_id", requestID))
		return response.BadRequest(c, "Failed to delete user", err.Error())
	}
	
	logger.Info("DeleteUser request completed", zap.String("request_id", requestID))
	return response.Success(c, "User deleted successfully", nil)
}

func (h *UserHandler) GetAllUsers(c echo.Context) error {
	requestID := c.Response().Header().Get(echo.HeaderXRequestID)
	logger.Info("GetAllUsers request started", zap.String("request_id", requestID))
	
	users, err := h.userService.GetAllUsers(c.Request().Context())
	if err != nil {
		logger.Error("Failed to get all users", zap.Error(err), zap.String("request_id", requestID))
		return response.InternalServerError(c, "Failed to get users", err.Error())
	}
	
	logger.Info("GetAllUsers request completed", zap.String("request_id", requestID))
	return response.Success(c, "Users retrieved successfully", users)
}