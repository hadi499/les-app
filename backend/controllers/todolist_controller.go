package controllers

import (
	"backend/database"
	"backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetTodoLists - Get all todo lists for the logged-in user
func GetTodoLists(c *gin.Context) {
	userID := c.GetUint("userID")

	var lists []models.TodoList
	if err := database.DB.Preload("Items").Where("user_id = ?", userID).Find(&lists).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch todo lists"})
		return
	}

	c.JSON(http.StatusOK, lists)
}

// GetTodoList - Get a single todo list by id
func GetTodoList(c *gin.Context) {
	userID := c.GetUint("userID")
	id := c.Param("id")

	var list models.TodoList
	if err := database.DB.Preload("Items").Where("id = ? AND user_id = ?", id, userID).First(&list).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Todo list not found"})
		return
	}

	c.JSON(http.StatusOK, list)
}

// CreateTodoList - Create a new todo list title
func CreateTodoList(c *gin.Context) {
	userID := c.GetUint("userID")

	var input struct {
		Title string `json:"title" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	list := models.TodoList{
		UserID: userID,
		Title:  input.Title,
	}

	if err := database.DB.Create(&list).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create todo list"})
		return
	}

	c.JSON(http.StatusOK, list)
}

// UpdateTodoList - Update a todo list title
func UpdateTodoList(c *gin.Context) {
	userID := c.GetUint("userID")
	id := c.Param("id")

	var list models.TodoList
	if err := database.DB.Where("id = ? AND user_id = ?", id, userID).First(&list).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Todo list not found"})
		return
	}

	var input struct {
		Title string `json:"title" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	list.Title = input.Title
	if err := database.DB.Save(&list).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update todo list"})
		return
	}

	c.JSON(http.StatusOK, list)
}

// DeleteTodoList - Delete a todo list and its items (cascade)
func DeleteTodoList(c *gin.Context) {
	userID := c.GetUint("userID")
	id := c.Param("id")

	var list models.TodoList
	if err := database.DB.Where("id = ? AND user_id = ?", id, userID).First(&list).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Todo list not found"})
		return
	}

	if err := database.DB.Delete(&list).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete todo list"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Todo list deleted"})
}

// CreateTodoItem - Add a new item to a todolist
func CreateTodoItem(c *gin.Context) {
	userID := c.GetUint("userID")
	listID := c.Param("id")

	// Verify the list belongs to the user
	var list models.TodoList
	if err := database.DB.Where("id = ? AND user_id = ?", listID, userID).First(&list).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Todo list not found"})
		return
	}

	var input struct {
		Text string `json:"text" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	item := models.TodoItem{
		TodoListID: list.ID,
		Text:       input.Text,
		Completed:  false,
	}

	if err := database.DB.Create(&item).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create todo item"})
		return
	}

	c.JSON(http.StatusOK, item)
}

// ToggleTodoItem - Toggle completion status
func ToggleTodoItem(c *gin.Context) {
	userID := c.GetUint("userID")
	listID := c.Param("id")
	itemID := c.Param("item_id")

	// Verify list ownership
	var list models.TodoList
	if err := database.DB.Where("id = ? AND user_id = ?", listID, userID).First(&list).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Todo list not found"})
		return
	}

	var item models.TodoItem
	if err := database.DB.Where("id = ? AND todo_list_id = ?", itemID, list.ID).First(&item).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Todo item not found"})
		return
	}

	var input struct {
		Completed bool `json:"completed"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	item.Completed = input.Completed
	if err := database.DB.Save(&item).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update item"})
		return
	}

	c.JSON(http.StatusOK, item)
}

// DeleteTodoItem - Delete an item
func DeleteTodoItem(c *gin.Context) {
	userID := c.GetUint("userID")
	listID := c.Param("id")
	itemID := c.Param("item_id")

	// Verify list ownership
	var list models.TodoList
	if err := database.DB.Where("id = ? AND user_id = ?", listID, userID).First(&list).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Todo list not found"})
		return
	}

	var item models.TodoItem
	if err := database.DB.Where("id = ? AND todo_list_id = ?", itemID, list.ID).First(&item).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Todo item not found"})
		return
	}

	if err := database.DB.Delete(&item).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete item"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Item deleted"})
}
