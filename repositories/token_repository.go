package repositories

import (
	"github.com/blog-platform/domain"
	"gorm.io/gorm"
)

type TokenRepository struct {
	DB *gorm.DB
}

func NewTokenRepository(db *gorm.DB) *TokenRepository {
	return &TokenRepository{
		DB: db,
	}
}

func (repo *TokenRepository) FetchByContent(content string) (domain.Token, error) {
	var token domain.Token
	result := repo.DB.First(&token, "content = ?", content)
	if result.Error != nil {
		return domain.Token{}, result.Error
	}

	return token, nil
}