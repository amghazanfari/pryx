package controllers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/amghazanfari/pryx/migrations"
	"github.com/amghazanfari/pryx/models"
	_ "github.com/jackc/pgx/v4/stdlib"
)

var db *sql.DB

func TestMain(m *testing.M) {
	cfg := models.PostgresConfig{
		Host:     "localhost",
		Port:     "5432",
		User:     "user1",
		Password: "changeme",
		Database: "pryx_test",
		SSLMode:  "disable",
	}

	var err error
	db, err = models.Open(cfg)
	if err != nil {
		panic(err)
	}

	err = models.MigrateFS(db, migrations.FS, ".")
	if err != nil {
		panic(err)
	}

	code := m.Run()

	err = models.TearDownDB(db, migrations.FS, ".")
	if err != nil {
		panic(err)
	}

	if db != nil {
		db.Close()
	}

	os.Exit(code)
}

func TestCreate(t *testing.T) {
	w := httptest.NewRecorder()

	data := map[string]string{
		"model_name":    "mamad2",
		"endpoint_name": "qwen/qwen-2.5-coder-32b-instruct:free",
		"api_key":       "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",
		"url_address":   "https://openrouter.ai/api/v1",
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		t.Fatalf("Error marshalling JSON: %v", err)
	}

	r, err := http.NewRequest(http.MethodPost, "/v1/models/add", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatalf("Failed to make a request: %s", err.Error())
	}

	testModelService := models.ModelService{
		DB: db,
	}

	testModel := Model{
		ModelService: &testModelService,
	}

	testModel.Create(w, r)
	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		t.Fatalf("model Create status=%d, want=200; body=%s", resp.StatusCode, string(body))
	}

	if ct := resp.Header.Get("Content-Type"); ct != "application/json" {
		t.Fatalf("model Create content-type=%s, want=application/json", ct)
	}
}
