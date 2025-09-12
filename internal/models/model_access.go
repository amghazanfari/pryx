package models

type ModelAccess struct {
	ID       uint `gorm:"primaryKey"`
	APIKeyID uint `gorm:"index;not null"`
	ModelID  uint `gorm:"index;not null"`
}
