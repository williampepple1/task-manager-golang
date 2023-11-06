package handlers

import (
	"fmt"
	"net/http"
	"strings"
	"task-manager/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
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

func CreateTask(db *gorm.DB, secretKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract the Bearer token from the Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			return
		}

		// Split the header to get the token part
		headerParts := strings.Split(authHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header format must be Bearer {token}"})
			return
		}

		// Parse the token
		tokenString := headerParts[1]
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(secretKey), nil
		})

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		// Extract the user ID from the token claims
		userIDStr, ok := claims["Id"].(string)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to extract user ID from token"})
			return
		}
		fmt.Println(userIDStr)
		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse user ID as UUID"})
			return
		}

		// Continue with the usual creation logic, now using the userID...
		var task models.Task
		if err := c.ShouldBindJSON(&task); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Here you would assign the userID to the task object, assuming it has an attribute for it
		task.UserID = userID // This assumes that the Task model has a UserID field of type uuid.UUID

		if err := db.Create(&task).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create task"})
			return
		}

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
