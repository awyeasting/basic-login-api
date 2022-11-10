package google

import (
	"context"
	"strings"

	log "github.com/sirupsen/logrus"
)

// Checks whether a given GoogleClaims can be used to log in.
func CheckLoginInfo(ctx context.Context, claims GoogleClaims) (string, error) {

	db := GetDBFromContext(ctx)

	var userID, pass string
	query := `
	SELECT id, password FROM users WHERE email=$1 AND emailConfirmed=true
	`
	err := db.QueryRow(query, claims.Email).Scan(&userID, &pass)
	if err != nil {
		log.Error("Error selecting account for google login. Likely account not yet made.", err)
		return "", InvalidGoogleLogin{}
	}

	if strings.ToLower(strings.TrimSpace(pass)) != "google" {
		return "", InvalidGoogleLogin{}
	}

	return userID, nil
}
