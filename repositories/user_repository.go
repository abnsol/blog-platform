package repositories

import (
	"errors"
	"strconv"

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
	err := ur.DB.Create(user).Error
	if err != nil {
		return domain.User{}, errors.New(err.Error())
	}
	return *user, nil
}

func (ur *UserRepository) FetchByEmail(email string) (domain.User, error) {
	var user domain.User
	err := ur.DB.Where("email = ?", email).First(&user).Error
	if err != nil {
		return domain.User{}, errors.New(err.Error())
	}
	return user, nil
}

func (ur *UserRepository) FetchByUsername(username string) (domain.User, error) {
	var user domain.User
	err := ur.DB.Where("username = ?", username).First(&user).Error
	if err != nil {
		return domain.User{}, errors.New(err.Error())
	}
	return user, nil
}

func (ur *UserRepository) ActivateAccount(idStr string) error {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return errors.New("invalid id")
	}

	result := ur.DB.Model(&domain.User{}).Where("id = ?", id).Update("status", "active")
	if result.Error != nil || result.RowsAffected == 0 {
		return errors.New(result.Error.Error())
	}

	return nil
}

func (ur *UserRepository) Fetch(idStr string) (domain.User, error) {
	var user domain.User
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return domain.User{}, errors.New("invalid id")
	}

	err = ur.DB.Where("id = ?", id).First(&user).Error
	if err != nil {
		return domain.User{}, errors.New(err.Error())
	}
	return user, nil
}

func (ur *UserRepository) Promote(idStr string) error {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return errors.New("invalid id")
	}

	result := ur.DB.Model(&domain.User{}).Where("id = ?", id).Update("role", "admin")
	if result.Error != nil || result.RowsAffected == 0 {
		return errors.New("failed to promote user")
	}

	return nil
}

func (ur *UserRepository) Demote(idStr string) error {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return errors.New("invalid id")
	}

	result := ur.DB.Model(&domain.User{}).Where("id = ?", id).Update("role", "user")
	if result.Error != nil || result.RowsAffected == 0 {
		return errors.New("failed to demote user")
	}

	return nil
}
