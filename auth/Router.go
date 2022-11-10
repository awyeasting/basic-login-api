package auth

import (
	"net/http"
	"os"
	"strings"

	"github.com/awyeasting/basic-login-api/auth/google"
	"github.com/go-chi/chi"
	"github.com/gorilla/sessions"
)

var (
	SESSION_SECRET   = ""
	DEV              = false
	cookieStore      *sessions.CookieStore
	EMAIL            = ""
	EMAIL_PASS       = ""
	MAIL_SERVER      = ""
	MAIL_SERVER_PORT = "587"
	TITLE            = ""
	FRONTEND_URL     = "http://localhost:3000"
)

// Create router for authentication paths
func AuthRouter() http.Handler {

	SESSION_SECRET = os.Getenv("SESSION_SECRET")
	EMAIL = os.Getenv("EMAIL")
	EMAIL_PASS = os.Getenv("EMAIL_PASS")
	MAIL_SERVER = os.Getenv("MAIL_SERVER")
	if os.Getenv("MAIL_SERVER_PORT") != "" {
		MAIL_SERVER_PORT = os.Getenv("MAIL_SERVER_PORT")
	}
	DEV = strings.ToLower(os.Getenv("DEV")) == "true"
	TITLE = os.Getenv("TITLE")
	if os.Getenv("FRONTEND_URL") != "" {
		FRONTEND_URL = os.Getenv("FRONTEND_URL")
	}
	cookieStore = sessions.NewCookieStore([]byte(SESSION_SECRET))

	r := chi.NewRouter()

	//r.Use(middleware.AllowContentType("application/json"))

	r.With(DecodeLoginInfo).Post("/login", LoginHandler)
	r.Post("/logout", LogoutHandler)
	r.With(DecodeNewUserInfo).Post("/register", RegisterHandler)
	r.With(RequireValidUser).Get("/user", UserHandler)
	r.With(RequireValidUser).With(DecodeUserInfo).Post("/user", UserInfoChangeHandler)
	r.With(DecodeUserConfirmation).Post("/confirm/{confirmationID}", UserConfirmationHandler)
	r.Mount("/google", google.GoogleRouter(CreateLoginSession))

	return r
}
