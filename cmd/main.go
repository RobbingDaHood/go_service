package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"todo-service/internal"
	"todo-service/internal/handler"
)

func main() {
	setupDatabase()
	err := setupServer().Run(":8080")
	if err != nil {
		fmt.Printf("Error starting server: %s\n", err)
	}
}

func setupDatabase() {
	internal.Init()
}

func setupServer() *gin.Engine {
	r := gin.Default()

	r.GET("/todos", handler.GetTodos)
	r.GET("/todos/:id", handler.GetTodoByID)
	r.POST("/todos", handler.CreateTodo)
	r.PUT("/todos/:id", handler.UpdateTodoByID)
	r.DELETE("/todos/:id", handler.DeleteTodoByID)

	return r // Start the server on port 8080
}
