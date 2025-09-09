package main

import (
	"context"
	"net/http"
	"log"
	"os"
	"encoding/json"

	"github.com/openai/openai-go/v2"
	"github.com/openai/openai-go/v2/option"
)

var (
API_KEY = os.Getenv("API_KEY")
)

func main() {
	http.HandleFunc("/", handler)
    log.Fatal(http.ListenAndServe(":8080", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	promptMessage := "hello"
	client := openai.NewClient(
		option.WithAPIKey(API_KEY), // defaults to os.LookupEnv("OPENAI_API_KEY")
		option.WithBaseURL("https://openrouter.ai/api/v1"),
	)
	chatCompletion, err := client.Chat.Completions.New(context.TODO(), openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(promptMessage),
		},
		Model: "openai/gpt-oss-120b:free",
	})
	if err != nil {
		panic(err.Error())
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(chatCompletion)
}