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

func (h *Handler) CompletionHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.WithFields(log.Fields{
			"method": r.Method,
			"path":   r.URL.Path,
		}).Info("incoming request")

		var chatRequest openai.ChatCompletionNewParams

		if err := json.NewDecoder(r.Body).Decode(&chatRequest); err != nil {
			http.Error(w, `{"error":"invalid JSON"}`, http.StatusBadRequest)
			return
		}
		client := openai.NewClient(
			option.WithAPIKey(API_KEY),
			option.WithBaseURL("https://openrouter.ai/api/v1"),
		)

		chatCompletion, err := client.Chat.Completions.New(r.Context(), chatRequest)
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
}
