package helpers

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_VerifyPassword(t *testing.T) {
	type args struct {
		userPassword string
		hashedPasswd string
	}
	type want struct {
		check bool
		msg   string
	}
	subtests := []struct {
		name string
		args args
		want
		compareHashWithPassword func(hashedPassword, password []byte) error
	}{
		{
			name: "Correct Password",
			args: args{
				userPassword: "test",
				hashedPasswd: "$2a$14$mE0dbCsxOlQSNXYgrZqzGu9fgSvOdttgKwExmtl3O0j6LINbTT16u",
			},
			want: want{
				check: true,
				msg:   "",
			},
			compareHashWithPassword: func(hashedPassword, password []byte) error {
				return nil
			},
		},
		{
			name: "Incorrect Password",
			args: args{
				userPassword: "abcd",
				hashedPasswd: "$2a$14$mE0dbCsxOlQSNXYgrZqzGu9fgSvOdttgKwExmtl3O0j6LINbTT16u",
			},
			want: want{
				check: false,
				msg:   "email or password is incorrect",
			},
			compareHashWithPassword: func(hashedPassword, password []byte) error {
				return errors.New("crypto/bcrypt: hashedPassword is not the hash of the given password")
			},
		},
	}
	for _, test := range subtests {
		compHashAndPass = test.compareHashWithPassword
		gotCheck, gotMsg := PasswordHelper.VerifyPassword(test.args.userPassword, test.args.hashedPasswd)
		var got want
		got.check = gotCheck
		got.msg = gotMsg
		assert.Equal(t, got, test.want)
	}

}

func Test_HashPassword(t *testing.T) {

	mockError := errors.New("mock error")

	testCases := []struct {
		name                 string
		password             string
		generateHashedPasswd func(password []byte, cost int) ([]byte, error)
		expectedErr          error
	}{
		{
			name:     "Success",
			password: "password",
			generateHashedPasswd: func(password []byte, cost int) ([]byte, error) {
				return nil, nil
			},
		},
		{
			name:     "Error",
			password: "password",
			generateHashedPasswd: func(password []byte, cost int) ([]byte, error) {
				return nil, mockError
			},
			expectedErr: mockError,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.name, func(t *testing.T) {
			genHashPasswd = tC.generateHashedPasswd
			_, err := PasswordHelper.HashPassword(tC.password)
			if !errors.Is(err, tC.expectedErr) {
				t.Errorf("expected error (%v), got error (%v)", tC.expectedErr, err)
			}
		})
	}
}
