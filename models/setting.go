package models

import "time"

type Setting struct {
	Key         string    `gorm:"primaryKey" json:"key"`
	Value       string    `json:"value"`
	LastValue   string    `json:"last_value"`
	Description string    `gorm:"type:text" json:"description"`
	UpdatedAt   time.Time `json:"updated_at"`
}
