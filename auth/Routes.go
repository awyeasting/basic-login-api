package auth

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	loginInfo := r.Context().Value("loginInfo").(LoginInfo)

	userID, err := CheckLoginInfo(r.Context(), &loginInfo)
	if err != nil {
		aErr, ok := err.(AuthError)
		if ok {
			if aErr.Type() == (UnconfirmedError{}).Type() {
				w.WriteHeader(http.StatusPreconditionRequired)
				log.Error(err)
				return
			}
			if aErr.Type() == (DeactivatedError{}).Type() {
				w.WriteHeader(http.StatusGone)
				log.Error(err)
				return
			}
		}
		w.WriteHeader(http.StatusUnauthorized)
		log.Error(err)
		return
	}

	// Create session and send back cookie
	CreateLoginSession(w, r, userID)
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {

	authSession, _ := cookieStore.Get(r, "authSession")
	delete(authSession.Values, "userID")
	authSession.Options.MaxAge = -1

	authSession.Save(r, w)
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	newUser := r.Context().Value("newUser").(NewUser)

	err := InsertNewUser(r.Context(), &newUser)
	if err != nil {
		pErr, ok := err.(*pq.Error)
		if ok {
			if pErr.Code.Name() == "unique_violation" {
				w.WriteHeader(http.StatusConflict)
				log.Debug("Attempt to create an account for an email that already is registered.")
				return
			}
			log.Info(pErr.Code.Name())
		}
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func UserHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)

	uInfo, err := GetUserInfo(r.Context(), userID)
	if err != nil {
		if errors.Is(err, DeactivatedError{}) {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		log.Error("Could not find user info for logged in user. ", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(uInfo)
}

func UserInfoChangeHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)
	userInfo := r.Context().Value("userInfo").(UserInfo)

	err := ChangeUserInfo(r.Context(), userID, &userInfo)
	if err != nil {
		pqErr, ok := err.(*pq.Error)
		if ok {
			if pqErr.Code.Name() == "unique_violation" {
				w.WriteHeader(http.StatusConflict)
				log.Debug("Attempt to change user info in a way which conflicts with constraints.")
				return
			}
		}
		log.Error("Error changing user information. ", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func UserConfirmationHandler(w http.ResponseWriter, r *http.Request) {
	confirmID := r.Context().Value("confirmationID").(string)

	log.Info("Processing email confirmation...")
	err := ConfirmEmail(r.Context(), confirmID)
	if err != nil {
		aErr, ok := err.(AuthError)
		if ok {
			if aErr.Type() == (ConfirmationExpiredError{}).Type() {
				w.WriteHeader(http.StatusGone)
				log.Info("Email confirmation expired")
				return
			}
			if aErr.Type() == (ConfirmationMissingError{}).Type() {
				w.WriteHeader(http.StatusNotFound)
				log.Info("Email confirmation gone")
				return
			}
		}
		log.Error("Unknown Confirmation Error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Info("Email Confirmed")
}
