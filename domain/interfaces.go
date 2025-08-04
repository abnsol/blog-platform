package domain

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