package db

import (
	"fmt"

	"gorm.io/gorm"
	"pryx/internal/models"
)

func AutoMigrateAll(db *gorm.DB) error {
	if err := db.AutoMigrate(
		&models.Model{},
		&models.User{},
		&models.APIKey{},
		&models.ModelAccess{},
	); err != nil {
		return fmt.Errorf("automigrate: %w", err)
	}
	return nil
}
