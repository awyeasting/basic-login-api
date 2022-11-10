package auth

import (
	"context"
	"database/sql"
	"time"

	log "github.com/sirupsen/logrus"
)

type dbHandle interface {
	QueryRow(string, ...any) *sql.Row
	Exec(string, ...any) (sql.Result, error)
}

func ConfirmEmail(ctx context.Context, confirmID string) error {
	db := GetDBFromContext(ctx)

	// Begin sql transaction so that failures are all or nothing
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Get the email confirmation information
	var userID string
	var expireTime, currentTime time.Time
	getConfirmationQuery := `
	SELECT userID, expireTime, NOW() as currentTime FROM emailConfirmations WHERE id=$1`
	err = tx.QueryRow(getConfirmationQuery, confirmID).Scan(
		&userID,
		&expireTime,
		&currentTime)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			log.Info("No email confirmation found")
			return ConfirmationMissingError{confirmID: confirmID}
		}
		return err
	}

	// Remove the confirmation if it has expired
	removeConfirmationQuery := `
	DELETE FROM emailConfirmations WHERE id=$1
	`
	if expireTime.Before(currentTime) {
		log.Info("Removing expired confirmation ", confirmID)
		log.Info(expireTime.Format("2006-01-02 3:04:05 pm -0700"), "\t", currentTime.Format("2006-01-02 3:04:05 pm -0700"))
		go db.Exec(removeConfirmationQuery, confirmID)
		return ConfirmationExpiredError{userID: userID}
	}

	// Otherwise confirm the email and remove the confirmation
	confirmEmailQuery := `
	UPDATE users SET emailConfirmed=true WHERE id=$1
	`
	log.Info("Attempting email confirmation for user ", userID)
	_, err = tx.Exec(confirmEmailQuery, userID)
	if err != nil {
		return err
	}

	_, err = tx.Exec(removeConfirmationQuery, confirmID)
	if err != nil {
		return err
	}

	err = tx.Commit()
	return err
}

func CreateEmailConfirmation(db dbHandle, userID string) (func(), error) {
	// Create email confirmation code
	getConfirmationInfoQuery := `
	SELECT gen_random_uuid(), email FROM users WHERE id=$1
	`
	var confirmID, email string
	err := db.QueryRow(getConfirmationInfoQuery).Scan(&confirmID, &email)
	if err != nil {
		return nil, err
	}

	// Create email confirmation table entry
	confirmEmailQuery := `
	INSERT INTO emailConfirmations(id, userID)
	VALUES ($1, $2)
	`
	_, err = db.Exec(confirmEmailQuery, confirmID, userID)
	if err != nil {
		return nil, err
	}

	// Send callback function to send email
	return func() {
		go SendRegistrationEmail(email, userID, confirmID)
	}, err
}

func CreateSafeEmailConfirmation(db dbHandle, userID string) (func(), error) {
	// Create email confirmation code
	getConfirmationInfoQuery := `
	SELECT gen_random_uuid(), U.email, E.id, E.expireTime, NOW()
	FROM users AS U 
	LEFT JOIN emailConfirmations AS E
	ON U.id = E.userID
	WHERE U.id=$1
	`
	var confirmID, email string
	var existingConfirmID sql.NullString
	var expireTime sql.NullTime
	var currentTime time.Time
	err := db.QueryRow(getConfirmationInfoQuery, userID).Scan(&confirmID, &email, &existingConfirmID, &expireTime, &currentTime)
	if err != nil {
		log.Error("Could not get user confirmation info")
		return nil, err
	}

	if !existingConfirmID.Valid {
		// Create email confirmation table entry
		confirmEmailQuery := `
		INSERT INTO emailConfirmations(id, userID)
		VALUES ($1, $2)
		`
		_, err = db.Exec(confirmEmailQuery, confirmID, userID)
		if err != nil {
			log.Error("Couldn't create new email confirmation")
			return nil, err
		}

		// Send callback function to send email
		return func() {
			go SendRegistrationEmail(email, userID, confirmID)
		}, err
	}

	// If the confirmation code hasn't expired yet then
	if expireTime.Time.After(currentTime) {
		return func() {
			go SendRegistrationEmail(email, userID, existingConfirmID.String)
		}, err
	}

	updateConfirmEmailQuery := `
	UPDATE emailConfirmations SET id=$1, expireTime=DEFAULT WHERE userID=$2
	`
	_, err = db.Exec(updateConfirmEmailQuery, confirmID, userID)
	if err != nil {
		log.Error("Couldn't update email confirmation")
		return nil, err
	}

	// Send callback function to send email
	return func() {
		go SendRegistrationEmail(email, userID, confirmID)
	}, err
}
