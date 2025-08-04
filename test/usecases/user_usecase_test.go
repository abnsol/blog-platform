package test

import (
	"errors"
	"os"
	"testing"

	"github.com/blog-platform/domain"
	"github.com/blog-platform/test/mocks"
	"github.com/blog-platform/usecases"
	"github.com/stretchr/testify/suite"
)

type UserUsecaseTestSuite struct {
	suite.Suite
	userRepo     *mocks.MockUserRepository
	emailService *mocks.MockEmailService
	userUsecase  domain.IUserUsecase
}

func (suite *UserUsecaseTestSuite) SetupTest() {
	suite.userRepo = new(mocks.MockUserRepository)
	suite.emailService = new(mocks.MockEmailService)
	suite.userUsecase = usecases.NewUserUsecase(suite.userRepo, suite.emailService)
	os.Setenv("PROTOCOL", "http")
	os.Setenv("DOMAIN", "localhost")
	os.Setenv("PORT", "8080")
	os.Setenv("EMAIL_SENDER", "test@example.com")
}

func (suite *UserUsecaseTestSuite) TestRegister_Success() {
	user := &domain.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "Password123!",
	}
	createdUser := *user
	createdUser.ID = 1

	suite.userRepo.On("FetchByUsername", user.Username).Return(domain.User{}, errors.New("not found"))
	suite.userRepo.On("FetchByEmail", user.Email).Return(domain.User{}, errors.New("not found"))
	suite.userRepo.On("Register", user).Return(createdUser, nil)
	suite.emailService.On("SendEmail", "test@example.com", []string{user.Email}, "http://localhost:8080/user/1/activate").Return(nil)

	_, err := suite.userUsecase.Register(user)
	suite.NoError(err)
}

func (suite *UserUsecaseTestSuite) TestRegister_MissingFields() {
	user := &domain.User{}
	_, err := suite.userUsecase.Register(user)
	suite.Error(err)
	suite.Equal("missing required fields", err.Error())
}

func (suite *UserUsecaseTestSuite) TestRegister_InvalidEmail() {
	user := &domain.User{
		Username: "testuser",
		Email:    "invalid-email",
		Password: "Password123!",
	}
	_, err := suite.userUsecase.Register(user)
	suite.Error(err)
	suite.Equal("invalid email format", err.Error())
}

func (suite *UserUsecaseTestSuite) TestRegister_InvalidPassword() {
	user := &domain.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "weak",
	}
	_, err := suite.userUsecase.Register(user)
	suite.Error(err)
	suite.Equal("password must be consisted of at least one uppercase character, one lowercase character, one punctuation character, one number and be at least of length 8", err.Error())
}

func (suite *UserUsecaseTestSuite) TestRegister_EmailInUse() {
	user := &domain.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "Password123!",
	}
	suite.userRepo.On("FetchByUsername", user.Username).Return(domain.User{}, nil)
	_, err := suite.userUsecase.Register(user)
	suite.Error(err)
	suite.Equal("this email is already in use", err.Error())
}

func (suite *UserUsecaseTestSuite) TestRegister_UsernameInUse() {
	user := &domain.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "Password123!",
	}
	suite.userRepo.On("FetchByUsername", user.Username).Return(domain.User{}, errors.New("not found"))
	suite.userRepo.On("FetchByEmail", user.Email).Return(domain.User{}, nil)
	_, err := suite.userUsecase.Register(user)
	suite.Error(err)
	suite.Equal("this username is already in use", err.Error())
}

func (suite *UserUsecaseTestSuite) TestRegister_RegistrationFails() {
	user := &domain.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "Password123!",
	}
	suite.userRepo.On("FetchByUsername", user.Username).Return(domain.User{}, errors.New("not found"))
	suite.userRepo.On("FetchByEmail", user.Email).Return(domain.User{}, errors.New("not found"))
	suite.userRepo.On("Register", user).Return(domain.User{}, errors.New("db error"))
	_, err := suite.userUsecase.Register(user)
	suite.Error(err)
	suite.Equal("unable to register user", err.Error())
}

func (suite *UserUsecaseTestSuite) TestRegister_SendEmailFails() {
	user := &domain.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "Password123!",
	}
	createdUser := *user
	createdUser.ID = 1

	suite.userRepo.On("FetchByUsername", user.Username).Return(domain.User{}, errors.New("not found"))
	suite.userRepo.On("FetchByEmail", user.Email).Return(domain.User{}, errors.New("not found"))
	suite.userRepo.On("Register", user).Return(createdUser, nil)
	suite.emailService.On("SendEmail", "test@example.com", []string{user.Email}, "http://localhost:8080/user/1/activate").Return(errors.New("email error"))

	_, err := suite.userUsecase.Register(user)
	suite.Error(err)
	suite.Equal("unable to send activation link", err.Error())
}

func (suite *UserUsecaseTestSuite) TestActivateAccount_Success() {
	suite.userRepo.On("Fetch", "1").Return(domain.User{}, nil)
	suite.userRepo.On("ActivateAccount", "1").Return(nil)
	err := suite.userUsecase.ActivateAccount("1")
	suite.NoError(err)
}

func (suite *UserUsecaseTestSuite) TestActivateAccount_UserNotFound() {
	suite.userRepo.On("Fetch", "1").Return(domain.User{}, errors.New("not found"))
	err := suite.userUsecase.ActivateAccount("1")
	suite.Error(err)
}

func (suite *UserUsecaseTestSuite) TestActivateAccount_ActivationFails() {
	suite.userRepo.On("Fetch", "1").Return(domain.User{}, nil)
	suite.userRepo.On("ActivateAccount", "1").Return(errors.New("db error"))
	err := suite.userUsecase.ActivateAccount("1")
	suite.Error(err)
}

func TestUserUsecase(t *testing.T) {
	suite.Run(t, new(UserUsecaseTestSuite))
}
