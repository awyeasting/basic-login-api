package auth

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/go-chi/chi"
	log "github.com/sirupsen/logrus"
)

func DecodeNewUserInfo(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var user NewUser
		body, _ := ioutil.ReadAll(r.Body)
		err := json.Unmarshal(body, &user)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Error(err)
			return
		}

		// Validate that the first and last names meet current limitations
		if !ValidateName(user.FirstName) || !ValidateName(user.LastName) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Validate that the username meets the current requirements
		if !ValidateUsername(user.Username) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Validate that the password meets current registration criteria
		if !ValidatePassword(user.Password) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Validate that the email is realistic
		if !ValidateEmail(user.Email) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), "newUser", user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func DecodeUserInfo(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var user UserInfo
		body, _ := ioutil.ReadAll(r.Body)
		err := json.Unmarshal(body, &user)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Error(err)
			return
		}

		// Validate that the first and last names meet current limitations
		if user.FirstName != nil {
			if !ValidateName(*user.FirstName) {
				w.WriteHeader(http.StatusBadRequest)
				log.Info("Invalid first name ", user.FirstName)
			}
		}
		if user.LastName != nil {
			if !ValidateName(*user.LastName) {
				w.WriteHeader(http.StatusBadRequest)
				log.Info("Invalid last name ", user.LastName)
			}
		}

		// Validate that the username meets the current requirements
		if user.Username != nil {
			if !ValidateUsername(*user.Username) {
				w.WriteHeader(http.StatusBadRequest)
				log.Info("Invalid username ", user.Username)
				return
			}
		}

		// Validate that the email is realistic
		if user.Email != nil {
			if !ValidateEmail(*user.Email) {
				w.WriteHeader(http.StatusBadRequest)
				log.Info("Invalid email ", user.Email)
				return
			}
		}

		ctx := context.WithValue(r.Context(), "userInfo", user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func DecodeLoginInfo(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var info LoginInfo
		body, _ := ioutil.ReadAll(r.Body)
		err := json.Unmarshal(body, &info)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Error(err)
			return
		}

		ctx := context.WithValue(r.Context(), "loginInfo", info)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func RequireValidUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, err := GetSessionUser(r)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			log.Info("Rejected request, ", err)
			return
		}

		ctx := context.WithValue(r.Context(), "userID", userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func DecodeUserConfirmation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		confirmationID := chi.URLParam(r, "confirmationID")
		if confirmationID == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), "confirmationID", confirmationID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
