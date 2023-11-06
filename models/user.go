package models

import (
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"time"
)

type User struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	Tasks     []Task    `gorm:"foreignkey:UserID"` // This indicates a one-to-many relationship
}

// BeforeCreate will set a UUID rather than numeric ID.
func (base *User) BeforeCreate(tx *gorm.DB) (err error) {
	base.ID = uuid.New()
	currentTime := time.Now()
	base.CreatedAt = currentTime
	base.UpdatedAt = currentTime
	return nil
}

// GORM V2 uses callbacks like BeforeUpdate to handle the update timestamp
func (user *User) BeforeUpdate(tx *gorm.DB) (err error) {
	user.UpdatedAt = time.Now()
	return
}
