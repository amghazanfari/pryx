package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"fmt"

	"pryx/internal/testutil"
	"pryx/internal/db"
	"pryx/internal/models"
)

func setup(t *testing.T) *Handler {
	t.Helper()
	g := testutil.NewTestDB(t)
	if err := db.AutoMigrateAll(g); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	return New(g)
}

func TestCreateUser(t *testing.T) {
	h := setup(t)
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/admin/users", bytes.NewBufferString(`{"email":"me@ex.com","name":"Me"}`))
	h.CreateUser().ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("want 200, got %d: %s", w.Code, w.Body.String())
	}
	var u models.User
	if err := json.Unmarshal(w.Body.Bytes(), &u); err != nil {
		t.Fatalf("json: %v", err)
	}
	if u.ID == 0 || u.Email != "me@ex.com" {
		t.Fatalf("bad user: %+v", u)
	}
}

func TestCreateAPIKey(t *testing.T) {
	h := setup(t)

	// seed user
	w1 := httptest.NewRecorder()
	r1 := httptest.NewRequest(http.MethodPost, "/admin/users", bytes.NewBufferString(`{"email":"a@b.c","name":"A"}`))
	h.CreateUser().ServeHTTP(w1, r1)
	var u models.User
	if err := json.Unmarshal(w1.Body.Bytes(), &u); err != nil {
		t.Fatalf("seed user: %v", err)
	}

	w := httptest.NewRecorder()
	body := fmt.Sprintf(`{"user_id":%d,"name":"cli","scopes":"completion:invoke"}`, u.ID)
	r := httptest.NewRequest(http.MethodPost, "/admin/keys", bytes.NewBufferString(body))
	h.CreateAPIKey().ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("want 200, got %d: %s", w.Code, w.Body.String())
	}
	var resp struct {
		ID     uint   `json:"id"`
		Prefix string `json:"prefix"`
		Key    string `json:"key"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("json: %v", err)
	}
	if resp.ID == 0 || resp.Prefix == "" || resp.Key == "" {
		t.Fatalf("bad resp: %+v", resp)
	}
}
