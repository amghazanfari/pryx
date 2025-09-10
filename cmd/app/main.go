package main

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"

	"github.com/openai/openai-go/v2"
	"github.com/openai/openai-go/v2/option"
	"pryx/internal/db"
)

var (
	API_KEY           = os.Getenv("API_KEY")
	POSTGRES_USER     = os.Getenv("POSTGRES_USER")
	POSTGRES_PASSWORD = os.Getenv("POSTGRES_PASSWORD")
	POSTGRES_HOST     = os.Getenv("POSTGRES_HOST")
	POSTGRES_PORT     = os.Getenv("POSTGRES_PORT")
)

func init() {
	// JSON logs to stdout, level from env
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)

	lvl, err := log.ParseLevel("info")
	if err != nil {
		lvl = log.InfoLevel
	}
	log.SetLevel(lvl)
}


func main() {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/postgres?sslmode=disable",
		POSTGRES_USER, POSTGRES_PASSWORD, POSTGRES_HOST, POSTGRES_PORT)

	conn, err := db.Connect(dsn)
	if err != nil {
		log.WithError(err).Fatal("db connect failed")
	}
	defer conn.Close()

	http.HandleFunc("/", handler)

	log.WithField("port", 8080).Info("starting server")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.WithError(err).Fatal("http server stopped")
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	log.WithFields(log.Fields{
		"method": r.Method,
		"path":   r.URL.Path,
	}).Info("incoming request")

	promptMessage := "hello"
	client := openai.NewClient(
		option.WithAPIKey(API_KEY),
		option.WithBaseURL("https://openrouter.ai/api/v1"),
	)

	chatCompletion, err := client.Chat.Completions.New(r.Context(), openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(promptMessage),
		},
		Model: "deepseek/deepseek-chat-v3-0324:free",
	})
	if err != nil {
		log.WithError(err).Error("chat completion failed")
		http.Error(w, "upstream error", http.StatusBadGateway)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(chatCompletion); err != nil {
		log.WithError(err).Error("write response failed")
	}
}
