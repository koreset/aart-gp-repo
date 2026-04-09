package models

import (
	"time"
)

type SystemLock struct {
	Name      string `gorm:"primaryKey"`
	Owner     string
	ExpiresAt time.Time
}
