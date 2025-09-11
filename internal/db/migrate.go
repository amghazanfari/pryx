package db

import (
	"fmt"

	"gorm.io/gorm"
	"pryx/internal/models"
)

func AutoMigrateAll(db *gorm.DB) error {
	if err := db.AutoMigrate(
		&models.Model{},
	); err != nil {
		return fmt.Errorf("automigrate: %w", err)
	}
	return nil
}
