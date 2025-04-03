// controllers/auth.go
package controllers

import (
	"cal-blog-service/auth"
	"cal-blog-service/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// AuthController handles user authentication
type AuthController struct {
	DB *gorm.DB
}

// Register handles user registration
func (authCtrl *AuthController) Register(context *gin.Context) {
	var user models.User

	if err := context.ShouldBindJSON(&user); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if username already exists
	var existingUser models.User
	if result := authCtrl.DB.Where("username = ?", user.Username).First(&existingUser); result.RowsAffected > 0 {
		context.JSON(http.StatusConflict, gin.H{"error": "Username already exists"})
		return
	}

	// Hash password before saving
	if err := user.HashPassword(user.Password); err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// Set default role
	if user.Role == "" {
		user.Role = "user"
	}

	// Create user in database
	if result := authCtrl.DB.Create(&user); result.Error != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	context.JSON(http.StatusCreated, gin.H{
		"id":       user.ID,
		"username": user.Username,
		"role":     user.Role,
	})
}

// Login authenticates a user and returns a JWT
func (authCtrl *AuthController) Login(context *gin.Context) {
	var credentials struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := context.ShouldBindJSON(&credentials); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Find user by username
	var user models.User
	if result := authCtrl.DB.Where("username = ?", credentials.Username).First(&user); result.Error != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username"})
		return
	}

	// Check password
	if !user.CheckPassword(credentials.Password) {
		context.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
		return
	}

	// Generate JWT token
	token, err := auth.GenerateToken(&user)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"token": token,
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"role":     user.Role,
		},
	})
}
