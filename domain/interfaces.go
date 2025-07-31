package domain

type IPasswordInfrastructure interface {
	HashPassword(password string) (string, error)
	ComparePassword(correctPassword []byte, inputPassword []byte) error
}