package handlers

import (
	"task-manager/models"

	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

// Helper function to retrieve and validate user ID and task
func getUserIDAndTask(c *gin.Context, db *gorm.DB, taskId string) (uuid.UUID, models.Task, bool) {
	var task models.Task

	// Retrieve the user ID from the context
	userID, exists := c.Get("userId")
	if !exists {
		c.JSON(401, gin.H{"error": "User ID not found"})
		return uuid.Nil, task, false
	}

	// Assert that userID is of type string
	userStrId, ok := userID.(string)
	if !ok {
		c.JSON(400, gin.H{"error": "User ID is not of type string"})
		return uuid.Nil, task, false
	}

	// Parse the user ID from string to uuid.UUID
	userUUID, err := uuid.Parse(userStrId)
	if err != nil {
		c.JSON(500, gin.H{"error": "User ID is not a valid UUID"})
		return uuid.Nil, task, false
	}

	// Find the task by ID and retrieve it along with its UserID
	if err := db.Where("id = ?", taskId).First(&task).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(404, gin.H{"error": "Task not found"})
		} else {
			c.JSON(500, gin.H{"error": "Failed to retrieve task"})
		}
		return uuid.Nil, task, false
	}

	// Check if the user is authorized to access the task
	if userUUID != task.UserID {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "You are not authorized to access this task"})
		return uuid.Nil, task, false
	}

	return userUUID, task, true
}

func ListTasks(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var tasks []models.Task
		if err := db.Find(&tasks).Error; err != nil {
			c.JSON(500, gin.H{"error": "Failed to retrieve tasks"})
			return
		}
		c.JSON(http.StatusOK, tasks)
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
			c.JSON(http.StatusNotFound, gin.H{"error": "User ID not found"})
			return
		}

		// Assert that userID is of type string
		userStrId, ok := userID.(string)
		if !ok {
			c.JSON(http.StatusNotAcceptable, gin.H{"error": "User ID is not of type string"})
			return
		}

		// Parse the user ID from string to uuid.UUID
		id, err := uuid.Parse(userStrId)
		if err != nil {
			c.JSON(http.StatusNotAcceptable, gin.H{"error": "User ID is not a valid UUID"})
			return
		}

		task.UserID = id // Assign the user's UUID to the task's UserID field

		// Create the task in the DB
		if err := db.Create(&task).Error; err != nil {
			c.JSON(http.StatusNotImplemented, gin.H{"error": "Failed to create task"})
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

		c.JSON(http.StatusCreated, task)
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

		c.JSON(http.StatusOK, task)
	}
}

func UpdateTask(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		// Use the helper function to get the user ID and task
		userUUID, task, ok := getUserIDAndTask(c, db, id)
		if !ok {
			return
		}
		// Now you can compare the UserID from the task with the user's UUID
		if userUUID != task.UserID {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "You are not allowed to update this task"})
			return
		}

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
		id := c.Param("id") // This is the task's ID

		userUUID, task, ok := getUserIDAndTask(c, db, id)
		if !ok {
			return
		}
		// Now you can compare the UserID from the task with the user's UUID
		if userUUID != task.UserID {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "You are not allowed to update this task"})
			return
		}

		// If the user IDs match, delete the task
		if err := db.Delete(&task).Error; err != nil {
			c.JSON(500, gin.H{"error": "Failed to delete task"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Task deleted successfully"})
	}
}
