package service

import (
	"context"
	"database/sql"
	"errors"
	"go-template/internal/dto"
	"go-template/internal/entity"
	"go-template/internal/logger"
	"go-template/internal/repository"

	"go.uber.org/zap"
)

type UserService interface {
	CreateUser(ctx context.Context, req dto.CreateUserRequest) (*dto.UserResponse, error)
	GetUserByID(ctx context.Context, id int) (*dto.UserResponse, error)
	UpdateUser(ctx context.Context, id int, req dto.UpdateUserRequest) (*dto.UserResponse, error)
	DeleteUser(ctx context.Context, id int) error
	GetAllUsers(ctx context.Context) ([]dto.UserResponse, error)
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

func (s *userService) CreateUser(ctx context.Context, req dto.CreateUserRequest) (*dto.UserResponse, error) {
	logger.Info("Creating new user", zap.String("email", req.Email))
	
	// Check if user already exists
	existingUser, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err == nil && existingUser != nil {
		logger.Warn("User already exists", zap.String("email", req.Email))
		return nil, errors.New("user with this email already exists")
	}
	
	// Create new user
	user, err := s.userRepo.Create(ctx, req.Name, req.Email)
	if err != nil {
		logger.Error("Failed to create user", zap.Error(err))
		return nil, err
	}
	
	logger.Info("User created successfully", zap.Int("user_id", user.ID))
	
	return s.mapUserToResponse(user), nil
}

func (s *userService) GetUserByID(ctx context.Context, id int) (*dto.UserResponse, error) {
	logger.Debug("Getting user by ID", zap.Int("user_id", id))
	
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.Warn("User not found", zap.Int("user_id", id))
			return nil, errors.New("user not found")
		}
		logger.Error("Failed to get user", zap.Error(err))
		return nil, err
	}
	
	return s.mapUserToResponse(user), nil
}

func (s *userService) UpdateUser(ctx context.Context, id int, req dto.UpdateUserRequest) (*dto.UserResponse, error) {
	logger.Info("Updating user", zap.Int("user_id", id))
	
	// Check if user exists
	_, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.Warn("User not found for update", zap.Int("user_id", id))
			return nil, errors.New("user not found")
		}
		logger.Error("Failed to get user for update", zap.Error(err))
		return nil, err
	}
	
	// Update user
	user, err := s.userRepo.Update(ctx, id, req.Name)
	if err != nil {
		logger.Error("Failed to update user", zap.Error(err))
		return nil, err
	}
	
	logger.Info("User updated successfully", zap.Int("user_id", id))
	
	return s.mapUserToResponse(user), nil
}

func (s *userService) DeleteUser(ctx context.Context, id int) error {
	logger.Info("Deleting user", zap.Int("user_id", id))
	
	// Check if user exists
	_, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.Warn("User not found for deletion", zap.Int("user_id", id))
			return errors.New("user not found")
		}
		logger.Error("Failed to get user for deletion", zap.Error(err))
		return err
	}
	
	if err := s.userRepo.Delete(ctx, id); err != nil {
		logger.Error("Failed to delete user", zap.Error(err))
		return err
	}
	
	logger.Info("User deleted successfully", zap.Int("user_id", id))
	
	return nil
}

func (s *userService) GetAllUsers(ctx context.Context) ([]dto.UserResponse, error) {
	logger.Debug("Getting all users")
	
	users, err := s.userRepo.GetAll(ctx)
	if err != nil {
		logger.Error("Failed to get all users", zap.Error(err))
		return nil, err
	}
	
	var userResponses []dto.UserResponse
	for _, user := range users {
		userResponses = append(userResponses, *s.mapUserToResponse(&user))
	}
	
	return userResponses, nil
}

func (s *userService) mapUserToResponse(user *entity.User) *dto.UserResponse {
	return &dto.UserResponse{
		ID:            user.ID,
		Name:          user.Name,
		Email:         user.Email,
		Role:          user.Role,
		EmailVerified: user.EmailVerified,
		CreatedAt:     user.CreatedAt,
		UpdatedAt:     user.UpdatedAt,
	}
}