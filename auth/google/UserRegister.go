package google

import (
	"context"
	"errors"
	"strings"
)

func InsertNewUser(ctx context.Context, claims GoogleClaims) error {
	if claims.Email == "" {
		return errors.New("Empty claims received in InsertNewUser for google registration")
	}

	db := GetDBFromContext(ctx)

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
	INSERT INTO users(password, hash_method, emailConfirmed, id, email, username, first_name, last_name)
	VALUES ('google', 'text', true, $1, $2, $3, $4, $5)`
	_, err = tx.Exec(createUserQuery,
		userID,
		claims.Email,
		strings.Split(claims.Email, "@")[0],
		claims.FirstName,
		claims.LastName)
	if err != nil {
		return err
	}

	// Commit transaction
	err = tx.Commit()
	return err
}
