package domain

import "github.com/blog-platform/infrastructure"

type IJWTInfrastructure interface {
	GenerateAccessToken(userID string, userRole string) (string, error)
	GenerateRefreshToken(userID string, userRole string) (string, error)
	ValidateToken(authHeader string, secret []byte) (*infrastructure.TokenClaims, error)
}