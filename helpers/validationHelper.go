package helpers

import (
	"errors"
	"login/models"
	"login/service"
	"net/mail"
)

type validationHepler interface {
	ValidateNewUser(user models.User) (bool, error)
}

type validationHeplerImpl struct{}

var (
	ValidationHepler validationHepler = &validationHeplerImpl{}
)
var (
	errEmailAlreadyExits    = errors.New("email already exits")
	errInvalidEmail         = errors.New("Invalid Email")
	errUsernameAlreadyExits = errors.New("Username already exits")
	errUserAlreadyExits     = errors.New("Email and Username already exits")
	errInvalidID            = errors.New("invalid id")
	errInvalidBody          = errors.New("invalid request body")
	errInvalidCred          = errors.New("invalid email or password")
	errInsertionFailed      = errors.New("error in the user insertion")
	errUpdationFailed       = errors.New("error in the user updation")
	errDeletionFailed       = errors.New("error in the user deletion")
)
var UserCollection = "user"

func (validationHepler *validationHeplerImpl) ValidateNewUser(user models.User) (bool, error) {

	if user.Email == "" || user.Password == "" || user.Name == "" {
		return false, errInvalidBody
	}
	isValid := ValidateEmail(user.Email)
	if !isValid {
		return false, errInvalidEmail
	}
	emailExists := service.UserService.Exists("email", user.Email, UserCollection)
	nameExists := service.UserService.Exists("name", user.Name, UserCollection)
	if emailExists && nameExists {
		return false, errUserAlreadyExits
	}
	if emailExists {
		return false, errEmailAlreadyExits

	}
	if nameExists {
		return false, errUsernameAlreadyExits
	}

	return true, nil
}

func ValidateEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}
