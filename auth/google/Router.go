package google

import (
	"net/http"
	"os"

	"github.com/go-chi/chi"
)

var (
	CREATE_LOGIN_SESSION func(http.ResponseWriter, *http.Request, string)
	GOOGLE_CLIENT_ID     string = ""
)

func GoogleRouter(loginSessionCreator func(http.ResponseWriter, *http.Request, string)) http.Handler {

	CREATE_LOGIN_SESSION = loginSessionCreator
	GOOGLE_CLIENT_ID = os.Getenv("GOOGLE_CLIENT_ID")

	r := chi.NewRouter()

	r.With(DecodeGoogleClaims).Post("/login", GoogleLoginHandler)
	r.With(DecodeGoogleClaims).Post("/register", GoogleRegisterHandler)

	return r
}
