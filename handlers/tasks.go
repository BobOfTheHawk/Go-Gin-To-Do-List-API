package handlers

import (
	"Go_Gin_To-Do_List_API/database"
	"Go_Gin_To-Do_List_API/models"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CreateTask adds a new task for the authenticated user.
func CreateTask(c *gin.Context) {
	var input struct {
		Title       string `json:"title" binding:"required"`
		Description string `json:"description"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := c.Get("userID")

	task := models.Task{
		Title:       input.Title,
		Description: input.Description,
		UserID:      userID.(uint),
		Status:      "pending",
	}

	database.DB.Create(&task)
	c.JSON(http.StatusCreated, task)
}

// GetTasks retrieves all tasks for the authenticated user.
func GetTasks(c *gin.Context) {
	userID, _ := c.Get("userID")
	var tasks []models.Task
	database.DB.Where("user_id = ?", userID).Find(&tasks)
	c.JSON(http.StatusOK, tasks)
}

// GetTask retrieves a single task by ID.
func GetTask(c *gin.Context) {
	userID, _ := c.Get("userID")
	taskID := c.Param("id")

	var task models.Task
	if err := database.DB.Where("id = ? AND user_id = ?", taskID, userID).First(&task).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	c.JSON(http.StatusOK, task)
}

// UpdateTask modifies an existing task.
func UpdateTask(c *gin.Context) {
	userID, _ := c.Get("userID")
	taskID := c.Param("id")

	var task models.Task
	if err := database.DB.Where("id = ? AND user_id = ?", taskID, userID).First(&task).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	var input struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Status      string `json:"status"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Use map to only update non-empty fields
	updates := make(map[string]interface{})
	if input.Title != "" {
		updates["title"] = input.Title
	}
	if input.Description != "" {
		updates["description"] = input.Description
	}
	if input.Status != "" {
		updates["status"] = input.Status
	}

	database.DB.Model(&task).Updates(updates)
	c.JSON(http.StatusOK, task)
}

// DeleteTask removes a task.
func DeleteTask(c *gin.Context) {
	userID, _ := c.Get("userID")
	taskID := c.Param("id")

	var task models.Task
	if err := database.DB.Where("id = ? AND user_id = ?", taskID, userID).First(&task).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	database.DB.Delete(&task)
	c.JSON(http.StatusOK, gin.H{"message": "Task deleted successfully"})
}
