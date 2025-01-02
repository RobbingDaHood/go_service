package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"todo-service/internal"
	"todo-service/internal/model"
)

func GetTodos(c *gin.Context) {
	var todos []model.Todo
	internal.DB.Find(&todos)
	c.JSON(http.StatusOK, todos)
}

func GetTodoByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var todo model.Todo
	if err := internal.DB.First(&todo, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
		return
	}
	c.JSON(http.StatusOK, todo)
}

func CreateTodo(c *gin.Context) {
	var todo model.Todo
	if err := c.ShouldBindJSON(&todo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	internal.DB.Create(&todo)
	c.JSON(http.StatusCreated, todo)
}

func UpdateTodoByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var todo model.Todo
	if err := internal.DB.First(&todo, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
		return
	}

	var input model.Todo
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	todo.Title = input.Title
	todo.Completed = input.Completed
	internal.DB.Save(&todo)
	c.JSON(http.StatusOK, todo)
}

func DeleteTodoByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := internal.DB.Delete(&model.Todo{}, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Todo deleted"})
}
