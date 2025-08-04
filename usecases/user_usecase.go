package usecases

import (
	"errors"
	"fmt"
	"net/mail"
	"os"
	"unicode"

	"github.com/blog-platform/domain"
)

type UserUsecase struct {
	userRepo domain.IUserRepository
	emailService domain.IEmailInfrastructure
}

func NewUserUsecase(ur domain.IUserRepository, es domain.IEmailInfrastructure) *UserUsecase {
	return &UserUsecase{
		userRepo: ur,
		emailService: es,
	}
}

func (uu *UserUsecase) Register(user *domain.User) (domain.User, error) {
	if user.Username == "" || user.Email == "" || user.Password == "" {
		return domain.User{}, errors.New("missing required fields")
	}

	_, err := mail.ParseAddress(user.Email)
	if err != nil {
		return domain.User{}, errors.New("invalid email format")
	}

	if !uu.validatePassword(user.Password) {
		return domain.User{}, errors.New("password must be consisted of at least one uppercase character, one lowercase character, one punctuation character, one number and be at least of length 8")
	}

	_, err = uu.userRepo.FetchByUsername(user.Username)
	if err == nil {
		return domain.User{}, errors.New("this email is already in use")
	}

	_, err = uu.userRepo.FetchByEmail(user.Email)
	if err == nil {
		return domain.User{}, errors.New("this username is already in use")
	}

	user.Status = "inactive"
	registeredUser, err := uu.userRepo.Register(user)
	if err != nil {
		return domain.User{}, errors.New("unable to register user")
	}

	emailContent := fmt.Sprintf("%v://%v:%v/user/%v/activate", os.Getenv("PROTOCOL"), os.Getenv("DOMAIN"), os.Getenv("PORT"), registeredUser.ID)
	err = uu.emailService.SendEmail(os.Getenv("EMAIL_SENDER"), []string{registeredUser.Email}, emailContent)
	if err != nil {
		return domain.User{}, errors.New("unable to send activation link")
	}

	return registeredUser, nil
}

func (uu *UserUsecase) validatePassword(password string) bool {
	var (
        hasMinLen  = false
        hasUpper   = false
        hasLower   = false
        hasNumber  = false
        hasSpecial = false
    )

    if len(password) >= 8 {
        hasMinLen = true
    }

    for _, c := range password {
        switch {
        case unicode.IsUpper(c):
            hasUpper = true
        case unicode.IsLower(c):
            hasLower = true
        case unicode.IsNumber(c):
            hasNumber = true
        case unicode.IsPunct(c) || unicode.IsSymbol(c):
            hasSpecial = true
        }
    }

    return hasMinLen && hasUpper && hasLower && hasNumber && hasSpecial
}

func (uu *UserUsecase) ActivateAccount(id string) error {
	_, err := uu.userRepo.Fetch(id)
	if err != nil {
		return err
	}

	err = uu.userRepo.ActivateAccount(id)
	if err != nil {
		return err
	}

	return nil
}