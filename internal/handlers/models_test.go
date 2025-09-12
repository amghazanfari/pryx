package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"pryx/internal/db"
	"pryx/internal/models"
	"pryx/internal/testutil"
)

func TestAddModelHandler_Happy(t *testing.T) {
	g := testutil.NewTestDB(t)
	if err := db.AutoMigrateAll(g); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	h := New(g)

	reqBody := `{"name":"OpenRouter GPT-4o Mini","model_name":"gpt-4o-mini","endpoint":"https://openrouter.ai/api/v1","api_key":""}`
	r := httptest.NewRequest(http.MethodPost, "/v1/models", bytes.NewBufferString(reqBody))
	w := httptest.NewRecorder()

	h.AddModelHandler().ServeHTTP(w, r)

	if w.Code != http.StatusCreated {
		t.Fatalf("want 201, got %d (%s)", w.Code, w.Body.String())
	}

	var m models.Model
	if err := json.Unmarshal(w.Body.Bytes(), &m); err != nil {
		t.Fatalf("json: %v", err)
	}
	if m.ID == 0 || m.Name == "" || m.ModelName == "" || m.Endpoint == "" {
		t.Fatalf("bad model: %+v", m)
	}
}

func TestAddModelHandler_ValidationAndMethod(t *testing.T) {
	g := testutil.NewTestDB(t)
	if err := db.AutoMigrateAll(g); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	h := New(g)

	// Wrong method
	r1 := httptest.NewRequest(http.MethodGet, "/v1/models", nil)
	w1 := httptest.NewRecorder()
	h.AddModelHandler().ServeHTTP(w1, r1)
	if w1.Code != http.StatusMethodNotAllowed {
		t.Fatalf("want 405, got %d", w1.Code)
	}

	// Bad payload
	r2 := httptest.NewRequest(http.MethodPost, "/v1/models", bytes.NewBufferString(`{"name":""}`))
	w2 := httptest.NewRecorder()
	h.AddModelHandler().ServeHTTP(w2, r2)
	if w2.Code != http.StatusBadRequest {
		t.Fatalf("want 400, got %d", w2.Code)
	}
}
