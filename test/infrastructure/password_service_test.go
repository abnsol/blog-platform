package test

import (
	"testing"
	"strings"

	"github.com/stretchr/testify/suite"
	"github.com/blog-platform/infrastructure"
)

type PasswordServiceTestSuite struct {
	suite.Suite
	service *infrastructure.PasswordInfrastructure
}

func (suite *PasswordServiceTestSuite) SetupTest() {
	suite.service = new(infrastructure.PasswordInfrastructure)
}

func (suite *PasswordServiceTestSuite) TestHashAndCompare_Success() {
	pwd := "secret123"
	hash, err := suite.service.HashPassword(pwd)
	suite.NoError(err)
	suite.NotEmpty(hash)

	err = suite.service.ComparePassword([]byte(hash), []byte(pwd))
	suite.NoError(err)
}

func (suite *PasswordServiceTestSuite) TestCompare_WrongPassword() {
	hash, err := suite.service.HashPassword("secret")
	suite.NoError(err)

	err = suite.service.ComparePassword([]byte(hash), []byte("invalid"))
	suite.Error(err)
}

func (suite *PasswordServiceTestSuite) TestCompare_InvalidHash() {
	err := suite.service.ComparePassword([]byte("not-a-valid-hash"), []byte("password"))
	suite.Error(err)
}

func (suite *PasswordServiceTestSuite) TestHash_Uniqueness() {
	h1, err := suite.service.HashPassword("dupassword")
	suite.NoError(err)
	suite.NotEmpty(h1)

	h2, err := suite.service.HashPassword("dupassword")
	suite.NoError(err)
	suite.NotEmpty(h2)

	suite.NotEqual(h1, h2)
}

func (suite *PasswordServiceTestSuite) TestHash_EdgeCases() {
	// empty password
	hEmpty, err := suite.service.HashPassword("")
	suite.NoError(err)
	suite.NotEmpty(hEmpty)

	// very long password
	longPwd := strings.Repeat("x", 72)
	hLong, err := suite.service.HashPassword(longPwd)
	suite.NoError(err)
	suite.NotEmpty(hLong)
}

func TestPasswordServiceTestSuite(t *testing.T) {
	suite.Run(t, new(PasswordServiceTestSuite))
}