package auth

import (
	"net/http"

	log "github.com/sirupsen/logrus"
)

func CreateLoginSession(w http.ResponseWriter, r *http.Request, userID string) {
	// Create session and send back cookie
	log.Info("Creatings session for user ", userID)
	authSession, _ := cookieStore.Get(r, "authSession")
	authSession.Options.Path = "/"
	authSession.Values["userID"] = userID
	if !DEV {
		authSession.Options.Secure = true
	}
	authSession.Options.HttpOnly = true
	authSession.Options.MaxAge = LOGIN_SESSION_DURATION * 24 * 60 * 60

	err := authSession.Save(r, w)
	if err != nil {
		log.Error("Could not create login session for userID:", userID)
	}
}
