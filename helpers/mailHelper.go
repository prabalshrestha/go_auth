package helpers

import (
	"fmt"
	"login/models"
	"net/smtp"
)

type mailHelper interface {
	SendRegistrationMail(user models.User, password string) error
	SendResetPasswordMail(user models.User, password string) error
}
type mailHelperImpl struct {
}

var (
	MailHelper mailHelper = &mailHelperImpl{}
)

func (helper *mailHelperImpl) SendRegistrationMail(user models.User, password string) error {
	// Sender data.
	from := "pstestmail04@gmail.com"
	emailPassword := "pstestmail04@test.com"

	// Receiver email address.
	to := []string{
		user.Email,
	}

	// smtp server configuration.
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	subject := "Subject: Account Created\n"
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body := "<html><body><h4>Hi " + user.Name + ",</h4><p>Your account was successfully created. Please use the following credentials to login and change your password once you are in</p><h5>Email: " + user.Email + "</h5><h5>Password: " + password + "</h5><br></br></body></html>"
	msg := []byte(subject + mime + body)

	// Authentication.
	auth := smtp.PlainAuth("", from, emailPassword, smtpHost)

	// Sending email.
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, msg)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println("Email Sent Successfully!")
	return nil
}

func (helper *mailHelperImpl) SendResetPasswordMail(user models.User, password string) error {
	// Sender data.
	from := "pstestmail04@gmail.com"
	emailPassword := "pstestmail04@test.com"

	// Receiver email address.
	to := []string{
		user.Email,
	}

	// smtp server configuration.
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	subject := "Subject: Password Reset\n"
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body := "<html><body><h4>Hi " + user.Name + ",</h4><p>Your password was successfully reset. Please use the following credentials to login and change your password once you are in</p><h5>New Password: " + password + "</h5><br></br></body></html>"
	msg := []byte(subject + mime + body)

	// Authentication.
	auth := smtp.PlainAuth("", from, emailPassword, smtpHost)

	// Sending email.
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, msg)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println("Email Sent Successfully!")
	return nil
}
