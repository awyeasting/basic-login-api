package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/awyeasting/basic-login-api/auth"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

var (
	DEV           = true
	TLS_CERT_PATH = ""
	TLS_KEY_PATH  = ""
	PORT          = "8080"
	DB_HOST       = "localhost"
	DB_PORT       = "5432"
	DB_USER       = "postgres"
	DB_PASS       = ""
	DB_NAME       = "basic-login-db"
	DB_CONN_STR   = ""
)

func SetDatabaseContext(client *sql.DB) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), "db", client)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func init() {
	// Load important .env file config
	godotenv.Load()

	if os.Getenv("DEV") != "" {
		DEV = strings.ToLower(os.Getenv("DEV")) == "true"
	}
	TLS_CERT_PATH = os.Getenv("TLS_CERT_PATH")
	TLS_KEY_PATH = os.Getenv("TLS_KEY_PATH")
	if os.Getenv("PORT") != "" {
		PORT = os.Getenv("PORT")
	}
	if os.Getenv("DB_HOST") != "" {
		DB_HOST = os.Getenv("DB_HOST")
	}
	if os.Getenv("DB_PORT") != "" {
		DB_PORT = os.Getenv("DB_PORT")
	}
	if os.Getenv("DB_USER") != "" {
		DB_USER = os.Getenv("DB_USER")
	}
	if os.Getenv("DB_NAME") != "" {
		DB_NAME = os.Getenv("DB_NAME")
	}
	DB_PASS = os.Getenv("DB_PASS")
	DB_CONN_STR = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s", DB_HOST, DB_PORT, DB_USER, DB_PASS, DB_NAME)
	if DEV {
		DB_CONN_STR += " sslmode=disable"
	}

	// Setup logger
	log.SetFormatter(&log.JSONFormatter{})

	log.SetOutput(os.Stdout)

	log.SetLevel(log.InfoLevel)
}

func main() {
	r := chi.NewRouter()

	log.Info("Connecting to postgres database...")
	db, err := sql.Open("postgres", DB_CONN_STR)
	if err != nil {
		log.Panic(err)
	}
	defer db.Close()

	log.Info("Pinging database...")
	err = db.Ping()
	if err != nil {
		log.Panic(err)
	}

	// Recovers from panics and returns an HTTP 500 status if possible
	log.Info("Setting up middleware...")
	r.Use(middleware.Recoverer)
	// Times out requests if they go on too long
	r.Use(middleware.Timeout(REQUEST_TIMEOUT))
	// Puts the database handle in request
	r.Use(SetDatabaseContext(db))
	// Configure CORS policy
	if DEV {
		r.Use(cors.Handler(cors.Options{
			// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
			AllowedOrigins: []string{"https://*", "http://*"},
			// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
			AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders: []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
			//ExposedHeaders:   []string{"Link"},
			AllowCredentials: true,
			MaxAge:           300, // Maximum value not ignored by any of major browsers
		}))
	} else {
		r.Use(cors.Handler(cors.Options{
			// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
			AllowedOrigins: []string{"https://*"},
			// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
			AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders: []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
			//ExposedHeaders:   []string{"Link"},
			AllowCredentials: true,
			MaxAge:           300, // Maximum value not ignored by any of major browsers
		}))
	}

	// Mount the subrouters
	log.Info("Mounting subrouters...")
	r.Mount(AUTH_PATH, auth.AuthRouter())

	log.Info(fmt.Sprintf("Listening on port:%v...", PORT))
	if DEV {
		err := http.ListenAndServe(fmt.Sprintf(":%v", PORT), r)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		err := http.ListenAndServeTLS(fmt.Sprintf(":%v", PORT), TLS_CERT_PATH, TLS_KEY_PATH, r)
		if err != nil {
			log.Fatal(err)
		}
	}
}
