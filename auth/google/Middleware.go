package google

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func DecodeGoogleClaims(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var gJWT GoogleJWT
		body, _ := ioutil.ReadAll(r.Body)
		err := json.Unmarshal(body, &gJWT)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Error("Could not decode Google JWT:", err)
			return
		}

		claims, err := ValidateGoogleJWT(gJWT.Token)
		if err != nil {
			w.WriteHeader(http.StatusForbidden)
			log.Error("Could not validate Google JWT:", err)
			return
		}

		ctx := context.WithValue(r.Context(), "googleClaims", claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
