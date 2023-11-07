package handlers

import (
	"fmt"
	"task-manager/models"

	"net/http"

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

func ListUserTasks(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract the user ID from the URL parameter.
		userIDStr := c.Param("userId")

		// Validate that the userID is a valid UUID
		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is not a valid UUID"})
			return
		}

		var tasks []models.Task
		// Find tasks where the UserID matches the provided UUID.
		if err := db.Where("user_id = ?", userID).Find(&tasks).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve tasks for the user"})
			return
		}

		c.JSON(http.StatusOK, tasks)
	}
}
func CreateTask(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var task models.Task
		if err := c.ShouldBindJSON(&task); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		// Retrieve the user ID from the context
		userID, exists := c.Get("userId")
		if !exists {
			c.JSON(401, gin.H{"error": "User ID not found"})
			return
		}

		// Assert that userID is of type string
		userStrId, ok := userID.(string)
		if !ok {
			c.JSON(400, gin.H{"error": "User ID is not of type string"})
			return
		}

		// Parse the user ID from string to uuid.UUID
		id, err := uuid.Parse(userStrId)
		if err != nil {
			c.JSON(500, gin.H{"error": "User ID is not a valid UUID"})
			return
		}

		task.UserID = id // Assign the user's UUID to the task's UserID field
		fmt.Println(task.UserID)
		// Create the task in the DB
		if err := db.Create(&task).Error; err != nil {
			c.JSON(500, gin.H{"error": "Failed to create task"})
			return
		}
		// Now load the user associated with the task and update the task.User
		var user models.User
		if err := db.Where("id = ?", task.UserID).First(&user).Error; err != nil {
			// Handle the error. Perhaps the user does not exist in the DB
			c.JSON(500, gin.H{"error": "Failed to load user data"})
			return
		}

		task.User = user // Assuming your Task struct has a User field to hold this data

		// Now task should have the User data included

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
