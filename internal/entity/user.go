package entity

import (
	"time"
)

type User struct {
	ID                         int        `json:"id"`
	Name                       string     `json:"name"`
	Email                      string     `json:"email"`
	PasswordHash               string     `json:"-"` // Never include in JSON responses
	Role                       string     `json:"role"`
	EmailVerified              bool       `json:"email_verified"`
	EmailVerificationToken     *string    `json:"-"` // Never include in JSON responses
	EmailVerificationExpiresAt *time.Time `json:"-"` // Never include in JSON responses
	PasswordResetToken         *string    `json:"-"` // Never include in JSON responses
	PasswordResetExpiresAt     *time.Time `json:"-"` // Never include in JSON responses
	CreatedAt                  time.Time  `json:"created_at"`
	UpdatedAt                  time.Time  `json:"updated_at"`
}
