package handlers

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"

	"github.com/openai/openai-go/v2"
	"github.com/openai/openai-go/v2/option"
)

var (
	API_KEY = os.Getenv("API_KEY")
)

func CompletionHandler(w http.ResponseWriter, r *http.Request) {
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
