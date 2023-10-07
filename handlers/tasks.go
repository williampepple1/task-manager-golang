package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"task-manager/models"
)

func ListTasks(c *gin.Context) {
	models.TasksMu.RLock()
	defer models.TasksMu.RUnlock()
	c.JSON(200, models.Tasks)
}

// ... other handler functions (e.g., CreateTask, GetTask, etc.) remain mostly unchanged

func createTask(c *gin.Context) {
	var task Task
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	tasksMu.Lock()
	task.ID = nextID
	nextID++
	tasks = append(tasks, &task)
	tasksMu.Unlock()

	c.JSON(201, task)
}

func getTask(c *gin.Context) {
	id := c.Param("id")

	tasksMu.RLock()
	defer tasksMu.RUnlock()

	for _, t := range tasks {
		if fmt.Sprintf("%d", t.ID) == id {
			c.JSON(200, t)
			return
		}
	}
	c.JSON(404, gin.H{"error": "Task not found"})
}

func updateTask(c *gin.Context) {
	id := c.Param("id")
	var task Task

	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	tasksMu.Lock()
	defer tasksMu.Unlock()

	for _, t := range tasks {
		if fmt.Sprintf("%d", t.ID) == id {
			t.Title = task.Title
			c.JSON(200, t)
			return
		}
	}

	c.JSON(404, gin.H{"error": "Task not found"})
}

func deleteTask(c *gin.Context) {
	id := c.Param("id")

	tasksMu.Lock()
	defer tasksMu.Unlock()

	for i, t := range tasks {
		if fmt.Sprintf("%d", t.ID) == id {
			tasks = append(tasks[:i], tasks[i+1:]...)
			c.JSON(200, gin.H{"message": "Task deleted"})
			return
		}
	}

	c.JSON(404, gin.H{"error": "Task not found"})
}
