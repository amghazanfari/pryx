package main

import (
	"net/http"

	"github.com/amghazanfari/pryx/controllers"
	"github.com/amghazanfari/pryx/migrations"
	"github.com/amghazanfari/pryx/models"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	_ "github.com/jackc/pgx/v4/stdlib"
)

func main() {
	// connect to database
	cfg := models.DefaultPostgresConfig()
	db, err := models.Open(cfg)
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

	//setup controller
	modelC := controllers.Model{
		ModelService: &modelService,
	}

	chatCompletionC := controllers.ChatCompletion{
		ChatCompletionService: &chatCompletionService,
		EndpointService:       &endpointService,
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Route("/v1/models", func(r chi.Router) {
		r.Post("/add", modelC.Create)
		r.Get("/", modelC.List)
		r.Get("/{model}", modelC.Retrieve)
	})
	r.Route("/v1/chat/completions", func(r chi.Router) {
		r.Get("/", chatCompletionC.Completion)
	})

	http.ListenAndServe(":8080", r)

}
