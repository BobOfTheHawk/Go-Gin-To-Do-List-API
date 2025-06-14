package models

import (
	"time"

	"gorm.io/gorm"
)

// User represents a user in the database
type User struct {
	gorm.Model
	Username          string `gorm:"uniqueIndex;not null"`
	Email             string `gorm:"uniqueIndex;not null"`
	PasswordHash      string `gorm:"not null"`
	IsVerified        bool   `gorm:"default:false"`
	VerificationToken string
	ResetToken        string
	ResetTokenExp     time.Time
	Tasks             []Task `gorm:"foreignKey:UserID"`
}

// Task represents a to-do item, belonging to a user
type Task struct {
	gorm.Model
	Title       string `gorm:"not null"`
	Description string
	Status      string `gorm:"default:'pending'"`
	UserID      uint   `gorm:"not null"`
}
