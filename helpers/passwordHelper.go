package helpers

import (
	"math/rand"

	"golang.org/x/crypto/bcrypt"
)

var genHashPasswd func(password []byte, cost int) ([]byte, error) = bcrypt.GenerateFromPassword
var compHashAndPass func(hashedPassword []byte, password []byte) error = bcrypt.CompareHashAndPassword

type passwordHelper interface {
	HashPassword(password string) (string, error)
	VerifyPassword(userPassword string, providedPassword string) (bool, string)
	GenerateRandomPassword() string
}

type passwordHelperImpl struct {
}

var (
	PasswordHelper passwordHelper = &passwordHelperImpl{}
)

func (helper *passwordHelperImpl) HashPassword(password string) (string, error) {
	bytes, err := genHashPasswd([]byte(password), 14)
	if err != nil {
		return "", err
	}
	return string(bytes), err
}

func (helper *passwordHelperImpl) VerifyPassword(userPassword string, providedPassword string) (bool, string) {
	err := compHashAndPass([]byte(providedPassword), []byte(userPassword))
	check := true
	msg := ""
	if err != nil {
		msg = "email or password is incorrect"
		check = false
	}
	return check, msg
}
func (helper *passwordHelperImpl) GenerateRandomPassword() string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*")
	n := 12
	s := make([]rune, n)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}
	return string(s)
}
