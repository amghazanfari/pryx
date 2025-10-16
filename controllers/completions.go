package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/amghazanfari/pryx/models"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

type ChatCompletionRequest struct {
	Model    string                                   `json:"model"`
	Messages []openai.ChatCompletionMessageParamUnion `json:"messages"`
}

type ChatCompletion struct {
	ChatCompletionService *models.ChatCompletionService
	EndpointService       *models.EndpointService
}

func (cc ChatCompletion) Completion(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Read the body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	var ccRequest ChatCompletionRequest
	err = json.Unmarshal(body, &ccRequest)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if ccRequest.Model == "" || ccRequest.Messages == nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	endpoints, err := cc.EndpointService.ListByModel(ccRequest.Model)
	if err != nil {
		http.Error(w, "Error getting endpoints", http.StatusInternalServerError)
		return
	}

	endpoint := (*endpoints)[0]

	client := openai.NewClient(
		option.WithAPIKey(endpoint.APIKey),
		option.WithBaseURL(endpoint.URLAdress),
	)
	chatCompletion, err := client.Chat.Completions.New(context.TODO(), openai.ChatCompletionNewParams{
		Messages: ccRequest.Messages,
		Model:    endpoint.Name,
	})
	if err != nil {
		fmt.Println(err)
	}

	chatCompletionBytes, err := json.Marshal(chatCompletion)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Write(chatCompletionBytes)
}
