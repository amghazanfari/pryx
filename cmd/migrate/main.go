package main

import (
	"log"
	"pryx/config"
	"pryx/internal/db"
	"pryx/internal/models"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func m0001() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "0001_init_models",
		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&models.Model{})
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropTable(&models.Model{})
		},
	}
}

func main() {
	cfg := config.DBFromEnv()
	conn, err := db.Open(cfg)
	if err != nil { log.Fatal(err) }

	m := gormigrate.New(conn, gormigrate.DefaultOptions, []*gormigrate.Migration{
		m0001(),
	})
	if err := m.Migrate(); err != nil {
		log.Fatal(err)
	}
	log.Println("migrations applied")
}