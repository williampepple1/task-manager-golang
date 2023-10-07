package handlers

import (
	"task-manager/models"

	"github.com/gin-gonic/gin"
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
