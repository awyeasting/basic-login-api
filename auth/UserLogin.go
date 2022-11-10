package auth

import (
	"context"

	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

func CheckLoginInfo(ctx context.Context, loginInfo *LoginInfo) (string, error) {

	db := GetDBFromContext(ctx)

	var id string
	var password string
	var hash_method string
	var deactivated, emailConfirmed bool
	query := `
	SELECT id, password, hash_method, deactivated, emailConfirmed FROM users WHERE email=$1`
	err := db.QueryRow(query, loginInfo.Email).Scan(&id, &password, &hash_method, &deactivated, &emailConfirmed)
	if err != nil {
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(password), []byte(loginInfo.Password))
	if err != nil {
		return "", err
	}
	if !emailConfirmed {
		log.Info("Attempt to log in to unconfirmed email. Attempting email confirmation creation...")
		err = EnsureVerification(ctx, id)
		return "", err
	}
	if deactivated {
		log.Info("Login attempted to deactivated account. ", err)
		return "", DeactivatedError{userID: id}
	}
	return id, err
}

func GetUserInfo(ctx context.Context, userID string) (*UserInfo, error) {
	db := GetDBFromContext(ctx)

	var deactivated, emailConfirmed bool
	var uInfo UserInfo
	query := `
	SELECT first_name, last_name, username, email, deactivated, emailConfirmed FROM users WHERE id=$1`
	err := db.QueryRow(query, userID).Scan(&uInfo.FirstName, &uInfo.LastName, &uInfo.Username, &uInfo.Email, &deactivated, &emailConfirmed)
	if err != nil {
		return nil, err
	}

	email := AnonymizeEmail(*(uInfo.Email))
	uInfo.Email = &email
	return &uInfo, err
}

func EnsureVerification(ctx context.Context, userID string) error {
	db := GetDBFromContext(ctx)

	// Create transaction to ensure that everything either works or doesn't
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		log.Error(err)
		return UnconfirmedAndNotSentError{userID: userID}
	}
	defer tx.Rollback()

	// Create an email confirmation
	cb, err := CreateSafeEmailConfirmation(tx, userID)
	if err != nil {
		log.Error(err)
		return UnconfirmedAndNotSentError{userID: userID}
	}
	err = tx.Commit()
	if err != nil {
		log.Error(err)
		return UnconfirmedAndNotSentError{userID: userID}
	}

	// If the email confirmation was created then send the email
	cb()
	return UnconfirmedError{userID: userID}
}
