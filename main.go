package main

import (
	"github.com/gin-gonic/gin"
	"task-manager/handlers"
)

func main() {
	r := gin.Default()

	r.GET("/tasks", handlers.ListTasks)
	r.POST("/tasks", handlers.CreateTask)
	r.GET("/tasks/:id", handlers.GetTask)
	r.PUT("/tasks/:id", handlers.UpdateTask)
	r.DELETE("/tasks/:id", handlers.DeleteTask)

	r.Run(":8080")
}
