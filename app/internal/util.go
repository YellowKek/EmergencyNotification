package internal

import (
	"app/internal/model"
	"net/mail"
)

func ValidateEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func ValidateUser(user model.User) bool {
	if !ValidateEmail(user.Email) {
		return false
	}
	if user.Name != "" && user.Surname != "" && user.Password != "" {
		return true
	}
	return false
}
