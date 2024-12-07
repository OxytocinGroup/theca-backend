package domain

import "time"

type Session struct {
	ID        string `gorm:"primaryKey;unique;not null" json:"id"`
	UserID    uint   `json:"user_id"`
	ExpiresAt time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}
