package auth

import "regexp"

var (
	emailRegex = regexp.MustCompile("(?:[a-z0-9!#$%&'*+/=?^_`{|}~-]+(?:\\.[a-z0-9!#$%&'*+/=?^_`{|}~-]+)*|\"(?:[\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x21\\x23-\\x5b\\x5d-\\x7f]|\\\\[\\x01-\\x09\\x0b\\x0c\\x0e-\\x7f])*\")@(?:(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\\.)+[a-z0-9](?:[a-z0-9-]*[a-z0-9])?|\\[(?:(?:(2(5[0-5]|[0-4][0-9])|1[0-9][0-9]|[1-9]?[0-9]))\\.){3}(?:(2(5[0-5]|[0-4][0-9])|1[0-9][0-9]|[1-9]?[0-9])|[a-z0-9-]*[a-z0-9]:(?:[\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x21-\\x5a\\x53-\\x7f]|\\\\[\\x01-\\x09\\x0b\\x0c\\x0e-\\x7f])+)\\])")
)

func ValidatePassword(password string) bool {

	if len(password) < 8 || len(password) > 72 {
		return false
	}

	return true
}

func ValidateEmail(email string) bool {

	if len(email) > 254 || len(email) < 5 {
		return false
	}

	return emailRegex.MatchString(email)
}

func ValidateName(name string) bool {

	if len(name) > 128 {
		return false
	}

	return true
}

func ValidateUsername(name string) bool {

	if len(name) > 40 || len(name) < 1 {
		return false
	}

	return true
}
