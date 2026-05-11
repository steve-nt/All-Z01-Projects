package controllers

import (
	"errors"
	"forum/src/models"

	"golang.org/x/crypto/bcrypt"
)

func CompareRegistrationPasswords(pass1, pass2 string) bool {
	return pass1 == pass2
}

func Auth(email, password string) error {
	var err error
	if !models.IsEmailRegistered(email) {
		return models.ErrorNotRegistered
	}
	var user models.User
	user, err = models.GetUserPasswordByEmail(email)
	if err != nil {
		return err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Hash), []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return models.ErrorWrongPassword
		}
		return err
	}
	return nil
}
