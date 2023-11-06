package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

// TaskStatus is a custom type that represents the status of a task
type TaskStatus string

// Constants for TaskStatus
const (
	StatusCompleted TaskStatus = "completed"
	StatusOngoing   TaskStatus = "ongoing"
	StatusPending   TaskStatus = "pending"
)

type Task struct {
	ID          uuid.UUID  `gorm:"type:uuid;primary_key;" json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	UserID      uuid.UUID  `gorm:"type:uuid;not null" json:"user_id"`                       // This is the foreign key
	User        User       `gorm:"foreignkey:UserID"`                                       // Use User as a reference for the foreign key
	Status      TaskStatus `gorm:"type:varchar(20);not null;default:'to do'" json:"status"` // Enum field
}

// BeforeCreate will set a UUID rather than numeric ID.
func (task *Task) BeforeCreate(tx *gorm.DB) (err error) {
	task.ID = uuid.New()
	currentTime := time.Now()
	task.CreatedAt = currentTime
	task.UpdatedAt = currentTime
	if task.Status == "" {
		task.Status = StatusPending // Set the default status if not provided
	}
	return nil
}

// GORM V2 uses callbacks like BeforeUpdate to handle the update timestamp
func (task *Task) BeforeUpdate(tx *gorm.DB) (err error) {
	task.UpdatedAt = time.Now()
	return
}
