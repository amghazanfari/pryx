package auth

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
	"fmt"

	"gorm.io/gorm"

	"pryx/internal/models"
	"pryx/internal/db"
	"pryx/internal/testutil"
)

func migrateAll(t *testing.T, g *gorm.DB) {
	t.Helper()
	if err := db.AutoMigrateAll(g); err != nil {
		t.Fatalf("automigrate: %v", err)
	}
}

func seedKey(t *testing.T, g *gorm.DB, scopes string, revoked bool) (plain string) {
	t.Helper()
	u := models.User{Email: fmt.Sprintf("%s-%d@example.com", t.Name(), time.Now().UnixNano()), Name: "X", IsActive: true}
	if err := g.Create(&u).Error; err != nil {
		t.Fatalf("create user: %v", err)
	}
	plain, prefix, hash, err := GenerateAPIKey()
	if err != nil {
		t.Fatalf("gen key: %v", err)
	}
	k := models.APIKey{
		UserID: u.ID, Name: "test", Prefix: prefix, Hash: hash,
		Scopes: strings.TrimSpace(scopes), Revoked: revoked,
	}
	if err := g.Create(&k).Error; err != nil {
		t.Fatalf("create key: %v", err)
	}
	return plain
}

func TestMiddleware_Happy(t *testing.T) {
	g := testutil.NewTestDB(t)
	migrateAll(t, g)
	plain := seedKey(t, g, "completion:invoke", false)

	guard := Middleware(g, "completion:invoke")
	ok := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, ok := From(r.Context()); !ok {
			t.Fatalf("no auth context in request")
		}
		w.WriteHeader(200)
	})

	req := httptest.NewRequest("GET", "/x", nil)
	req.Header.Set("Authorization", "Bearer "+plain)
	rr := httptest.NewRecorder()
	guard(ok).ServeHTTP(rr, req)

	if rr.Code != 200 {
		t.Fatalf("want 200, got %d body=%s", rr.Code, rr.Body.String())
	}

	// last_used_at should be set
	var ak models.APIKey
	if err := g.First(&ak).Error; err != nil {
		t.Fatalf("read key: %v", err)
	}
	if ak.LastUsedAt == nil || time.Since(*ak.LastUsedAt) > time.Minute {
		t.Fatalf("last_used_at not updated")
	}
}

func TestMiddleware_ForbiddenMissingScope(t *testing.T) {
	g := testutil.NewTestDB(t)
	migrateAll(t, g)
	plain := seedKey(t, g, "model:write", false)

	guard := Middleware(g, "completion:invoke")
	req := httptest.NewRequest("GET", "/x", nil)
	req.Header.Set("Authorization", "Bearer "+plain)
	rr := httptest.NewRecorder()
	guard(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})).ServeHTTP(rr, req)

	if rr.Code != http.StatusForbidden {
		t.Fatalf("want 403, got %d body=%s", rr.Code, rr.Body.String())
	}
}

func TestMiddleware_RevokedOrBad(t *testing.T) {
	g := testutil.NewTestDB(t)
	migrateAll(t, g)
	plain := seedKey(t, g, "completion:invoke", true) 

	for _, token := range []string{
		"", "Bearer short", "Bearer " + plain, 
	} {
		req := httptest.NewRequest("GET", "/x", nil)
		if token != "" {
			req.Header.Set("Authorization", token)
		}
		rr := httptest.NewRecorder()
		Middleware(g, "completion:invoke")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})).
			ServeHTTP(rr, req)

		if rr.Code != http.StatusUnauthorized {
			t.Fatalf("want 401, got %d token=%q body=%s", rr.Code, token, rr.Body.String())
		}
	}
}
