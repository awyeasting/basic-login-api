package google

import (
	"net/http"

	"github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

// Handles internal login checks after google claims clears for a given user.
func GoogleLoginHandler(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value("googleClaims").(GoogleClaims)

	userID, err := CheckLoginInfo(r.Context(), claims)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		log.Error("Could not log user in with google claims")
		return
	}

	// Create session and send back cookie
	CREATE_LOGIN_SESSION(w, r, userID)
}

// Handles internal registration after google claims clears for a given user.
func GoogleRegisterHandler(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value("googleClaims").(GoogleClaims)

	err := InsertNewUser(r.Context(), claims)
	if err != nil {
		pErr, ok := err.(*pq.Error)
		if ok {
			if pErr.Code.Name() == "unique_violation" {
				w.WriteHeader(http.StatusConflict)
				log.Debug("Attempt to create an account for a google email that already is registered.")
				return
			}
			log.Info(pErr.Code.Name())
		}
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
