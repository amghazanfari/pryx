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
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	csrf "github.com/utrack/gin-csrf"

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

	// setup controllers
	userC := controllers.Users{
		UserService:          &userService,
		SessionService:       &sessionService,
		PasswordResetService: &pwResetService,
		EmailService:         emailService,
		EndpointService:      &endpointService,
	}

	ModelListTpl, err := views.ParseFS(templates.FS, "layout.gohtml", "proxy-list.gohtml")
	if err != nil {
		panic(err)
	}
	userC.Templates.ModelList = ModelListTpl

	r := gin.Default()
	r.Use(umw.SetUser())
	r.Static("/static", "./static")

	v1 := r.Group("/v1")
	{
		modelsGroup := v1.Group("/models")
		{
			modelsGroup.POST("/add", gin.WrapF(modelC.Create))
			modelsGroup.GET("/", gin.WrapF(modelC.List))
			modelsGroup.GET("/:model", gin.WrapF(modelC.Retrieve))
		}

		endpointsGroup := v1.Group("/endpoints")
		{
			endpointsGroup.GET("/", gin.WrapF(endpointC.List))
			endpointsGroup.DELETE("/", gin.WrapF(endpointC.Delete))
		}

		chatGroup := v1.Group("/chat/completions")
		{
			chatGroup.GET("/", gin.WrapF(chatCompletionC.Completion))
		}
	}
	ui := r.Group("/ui")
	store := cookie.NewStore([]byte("secret"))
	ui.Use(sessions.Sessions("mysession", store))
	ui.Use(csrf.Middleware(csrf.Options{
		Secret: cfg.CSRF.Key,
		ErrorFunc: func(c *gin.Context) {
			c.String(400, "CSRF token mismatch")
			c.Abort()
		},
	}))
	{
		uiGroup := ui.Group("/models")
		{
			uiGroup.GET("/", gin.WrapF(userC.ModelList))
		}
	}

	http.ListenAndServe(":8080", r)

}
