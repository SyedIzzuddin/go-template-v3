package handler

import (
	"strconv"

	"go-template/internal/dto"
	"go-template/internal/logger"
	"go-template/internal/service"
	"go-template/pkg/pagination"
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

// CreateUser godoc
// @Summary Create a new user (Admin only)
// @Description Create a new user. Requires admin role.
// @Tags User Management
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreateUserRequest true "User creation data"
// @Success 201 {object} response.Response{data=dto.UserResponse} "User created successfully"
// @Failure 400 {object} response.Response{error=[]dto.ValidationError} "Validation error"
// @Failure 401 {object} response.Response "Unauthorized"
// @Failure 403 {object} response.Response "Forbidden - Admin access required"
// @Failure 409 {object} response.Response "User already exists"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /users [post]
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

// GetAllUsers godoc
// @Summary Get all users with pagination and filtering (Moderator+ only)
// @Description Get a paginated list of all users with optional filtering and search. Requires moderator or admin role.
// @Tags User Management
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Items per page (default: 10, max: 100)"
// @Param sort query string false "Sort field: id, name, email, created_at, role (default: id)"
// @Param order query string false "Sort order: ASC, DESC (default: DESC)"
// @Param search query string false "Search in name and email"
// @Param name query string false "Filter by name (partial match)"
// @Param email query string false "Filter by email (partial match)"
// @Param role query string false "Filter by role (exact match)"
// @Param email_verified query bool false "Filter by email verification status"
// @Param created_after query string false "Filter by creation date (RFC3339 format)"
// @Param created_before query string false "Filter by creation date (RFC3339 format)"
// @Success 200 {object} response.Response{data=[]dto.UserResponse,pagination=pagination.PaginationMeta} "Users retrieved successfully"
// @Failure 401 {object} response.Response "Unauthorized"
// @Failure 403 {object} response.Response "Forbidden - Moderator+ access required"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /users [get]
func (h *UserHandler) GetAllUsers(c echo.Context) error {
	requestID := c.Response().Header().Get(echo.HeaderXRequestID)
	logger.Info("GetAllUsers request started", zap.String("request_id", requestID))
	
	// Parse pagination and filter parameters
	paginationParams := pagination.GetPaginationParams(c)
	filterParams := pagination.GetFilterParams(c)
	
	users, paginationMeta, err := h.userService.GetAllUsersWithPagination(c.Request().Context(), paginationParams, filterParams)
	if err != nil {
		logger.Error("Failed to get all users", zap.Error(err), zap.String("request_id", requestID))
		return response.InternalServerError(c, "Failed to get users", err.Error())
	}
	
	logger.Info("GetAllUsers request completed", 
		zap.String("request_id", requestID),
		zap.Int("total_users", paginationMeta.TotalRecords),
		zap.Int("page", paginationMeta.CurrentPage))
	return response.SuccessWithPagination(c, "Users retrieved successfully", users, paginationMeta)
}