package mocks

import "github.com/stretchr/testify/mock"

type MockEmailService struct {
	mock.Mock
}

func (m *MockEmailService) SendEmail(from string, to []string, content string) error {
	args := m.Called(from, to, content)
	return args.Error(0)
}
