package main

import (
	"net/http"
	"os"
	"strconv"

	"github.com/amghazanfari/pryx/controllers"
	"github.com/amghazanfari/pryx/migrations"
	"github.com/amghazanfari/pryx/models"
	"github.com/amghazanfari/pryx/templates"
	"github.com/amghazanfari/pryx/views"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/csrf"
	"github.com/joho/godotenv"

	_ "github.com/jackc/pgx/v4/stdlib"
)

type config struct {
	PSQL models.PostgresConfig
	SMTP models.SMTPConfig
	CSRF struct {
		Key    string
		Secure bool
	}
	SERVER struct {
		Address string
	}
}

func loadEnvConfig() (config, error) {
	var cfg config

	err := godotenv.Load()
	if err != nil {
		return cfg, err
	}
	cfg.PSQL = models.DefaultPostgresConfig()

	cfg.SMTP.Host = os.Getenv("SMTP_HOST")
	portString := os.Getenv("SMTP_PORT")
	cfg.SMTP.Port, err = strconv.Atoi(portString)

	if err != nil {
		return cfg, err
	}

	cfg.SMTP.Username = os.Getenv("SMTP_USERNAME")
	cfg.SMTP.Password = os.Getenv("SMTP_PASSWORD")

	cfg.CSRF.Key = "Nk2uFnisr5156l3xeXKYtj4HS4o5CTAV"
	cfg.CSRF.Secure = false
	cfg.SERVER.Address = ":8080"
	return cfg, nil
}

func main() {
	// connect to database
	cfg, err := loadEnvConfig()
	if err != nil {
		panic(err)
	}

	db, err := models.Open(cfg.PSQL)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	//migrate database
	err = models.MigrateFS(db, migrations.FS, ".")
	if err != nil {
		panic(err)
	}

	// setup services
	modelService := models.ModelService{
		DB: db,
	}
	endpointService := models.EndpointService{
		DB: db,
	}
	chatCompletionService := models.ChatCompletionService{
		DB: db,
	}
	userService := models.UserService{
		DB: db,
	}

	sessionService := models.SessionService{
		DB: db,
	}

	pwResetService := models.PasswordResetService{
		DB: db,
	}

	emailService := models.NewEmailService(cfg.SMTP)

	//setup controller
	modelC := controllers.Model{
		ModelService: &modelService,
	}

	endpointC := controllers.Endpoint{
		EndpointService: &endpointService,
	}

	//setup middleware
	umw := controllers.UserMiddleware{
		SessionService: &sessionService,
	}

	chatCompletionC := controllers.ChatCompletion{
		ChatCompletionService: &chatCompletionService,
		EndpointService:       &endpointService,
	}

	csrfKey := cfg.CSRF.Key
	csrfMw := csrf.Protect(
		[]byte(csrfKey),
		csrf.Secure(cfg.CSRF.Secure),
	)

	// setup controllers
	userC := controllers.Users{
		UserService:          &userService,
		SessionService:       &sessionService,
		PasswordResetService: &pwResetService,
		EmailService:         emailService,
	}

	ModelListTpl, err := views.ParseFS(templates.FS, "layout.gohtml", "proxy-list.gohtml")
	if err != nil {
		panic(err)
	}
	userC.Templates.ModelList = ModelListTpl

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(csrfMw)
	r.Use(umw.SetUser)

	r.Route("/v1/models", func(r chi.Router) {
		r.Post("/add", modelC.Create)
		r.Get("/", modelC.List)
		r.Get("/{model}", modelC.Retrieve)
	})
	r.Route("/v1/endpoints", func(r chi.Router) {
		r.Get("/", endpointC.List)
	})
	r.Route("/v1/chat/completions", func(r chi.Router) {
		r.Get("/", chatCompletionC.Completion)
	})
	r.Route("/ui", func(r chi.Router) {
		r.Get("/models", userC.ModelList)
	})

	http.ListenAndServe(":8080", r)

}
