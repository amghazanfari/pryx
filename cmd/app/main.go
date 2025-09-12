package main

import (
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"

	"pryx/config"
	"pryx/internal/db"
	"pryx/internal/handlers"
	"pryx/internal/auth"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

var (
	API_KEY           = os.Getenv("API_KEY")
	POSTGRES_USER     = os.Getenv("POSTGRES_USER")
	POSTGRES_PASSWORD = os.Getenv("POSTGRES_PASSWORD")
	POSTGRES_HOST     = os.Getenv("POSTGRES_HOST")
	POSTGRES_PORT     = os.Getenv("POSTGRES_PORT")
	DB_AUTOMIGRATE    = os.Getenv("DB_AUTOMIGRATE")
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)

	lvl, err := log.ParseLevel("info")
	if err != nil {
		lvl = log.InfoLevel
	}
	log.SetLevel(lvl)
}

func main() {
	db_cfg := config.DBFromEnv()
	conn, err := db.Open(db_cfg)
	if err != nil {
		log.WithError(err).Fatal("db connect failed")
	}

	h := handlers.New(conn)
	r := chi.NewRouter()

	r.Use(middleware.Logger)

	r.Post("/admin/users", h.CreateUser())
	r.Post("/admin/keys", h.CreateAPIKey())

	protected := chi.NewRouter()

	protected.Use(auth.Middleware(conn, "completion:invoke"))
	protected.Post("/", h.CompletionHandler())

	modelsRouter := chi.NewRouter()
	modelsRouter.Use(auth.Middleware(conn, "model:write"))
	modelsRouter.Post("/", h.AddModelHandler())

	r.Mount("/", protected)
	r.Mount("/models", modelsRouter)

	if DB_AUTOMIGRATE == "true" {
		if err := db.AutoMigrateAll(conn); err != nil {
			log.WithError(err).Fatal("auto-migrate failed")
		}
		log.Info("auto-migrate completed")
	}

	http.ListenAndServe(":8080", r)
}
