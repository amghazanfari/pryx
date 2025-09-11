package models

import "time"

type Model struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" gorm:"type:text;not null;index"`
	ModelName string    `json:"model_name" gorm:"type:text;not null"`
	Endpoint  string    `json:"endpoint" gorm:"type:text;not null"`
	APIKey    string    `json:"api_key" gorm:"type:text"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
