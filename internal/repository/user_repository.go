package repository

import (
	"context"
	"database/sql"
	"time"

	db "go-template/db/sqlc"

	"go-template/internal/entity"
)

type UserRepository interface {
	Create(ctx context.Context, name, email string) (*entity.User, error)
	CreateWithPassword(ctx context.Context, name, email, passwordHash string) (*entity.User, error)
	CreateWithPasswordAndRole(ctx context.Context, name, email, passwordHash, role, emailVerificationToken string, emailVerificationExpiresAt *time.Time) (*entity.User, error)
	GetByID(ctx context.Context, id int) (*entity.User, error)
	GetByEmail(ctx context.Context, email string) (*entity.User, error)
	GetByEmailWithPassword(ctx context.Context, email string) (*entity.User, error)
	GetByVerificationToken(ctx context.Context, token string) (*entity.User, error)
	GetByPasswordResetToken(ctx context.Context, token string) (*entity.User, error)
	Update(ctx context.Context, id int, name string) (*entity.User, error)
	VerifyEmail(ctx context.Context, token string) error
	UpdateVerificationToken(ctx context.Context, userID int, token string, expiresAt *time.Time) error
	UpdatePasswordResetToken(ctx context.Context, userID int, token string, expiresAt *time.Time) error
	ResetPassword(ctx context.Context, token, passwordHash string) error
	Delete(ctx context.Context, id int) error
	GetAll(ctx context.Context) ([]entity.User, error)
}

type userRepository struct {
	db      *sql.DB
	queries *db.Queries
}

func NewUserRepository(dbConn *sql.DB) UserRepository {
	return &userRepository{
		db:      dbConn,
		queries: db.New(dbConn),
	}
}

func (r *userRepository) Create(ctx context.Context, name, email string) (*entity.User, error) {
	createdUser, err := r.queries.CreateUser(ctx, db.CreateUserParams{
		Name:  name,
		Email: email,
	})
	if err != nil {
		return nil, err
	}

	return &entity.User{
		ID:                         int(createdUser.ID),
		Name:                       createdUser.Name,
		Email:                      createdUser.Email,
		PasswordHash:               createdUser.PasswordHash,
		Role:                       createdUser.Role,
		EmailVerified:              createdUser.EmailVerified,
		EmailVerificationToken:     nullStringToPtr(createdUser.EmailVerificationToken),
		EmailVerificationExpiresAt: nullTimeToPtr(createdUser.EmailVerificationExpiresAt),
		PasswordResetToken:         nullStringToPtr(createdUser.PasswordResetToken),
		PasswordResetExpiresAt:     nullTimeToPtr(createdUser.PasswordResetExpiresAt),
		CreatedAt:                  createdUser.CreatedAt.Time,
		UpdatedAt:                  createdUser.UpdatedAt.Time,
	}, nil
}

func (r *userRepository) GetByID(ctx context.Context, id int) (*entity.User, error) {
	user, err := r.queries.GetUser(ctx, int32(id))
	if err != nil {
		return nil, err
	}

	return &entity.User{
		ID:                         int(user.ID),
		Name:                       user.Name,
		Email:                      user.Email,
		PasswordHash:               user.PasswordHash,
		Role:                       user.Role,
		EmailVerified:              user.EmailVerified,
		EmailVerificationToken:     nullStringToPtr(user.EmailVerificationToken),
		EmailVerificationExpiresAt: nullTimeToPtr(user.EmailVerificationExpiresAt),
		PasswordResetToken:         nullStringToPtr(user.PasswordResetToken),
		PasswordResetExpiresAt:     nullTimeToPtr(user.PasswordResetExpiresAt),
		CreatedAt:                  user.CreatedAt.Time,
		UpdatedAt:                  user.UpdatedAt.Time,
	}, nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	user, err := r.queries.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	return &entity.User{
		ID:                         int(user.ID),
		Name:                       user.Name,
		Email:                      user.Email,
		PasswordHash:               user.PasswordHash,
		Role:                       user.Role,
		EmailVerified:              user.EmailVerified,
		EmailVerificationToken:     nullStringToPtr(user.EmailVerificationToken),
		EmailVerificationExpiresAt: nullTimeToPtr(user.EmailVerificationExpiresAt),
		PasswordResetToken:         nullStringToPtr(user.PasswordResetToken),
		PasswordResetExpiresAt:     nullTimeToPtr(user.PasswordResetExpiresAt),
		CreatedAt:                  user.CreatedAt.Time,
		UpdatedAt:                  user.UpdatedAt.Time,
	}, nil
}

func (r *userRepository) Update(ctx context.Context, id int, name string) (*entity.User, error) {
	updatedUser, err := r.queries.UpdateUser(ctx, db.UpdateUserParams{
		ID:   int32(id),
		Name: name,
	})
	if err != nil {
		return nil, err
	}

	return &entity.User{
		ID:                         int(updatedUser.ID),
		Name:                       updatedUser.Name,
		Email:                      updatedUser.Email,
		PasswordHash:               updatedUser.PasswordHash,
		Role:                       updatedUser.Role,
		EmailVerified:              updatedUser.EmailVerified,
		EmailVerificationToken:     nullStringToPtr(updatedUser.EmailVerificationToken),
		EmailVerificationExpiresAt: nullTimeToPtr(updatedUser.EmailVerificationExpiresAt),
		PasswordResetToken:         nullStringToPtr(updatedUser.PasswordResetToken),
		PasswordResetExpiresAt:     nullTimeToPtr(updatedUser.PasswordResetExpiresAt),
		CreatedAt:                  updatedUser.CreatedAt.Time,
		UpdatedAt:                  updatedUser.UpdatedAt.Time,
	}, nil
}

func (r *userRepository) Delete(ctx context.Context, id int) error {
	return r.queries.DeleteUser(ctx, int32(id))
}

func (r *userRepository) CreateWithPassword(ctx context.Context, name, email, passwordHash string) (*entity.User, error) {
	createdUser, err := r.queries.CreateUserWithPassword(ctx, db.CreateUserWithPasswordParams{
		Name:                       name,
		Email:                      email,
		PasswordHash:               passwordHash,
		Role:                       "user", // default role
		EmailVerificationToken:     sql.NullString{Valid: false},
		EmailVerificationExpiresAt: sql.NullTime{Valid: false},
	})
	if err != nil {
		return nil, err
	}

	return &entity.User{
		ID:                         int(createdUser.ID),
		Name:                       createdUser.Name,
		Email:                      createdUser.Email,
		PasswordHash:               createdUser.PasswordHash,
		Role:                       createdUser.Role,
		EmailVerified:              createdUser.EmailVerified,
		EmailVerificationToken:     nullStringToPtr(createdUser.EmailVerificationToken),
		EmailVerificationExpiresAt: nullTimeToPtr(createdUser.EmailVerificationExpiresAt),
		PasswordResetToken:         nullStringToPtr(createdUser.PasswordResetToken),
		PasswordResetExpiresAt:     nullTimeToPtr(createdUser.PasswordResetExpiresAt),
		CreatedAt:                  createdUser.CreatedAt.Time,
		UpdatedAt:                  createdUser.UpdatedAt.Time,
	}, nil
}

func (r *userRepository) GetByEmailWithPassword(ctx context.Context, email string) (*entity.User, error) {
	user, err := r.queries.GetUserByEmailWithPassword(ctx, email)
	if err != nil {
		return nil, err
	}

	return &entity.User{
		ID:                         int(user.ID),
		Name:                       user.Name,
		Email:                      user.Email,
		PasswordHash:               user.PasswordHash,
		Role:                       user.Role,
		EmailVerified:              user.EmailVerified,
		EmailVerificationToken:     nullStringToPtr(user.EmailVerificationToken),
		EmailVerificationExpiresAt: nullTimeToPtr(user.EmailVerificationExpiresAt),
		PasswordResetToken:         nullStringToPtr(user.PasswordResetToken),
		PasswordResetExpiresAt:     nullTimeToPtr(user.PasswordResetExpiresAt),
		CreatedAt:                  user.CreatedAt.Time,
		UpdatedAt:                  user.UpdatedAt.Time,
	}, nil
}

func (r *userRepository) GetAll(ctx context.Context) ([]entity.User, error) {
	// Get all users without pagination
	userList, err := r.queries.GetAllUsers(ctx)
	if err != nil {
		return nil, err
	}

	users := make([]entity.User, len(userList))
	for i, dbUser := range userList {
		users[i] = entity.User{
			ID:                         int(dbUser.ID),
			Name:                       dbUser.Name,
			Email:                      dbUser.Email,
			PasswordHash:               dbUser.PasswordHash,
			Role:                       dbUser.Role,
			EmailVerified:              dbUser.EmailVerified,
			EmailVerificationToken:     nullStringToPtr(dbUser.EmailVerificationToken),
			EmailVerificationExpiresAt: nullTimeToPtr(dbUser.EmailVerificationExpiresAt),
			PasswordResetToken:         nullStringToPtr(dbUser.PasswordResetToken),
			PasswordResetExpiresAt:     nullTimeToPtr(dbUser.PasswordResetExpiresAt),
			CreatedAt:                  dbUser.CreatedAt.Time,
			UpdatedAt:                  dbUser.UpdatedAt.Time,
		}
	}

	return users, nil
}

func (r *userRepository) CreateWithPasswordAndRole(ctx context.Context, name, email, passwordHash, role, emailVerificationToken string, emailVerificationExpiresAt *time.Time) (*entity.User, error) {
	createdUser, err := r.queries.CreateUserWithPassword(ctx, db.CreateUserWithPasswordParams{
		Name:                       name,
		Email:                      email,
		PasswordHash:               passwordHash,
		Role:                       role,
		EmailVerificationToken:     ptrToNullString(&emailVerificationToken),
		EmailVerificationExpiresAt: ptrToNullTime(emailVerificationExpiresAt),
	})
	if err != nil {
		return nil, err
	}

	return &entity.User{
		ID:                         int(createdUser.ID),
		Name:                       createdUser.Name,
		Email:                      createdUser.Email,
		PasswordHash:               createdUser.PasswordHash,
		Role:                       createdUser.Role,
		EmailVerified:              createdUser.EmailVerified,
		EmailVerificationToken:     nullStringToPtr(createdUser.EmailVerificationToken),
		EmailVerificationExpiresAt: nullTimeToPtr(createdUser.EmailVerificationExpiresAt),
		PasswordResetToken:         nullStringToPtr(createdUser.PasswordResetToken),
		PasswordResetExpiresAt:     nullTimeToPtr(createdUser.PasswordResetExpiresAt),
		CreatedAt:                  createdUser.CreatedAt.Time,
		UpdatedAt:                  createdUser.UpdatedAt.Time,
	}, nil
}

func (r *userRepository) GetByVerificationToken(ctx context.Context, token string) (*entity.User, error) {
	user, err := r.queries.GetUserByVerificationToken(ctx, sql.NullString{String: token, Valid: true})
	if err != nil {
		return nil, err
	}

	return &entity.User{
		ID:                         int(user.ID),
		Name:                       user.Name,
		Email:                      user.Email,
		PasswordHash:               user.PasswordHash,
		Role:                       user.Role,
		EmailVerified:              user.EmailVerified,
		EmailVerificationToken:     nullStringToPtr(user.EmailVerificationToken),
		EmailVerificationExpiresAt: nullTimeToPtr(user.EmailVerificationExpiresAt),
		PasswordResetToken:         nullStringToPtr(user.PasswordResetToken),
		PasswordResetExpiresAt:     nullTimeToPtr(user.PasswordResetExpiresAt),
		CreatedAt:                  user.CreatedAt.Time,
		UpdatedAt:                  user.UpdatedAt.Time,
	}, nil
}

func (r *userRepository) VerifyEmail(ctx context.Context, token string) error {
	return r.queries.VerifyEmailByToken(ctx, sql.NullString{String: token, Valid: true})
}

func (r *userRepository) UpdateVerificationToken(ctx context.Context, userID int, token string, expiresAt *time.Time) error {
	return r.queries.UpdateVerificationToken(ctx, db.UpdateVerificationTokenParams{
		ID:                         int32(userID),
		EmailVerificationToken:     ptrToNullString(&token),
		EmailVerificationExpiresAt: ptrToNullTime(expiresAt),
	})
}

// Helper functions to convert between sql.Null* and pointers
func nullStringToPtr(ns sql.NullString) *string {
	if ns.Valid {
		return &ns.String
	}
	return nil
}

func nullTimeToPtr(nt sql.NullTime) *time.Time {
	if nt.Valid {
		return &nt.Time
	}
	return nil
}

func ptrToNullString(s *string) sql.NullString {
	if s != nil {
		return sql.NullString{String: *s, Valid: true}
	}
	return sql.NullString{Valid: false}
}

func ptrToNullTime(t *time.Time) sql.NullTime {
	if t != nil {
		return sql.NullTime{Time: *t, Valid: true}
	}
	return sql.NullTime{Valid: false}
}

func (r *userRepository) GetByPasswordResetToken(ctx context.Context, token string) (*entity.User, error) {
	user, err := r.queries.GetUserByPasswordResetToken(ctx, sql.NullString{String: token, Valid: true})
	if err != nil {
		return nil, err
	}

	return &entity.User{
		ID:                         int(user.ID),
		Name:                       user.Name,
		Email:                      user.Email,
		PasswordHash:               user.PasswordHash,
		Role:                       user.Role,
		EmailVerified:              user.EmailVerified,
		EmailVerificationToken:     nullStringToPtr(user.EmailVerificationToken),
		EmailVerificationExpiresAt: nullTimeToPtr(user.EmailVerificationExpiresAt),
		PasswordResetToken:         nullStringToPtr(user.PasswordResetToken),
		PasswordResetExpiresAt:     nullTimeToPtr(user.PasswordResetExpiresAt),
		CreatedAt:                  user.CreatedAt.Time,
		UpdatedAt:                  user.UpdatedAt.Time,
	}, nil
}

func (r *userRepository) UpdatePasswordResetToken(ctx context.Context, userID int, token string, expiresAt *time.Time) error {
	return r.queries.UpdatePasswordResetToken(ctx, db.UpdatePasswordResetTokenParams{
		ID:                     int32(userID),
		PasswordResetToken:     ptrToNullString(&token),
		PasswordResetExpiresAt: ptrToNullTime(expiresAt),
	})
}

func (r *userRepository) ResetPassword(ctx context.Context, token, passwordHash string) error {
	return r.queries.ResetPassword(ctx, db.ResetPasswordParams{
		PasswordResetToken: sql.NullString{String: token, Valid: true},
		PasswordHash:       passwordHash,
	})
}
