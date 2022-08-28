package helpers_test

import (
	helper "login/helpers"
	"testing"

	"github.com/dgrijalva/jwt-go"
)

type SignedDetails struct {
	Name   string
	UserId string `json:"userId" bson:"userId"`
	jwt.StandardClaims
}

func TestValidateToken(t *testing.T) {

	testCases := []struct {
		name        string
		token       string
		expectError bool
	}{
		{
			name:        "validToken",
			token:       "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJOYW1lIjoidW5pdFRlc3RUb2tlbiIsInVzZXJJZCI6IjYxYWYxOGI4OTNmYTRhYzJhMjU2ODI4ZCIsImV4cCI6NTIzODg2MTQ4MH0.x0Ou-k0YOKUv-v2NFKGo5v0p2_ttrNJW_yuTXV5IGf4",
			expectError: false,
		},
		{
			name:        "tokenExpired",
			token:       "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJOYW1lIjoiVGVzdFVzZXIxIiwidXNlcklkIjoiNjE1NDNlNzM5M2ZhNGEyM2U3MWUyYmUxIiwiZXhwIjoxNjMyOTk3MzYzfQ.04Ow4SJ_aLSpb0Qo3lrnnKJcmjd7E5JgaAFQAv6pnR8",
			expectError: true,
		},
		{
			name:        "InvalidSignature",
			token:       "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c",
			expectError: true,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.name, func(t *testing.T) {
			_, err := helper.TokenHelper.ValidateToken(tC.token)
			if tC.expectError && err == "" {
				t.Errorf(tC.name + "error expected")
			}
		})
	}
}
