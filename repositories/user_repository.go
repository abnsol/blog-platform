package repositories

import (
	"errors"

	"github.com/blog-platform/domain"
	"gorm.io/gorm"
)

type UserRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		DB: db,
	}
}

func (ur *UserRepository) Register(user *domain.User) (domain.User, error) {
	err := ur.DB.Create(user)
	if err != nil {
		return domain.User{}, errors.New("unable to register user")
	}
	return *user, nil
}

func (ur *UserRepository) FetchByEmail(email string) (domain.User, error) {
	var user domain.User
	err := ur.DB.Where("email = ?", email).First(&user).Error
	if err != nil {
		return domain.User{}, errors.New("user not found")
	}
	return user, nil
}

func (ur *UserRepository) FetchByUsername(username string) (domain.User, error) {
	var user domain.User
	err := ur.DB.Where("username = ?", username).First(&user).Error
	if err != nil {
		return domain.User{}, errors.New("user not found")
	}
	return user, nil
}