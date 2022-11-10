package auth

import (
	"context"
	"fmt"
	"net/smtp"

	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

// Inserts new user into database
func InsertNewUser(ctx context.Context, newUser *NewUser) error {

	db := GetDBFromContext(ctx)

	// Hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), CURRENT_HASH_COST)
	if err != nil {
		return err
	}

	// Begin sql transaction so that failures are all or nothing
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Generate userID
	var userID string
	getUUIDQuery := `
	SELECT gen_random_uuid()`
	err = tx.QueryRow(getUUIDQuery).Scan(&userID)
	if err != nil {
		return err
	}

	// Create user
	createUserQuery := `
	INSERT INTO users(id, email, username, first_name, last_name, password, hash_method)
	VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err = tx.Exec(createUserQuery,
		userID,
		newUser.Email,
		newUser.Username,
		newUser.FirstName,
		newUser.LastName,
		hash,
		CURRENT_HASH_METHOD)
	if err != nil {
		return err
	}

	// Create email confirmation
	cb, err := CreateEmailConfirmation(tx, userID)
	if err != nil {
		return err
	}

	// Commit transaction
	err = tx.Commit()
	if err != nil {
		return err
	}

	// Send email confirmation
	cb()
	return err
}

// Send an email confirmation email to a given user
func SendRegistrationEmail(destEmail, userID, emailID string) error {

	subject := fmt.Sprintf("Confirm Email for %v", TITLE)
	link := fmt.Sprintf("%v/confirm/%v", FRONTEND_URL, emailID)
	body := fmt.Sprintf("Click below to confirm your email and use the site.\n\n%v",
		link)
	msg := fmt.Sprintf("From: %v\nTo: %v\nSubject: %v\n\n%v", EMAIL, destEmail, subject, body)

	log.Info(fmt.Sprintf("MAIL_SERVER: %v, MAIL_SERVER_PASS: %v, EMAIL: %v, EMAIL_PASS: %v", MAIL_SERVER, MAIL_SERVER_PORT, EMAIL, EMAIL_PASS))
	err := smtp.SendMail(fmt.Sprintf("%v:%v", MAIL_SERVER, MAIL_SERVER_PORT),
		smtp.PlainAuth("", EMAIL, EMAIL_PASS, MAIL_SERVER),
		EMAIL, []string{destEmail}, []byte(msg))

	if err != nil {
		log.Error("Error sending mail: ", err)
	}

	return err
}
