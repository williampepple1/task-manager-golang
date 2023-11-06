package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

type Task struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;" json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	UserID      uuid.UUID `gorm:"type:uuid;not null" json:"user_id"` // This is the foreign key
	User        User      `gorm:"foreignkey:UserID"`                 // Use User as a reference for the foreign key
	// UserId      string `json:"user_id"`
}

// BeforeCreate will set a UUID rather than numeric ID.
func (base *Task) BeforeCreate(tx *gorm.DB) (err error) {
	base.ID = uuid.New()
	currentTime := time.Now()
	base.CreatedAt = currentTime
	base.UpdatedAt = currentTime
	return nil
}

// GORM V2 uses callbacks like BeforeUpdate to handle the update timestamp
func (base *Task) BeforeUpdate(tx *gorm.DB) (err error) {
	base.UpdatedAt = time.Now()
	return
}
