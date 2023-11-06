package handlers

import (
	"fmt"
	"reflect"
	"task-manager/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

func ListTasks(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var tasks []models.Task
		if err := db.Find(&tasks).Error; err != nil {
			c.JSON(500, gin.H{"error": "Failed to retrieve tasks"})
			return
		}
		c.JSON(200, tasks)
	}
}

func CreateTask(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var task models.Task
		if err := c.ShouldBindJSON(&task); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		userID , exists := c.Get("userId")
		fmt.Println(userID)
		fmt.Println("type:", reflect.TypeOf(userID))
		if !exists {
			c.JSON(401, gin.H{"error": "User ID not found"})
			return
		}

		id, ok := userID.(uuid.UUID)
		if !ok {
			c.JSON(500, gin.H{"error": "User ID is not a valid UUID"})
			return
		}

		task.UserID = id // Assign the user's UUID to the task's ID

		if err := db.Create(&task).Error; err != nil {
			c.JSON(500, gin.H{"error": "Failed to create task"})
			return
		}

		c.JSON(201, task)
	}
}

func GetTask(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		var task models.Task
		if err := db.Where("id = ?", id).First(&task).Error; err != nil {
			c.JSON(404, gin.H{"error": "Task not found"})
			return
		}

		c.JSON(200, task)
	}
}

func UpdateTask(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var task models.Task

		if err := c.ShouldBindJSON(&task); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		if err := db.Where("id = ?", id).Save(&task).Error; err != nil {
			c.JSON(500, gin.H{"error": "Failed to update task"})
			return
		}

		c.JSON(200, task)
	}
}

func DeleteTask(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		if err := db.Where("id = ?", id).Delete(&models.Task{}).Error; err != nil {
			c.JSON(500, gin.H{"error": "Failed to delete task"})
			return
		}

		c.JSON(200, gin.H{"message": "Task deleted successfully"})
	}
}
