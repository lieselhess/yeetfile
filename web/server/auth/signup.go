package auth

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"yeetfile/shared"
	"yeetfile/web/db"
	"yeetfile/web/mail"
	"yeetfile/web/utils"
)

var MissingField = errors.New("missing username or email")

// SignupWithEmail uses values from the Signup struct to complete registration
// of a new user. A hash is generated from the provided password and entered
// into the "users" db table.
func SignupWithEmail(signup shared.Signup) error {
	// When signing up with email, no part of the signup struct can be empty
	if utils.IsStructMissingAnyField(signup) {
		return MissingField
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(signup.LoginKeyHash), 8)
	if err != nil {
		return err
	}

	code, err := db.NewVerification(signup.Identifier, hash, signup.ProtectedKey, false)
	if err != nil {
		return err
	}

	err = mail.SendVerificationEmail(code, signup.Identifier)
	return err
}

// SignupAccountIDOnly creates a new user with only an account ID as the user's
// login credential. Returns the user's (temporary) account ID, an image
// of their captcha code, and an error.
func SignupAccountIDOnly() (string, string, error) {
	id := db.CreateUniqueUserID()

	code, err := db.NewVerification(id, nil, nil, false)
	if err != nil {
		return "", "", err
	}

	captchaBase64 := GenerateCaptchaImage(code)
	return id, captchaBase64, nil
}
