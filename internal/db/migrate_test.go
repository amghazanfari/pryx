package db

import (
	"testing"

	"pryx/internal/models"
	"pryx/internal/testutil"
)

func TestAutoMigrateAll(t *testing.T) {
	g := testutil.NewTestDB(t)
	if err := AutoMigrateAll(g); err != nil {
		t.Fatalf("automigrate: %v", err)
	}
	// smoke: create rows
	if err := g.Create(&models.User{Email: "a@b", Name: "A"}).Error; err != nil {
		t.Fatalf("insert user: %v", err)
	}
	if err := g.Create(&models.Model{Name: "X", ModelName: "m", Endpoint: "http://e"}).Error; err != nil {
		t.Fatalf("insert model: %v", err)
	}
}
