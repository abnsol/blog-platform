package domain

import (
	"context"
)

type IJWTInfrastructure interface {
	GenerateAccessToken(userID string, userRole string) (string, error)
	GenerateRefreshToken(userID string, userRole string) (string, error)
	ValidateAccessToken(authHeader string) (*TokenClaims, error)
	ValidateRefreshToken(token string) (*TokenClaims, error)
}

type ITokenRepository interface {
	FetchByContent(content string) (Token, error)
}

type IPasswordInfrastructure interface {
	HashPassword(password string) (string, error)
	ComparePassword(correctPassword []byte, inputPassword []byte) error
}

type IUserUsecase interface {
	Register(user *User) (User, error)
}

type IUserRepository interface {
	Register(user *User) (User, error)
	FetchByUsername(username string) (User, error)
	FetchByEmail(email string) (User, error)
}

type IUserController interface {
	Register(ctx *context.Context)
}

type IEmailInfrastructure interface {
	SendEmail(from string, to []string, content string) error
}