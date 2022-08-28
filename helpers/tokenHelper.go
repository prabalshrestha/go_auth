package helpers

import (
	"log"
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

type tokenHelper interface {
	GenerateAllToken(name string, userId string, userRole []string) (signedToken string, signedRefreshToken string, err error)
	ValidateToken(siginedToken string) (claims *SignedDetails, msg string)
}

type tokenHelperImpl struct {
}

var (
	TokenHelper tokenHelper = &tokenHelperImpl{}
)

type SignedDetails struct {
	Name     string   `json:"name" bson:"name"`
	UserId   string   `json:"userId" bson:"userId"`
	UserRole []string `json:"userRole" bson:"userRole"`
	jwt.StandardClaims
}

var SECRET_KEY string = os.Getenv("SECRET_KEY")

// var SECRET_KEY string = "secret"

func (helper *tokenHelperImpl) GenerateAllToken(name string, userId string, userRole []string) (signedToken string, signedRefreshToken string, err error) {
	claims := &SignedDetails{
		Name:     name,
		UserId:   userId,
		UserRole: userRole,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(10)).Unix(),
		},
	}
	refreshClaims := &SignedDetails{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Duration(1000)).Unix(),
		},
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(SECRET_KEY))
	if err != nil {
		log.Panic(err)
		return
	}
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(SECRET_KEY))
	if err != nil {
		log.Panic(err)
		return
	}
	return token, refreshToken, err
}

var (
	jwtParseWithClaims = jwt.ParseWithClaims
)

func (helper *tokenHelperImpl) ValidateToken(siginedToken string) (claims *SignedDetails, msg string) {
	token, err := jwtParseWithClaims(
		siginedToken,
		&SignedDetails{},
		func(t *jwt.Token) (interface{}, error) {
			return []byte(SECRET_KEY), nil
		},
	)
	if err != nil {
		msg = err.Error()
	}
	claims, ok := token.Claims.(*SignedDetails)
	if !ok {
		msg = "Token is invalid"
		return
	}
	if claims.ExpiresAt < time.Now().Local().Unix() {
		msg = "Token is expired"
		return
	}
	return claims, msg
}
