package auth

import (
	"errors"
	"net/http"

	log "github.com/sirupsen/logrus"
)

// Gets the currently logged in user
func GetSessionUser(r *http.Request) (string, error) {
	authSession, err := cookieStore.Get(r, "authSession")
	userID, ok := authSession.Values["userID"]
	if !ok {
		return "", errors.New("User not logged in")
	}
	user, ok := userID.(string)
	if !ok {
		log.WithField("userID", userID).Error("UserID stored in improper form")
		err = errors.New("UserID stored in improper form")
	}
	return user, err
}
