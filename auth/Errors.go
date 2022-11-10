package auth

import "fmt"

type AuthError interface {
	Error() string
	Type() string
}

type DeactivatedError struct {
	userID string
}

type UnconfirmedError struct {
	userID string
	email  string
}

type UnconfirmedAndNotSentError struct {
	userID string
	email  string
}

type ConfirmationExpiredError struct {
	userID string
}

type ConfirmationMissingError struct {
	confirmID string
}

func (e DeactivatedError) Error() string {
	return fmt.Sprintf("Account deactivated for user %v", e.userID)
}

func (e DeactivatedError) Type() string {
	return "deactivated_error"
}

func (e UnconfirmedError) Error() string {
	return fmt.Sprintf("Email unconfirmed for user %v with email %v", e.userID, e.email)
}

func (e UnconfirmedError) Type() string {
	return "unconfirmed_error"
}

func (e UnconfirmedAndNotSentError) Error() string {
	return fmt.Sprintf("Email unconfirmed and not sent and for user %v with email %v", e.userID, e.email)
}

func (e UnconfirmedAndNotSentError) Type() string {
	return "unconfirmed_and_not_sent_error"
}

func (e ConfirmationExpiredError) Error() string {
	return fmt.Sprintf("Email confirmation expired for user %v", e.userID)
}

func (e ConfirmationExpiredError) Type() string {
	return "confirmation_expired_error"
}

func (e ConfirmationMissingError) Error() string {
	return fmt.Sprintf("Email confirmation code missing %v", e.confirmID)
}

func (e ConfirmationMissingError) Type() string {
	return "confirmation_missing_error"
}
