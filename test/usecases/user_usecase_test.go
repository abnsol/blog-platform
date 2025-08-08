package test

import (
	"errors"
	"os"
	"testing"

	"github.com/blog-platform/domain"
	"github.com/blog-platform/test/mocks"
	"github.com/blog-platform/usecases"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type UserUsecaseTestSuite struct {
	suite.Suite
	userRepo     *mocks.MockUserRepository
	emailService *mocks.MockEmailService
	pwdService   *mocks.MockPasswordService
	jwtService   *mocks.MockJWTService
	tokenRepo    *mocks.MockTokenRepository
	userUsecase  domain.IUserUsecase
}

func (suite *UserUsecaseTestSuite) SetupTest() {
	suite.userRepo = new(mocks.MockUserRepository)
	suite.emailService = new(mocks.MockEmailService)
	suite.pwdService = new(mocks.MockPasswordService)
	suite.jwtService = new(mocks.MockJWTService)
	suite.tokenRepo = new(mocks.MockTokenRepository)
	suite.userUsecase = usecases.NewUserUsecase(suite.userRepo, suite.emailService, suite.pwdService, suite.jwtService, suite.tokenRepo)
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
	suite.pwdService.On("HashPassword", user.Password).Return("hashedpassword", nil)
	suite.userRepo.On("Register", mock.AnythingOfType("*domain.User")).Return(createdUser, nil)
	suite.emailService.On("SendEmail", []string{user.Email}, "Activate Account", "http://localhost:8080/user/1/activate").Return(nil)

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

func (suite *UserUsecaseTestSuite) TestRegister_UsernameInUse() {
	user := &domain.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "Password123!",
	}
	suite.userRepo.On("FetchByUsername", user.Username).Return(domain.User{}, nil)
	_, err := suite.userUsecase.Register(user)
	suite.Error(err)
	suite.Equal("this username is already in use", err.Error())
}

func (suite *UserUsecaseTestSuite) TestRegister_EmailInUse() {
	user := &domain.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "Password123!",
	}
	suite.userRepo.On("FetchByUsername", user.Username).Return(domain.User{}, errors.New("not found"))
	suite.userRepo.On("FetchByEmail", user.Email).Return(domain.User{}, nil)
	_, err := suite.userUsecase.Register(user)
	suite.Error(err)
	suite.Equal("this email is already in use", err.Error())
}

func (suite *UserUsecaseTestSuite) TestRegister_RegistrationFails_HashError() {
	user := &domain.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "Password123!",
	}
	suite.userRepo.On("FetchByUsername", user.Username).Return(domain.User{}, errors.New("not found"))
	suite.userRepo.On("FetchByEmail", user.Email).Return(domain.User{}, errors.New("not found"))
	suite.pwdService.On("HashPassword", user.Password).Return("", errors.New("hash error"))

	_, err := suite.userUsecase.Register(user)
	suite.Error(err)
}

func (suite *UserUsecaseTestSuite) TestRegister_RegistrationFails_RepoError() {
	user := &domain.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "Password123!",
	}
	suite.userRepo.On("FetchByUsername", user.Username).Return(domain.User{}, errors.New("not found"))
	suite.userRepo.On("FetchByEmail", user.Email).Return(domain.User{}, errors.New("not found"))
	suite.pwdService.On("HashPassword", user.Password).Return("hashedpassword", nil)
	suite.userRepo.On("Register", mock.AnythingOfType("*domain.User")).Return(domain.User{}, errors.New("db error"))

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
	suite.pwdService.On("HashPassword", user.Password).Return("hashedpassword", nil)
	suite.userRepo.On("Register", mock.AnythingOfType("*domain.User")).Return(createdUser, nil)
	suite.emailService.On("SendEmail", []string{user.Email}, "Activate Account", "http://localhost:8080/user/1/activate").Return(errors.New("email error"))

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

func (suite *UserUsecaseTestSuite) TestLogin_Success() {
	user := &domain.User{
		ID:       1,
		Username: "testuser",
		Password: "hashedpassword",
		Role:     "user",
	}
	suite.userRepo.On("FetchByUsername", "testuser").Return(*user, nil)
	suite.pwdService.On("ComparePassword", []byte(user.Password), []byte("Password123!")).Return(nil)
	suite.jwtService.On("GenerateAccessToken", "1", "user").Return("access_token", nil)
	suite.jwtService.On("GenerateRefreshToken", "1", "user").Return("refresh_token", nil)
	suite.tokenRepo.On("Save", mock.AnythingOfType("*domain.Token")).Return(nil).Twice()

	accessToken, refreshToken, err := suite.userUsecase.Login("testuser", "Password123!")

	suite.NoError(err)
	suite.Equal("access_token", accessToken)
	suite.Equal("refresh_token", refreshToken)
	suite.tokenRepo.AssertNumberOfCalls(suite.T(), "Save", 2)
}

func (suite *UserUsecaseTestSuite) TestLogin_InvalidIdentifier() {
	suite.userRepo.On("FetchByUsername", "unknown").Return(domain.User{}, errors.New("not found"))
	_, _, err := suite.userUsecase.Login("unknown", "Password123!")
	suite.Error(err)
}

func (suite *UserUsecaseTestSuite) TestLogin_InvalidPassword() {
	user := &domain.User{
		ID:       1,
		Username: "testuser",
		Password: "hashedpassword",
	}
	suite.userRepo.On("FetchByUsername", "testuser").Return(*user, nil)
	suite.pwdService.On("ComparePassword", []byte(user.Password), []byte("WrongPassword!")).Return(errors.New("wrong password"))
	_, _, err := suite.userUsecase.Login("testuser", "WrongPassword!")
	suite.Error(err)
}

func (suite *UserUsecaseTestSuite) TestLogin_GenerateAccessTokenError() {
	user := &domain.User{
		ID:       1,
		Username: "testuser",
		Password: "hashedpassword",
		Role:     "user",
	}
	suite.userRepo.On("FetchByUsername", "testuser").Return(*user, nil)
	suite.pwdService.On("ComparePassword", []byte(user.Password), []byte("Password123!")).Return(nil)
	suite.jwtService.On("GenerateAccessToken", "1", "user").Return("", errors.New("jwt error"))

	_, _, err := suite.userUsecase.Login("testuser", "Password123!")
	suite.Error(err)
}

func (suite *UserUsecaseTestSuite) TestLogin_GenerateRefreshTokenError() {
	user := &domain.User{
		ID:       1,
		Username: "testuser",
		Password: "hashedpassword",
		Role:     "user",
	}
	suite.userRepo.On("FetchByUsername", "testuser").Return(*user, nil)
	suite.pwdService.On("ComparePassword", []byte(user.Password), []byte("Password123!")).Return(nil)
	suite.jwtService.On("GenerateAccessToken", "1", "user").Return("access_token", nil)
	suite.jwtService.On("GenerateRefreshToken", "1", "user").Return("", errors.New("jwt error"))

	_, _, err := suite.userUsecase.Login("testuser", "Password123!")
	suite.Error(err)
}

func (suite *UserUsecaseTestSuite) TestLogin_SaveAccessTokenError() {
	user := &domain.User{
		ID:       1,
		Username: "testuser",
		Password: "hashedpassword",
		Role:     "user",
	}
	suite.userRepo.On("FetchByUsername", "testuser").Return(*user, nil)
	suite.pwdService.On("ComparePassword", []byte(user.Password), []byte("Password123!")).Return(nil)
	suite.jwtService.On("GenerateAccessToken", "1", "user").Return("access_token", nil)
	suite.jwtService.On("GenerateRefreshToken", "1", "user").Return("refresh_token", nil)
	suite.tokenRepo.On("Save", mock.AnythingOfType("*domain.Token")).Return(errors.New("db error")).Once()

	_, _, err := suite.userUsecase.Login("testuser", "Password123!")
	suite.Error(err)
}

func (suite *UserUsecaseTestSuite) TestLogin_SaveRefreshTokenError() {
	user := &domain.User{
		ID:       1,
		Username: "testuser",
		Password: "hashedpassword",
		Role:     "user",
	}
	suite.userRepo.On("FetchByUsername", "testuser").Return(*user, nil)
	suite.pwdService.On("ComparePassword", []byte(user.Password), []byte("Password123!")).Return(nil)
	suite.jwtService.On("GenerateAccessToken", "1", "user").Return("access_token", nil)
	suite.jwtService.On("GenerateRefreshToken", "1", "user").Return("refresh_token", nil)
	suite.tokenRepo.On("Save", mock.AnythingOfType("*domain.Token")).Return(nil).Once()
	suite.tokenRepo.On("Save", mock.AnythingOfType("*domain.Token")).Return(errors.New("db error")).Once()

	_, _, err := suite.userUsecase.Login("testuser", "Password123!")
	suite.Error(err)
}

func (suite *UserUsecaseTestSuite) TestPromote_Success() {
	suite.userRepo.On("Fetch", "1").Return(domain.User{}, nil)
	suite.userRepo.On("Promote", "1").Return(nil)
	err := suite.userUsecase.Promote("1")
	suite.NoError(err)
}

func (suite *UserUsecaseTestSuite) TestPromote_UserNotFound() {
	suite.userRepo.On("Fetch", "1").Return(domain.User{}, errors.New("not found"))
	err := suite.userUsecase.Promote("1")
	suite.Error(err)
	suite.Equal("user not found", err.Error())
}

func (suite *UserUsecaseTestSuite) TestDemote_Success() {
	suite.userRepo.On("Fetch", "1").Return(domain.User{}, nil)
	suite.userRepo.On("Demote", "1").Return(nil)
	err := suite.userUsecase.Demote("1")
	suite.NoError(err)
}

func (suite *UserUsecaseTestSuite) TestDemote_UserNotFound() {
	suite.userRepo.On("Fetch", "1").Return(domain.User{}, errors.New("not found"))
	err := suite.userUsecase.Demote("1")
	suite.Error(err)
	suite.Equal("user not found", err.Error())
}

func TestUserUsecase(t *testing.T) {
	suite.Run(t, new(UserUsecaseTestSuite))
}
