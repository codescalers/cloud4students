package validator

import (
	"net/mail"

	"github.com/caitlin615/nist-password-validator/password"
)

func ValidateMail(address string) bool {
	_, err := mail.ParseAddress(address)
	return err == nil
}

func ValidatePassword(Password string) error {
	// password should be ACII , min 5 , max 10
	validator := password.NewValidator(true, 5, 10)
	err := validator.ValidatePassword(Password)
	if err != nil {
		return err
	}
	return nil
}
