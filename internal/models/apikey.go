package models

import "time"

type APIKey struct {
	ID         uint       `json:"id" gorm:"primaryKey"`
	UserID     uint       `json:"user_id" gorm:"index;not null"`
	Name       string     `json:"name" gorm:"type:text"`
	Prefix     string     `json:"prefix" gorm:"type:char(8);index"`
	Hash       string     `json:"-" gorm:"type:char(64);uniqueIndex"`
	Scopes     string     `json:"scopes" gorm:"type:text"`
	Revoked    bool       `json:"revoked" gorm:"default:false"`
	LastUsedAt *time.Time `json:"last_used_at"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}
