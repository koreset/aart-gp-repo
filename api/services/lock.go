package services

import (
	"api/models"
	"errors"
	"os"
	"time"

	"gorm.io/gorm"
)

// AcquireLock attempts to acquire a named lock in the database.
// It returns true if the lock was acquired, false otherwise.
func AcquireLock(name string, duration time.Duration) bool {
	owner, _ := os.Hostname()
	now := time.Now()
	expiresAt := now.Add(duration)

	var lock models.SystemLock
	err := DB.Where("name = ?", name).First(&lock).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		// Try to create the lock
		newLock := models.SystemLock{
			Name:      name,
			Owner:     owner,
			ExpiresAt: expiresAt,
		}
		if err := DB.Create(&newLock).Error; err == nil {
			return true
		}
		// If create failed, someone else might have created it simultaneously
		// Fall through to check if it's expired
		err = DB.Where("name = ?", name).First(&lock).Error
	}

	if err != nil {
		return false
	}

	// If the lock is held by someone else but has expired, we can take it
	if lock.ExpiresAt.Before(now) {
		result := DB.Model(&models.SystemLock{}).
			Where("name = ? AND expires_at = ?", name, lock.ExpiresAt).
			Updates(models.SystemLock{
				Owner:     owner,
				ExpiresAt: expiresAt,
			})
		return result.RowsAffected > 0
	}

	return false
}

// ReleaseLock releases a named lock if owned by the caller.
func ReleaseLock(name string) {
	owner, _ := os.Hostname()
	DB.Where("name = ? AND owner = ?", name, owner).Delete(&models.SystemLock{})
}
