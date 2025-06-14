package handlers

import (
	"Go_Gin_To-Do_List_API/auth"
	"Go_Gin_To-Do_List_API/database"
	"Go_Gin_To-Do_List_API/models"
	"Go_Gin_To-Do_List_API/utils"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// Register creates a new user.
func Register(c *gin.Context) {
	var input struct {
		Username string `json:"username" binding:"required"`
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	verificationTokenBytes := make([]byte, 16)
	if _, err := rand.Read(verificationTokenBytes); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate verification token"})
		return
	}
	verificationToken := hex.EncodeToString(verificationTokenBytes)

	user := models.User{
		Username:          input.Username,
		Email:             input.Email,
		PasswordHash:      string(hashedPassword),
		VerificationToken: verificationToken,
	}

	result := database.DB.Create(&user)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user", "details": result.Error.Error()})
		return
	}

	appURL := os.Getenv("APP_URL")
	if appURL == "" {
		appURL = "http://localhost:8080"
	}
	verificationLink := fmt.Sprintf("%s/api/v1/verify-email?token=%s", appURL, verificationToken)

	log.Printf("Verification link for user %s: %s", user.Email, verificationLink)

	emailSubject := "Verify Your Email for To-Do List API"
	emailBody := fmt.Sprintf("Hi %s,\n\nPlease verify your email by clicking on the following link:\n%s", user.Username, verificationLink)
	if err := utils.SendEmail(user.Email, emailSubject, emailBody); err != nil {
		log.Printf("Failed to send verification email to %s: %v", user.Email, err)
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully. Please check your email to verify your account."})
}

// VerifyEmail confirms a user's email address.
func VerifyEmail(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Verification token is required"})
		return
	}

	var user models.User
	if err := database.DB.Where("verification_token = ?", token).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid verification token"})
		return
	}

	user.IsVerified = true
	user.VerificationToken = "" // Invalidate the token
	database.DB.Save(&user)

	// Return a JSON response instead of rendering HTML
	c.JSON(http.StatusOK, gin.H{"message": "Email verified successfully. You can now log in."})
}

// Other handlers (Login, ForgotPassword, ResetPassword) remain the same...

// Login authenticates a user and returns a JWT.
func Login(c *gin.Context) {
	var input struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := database.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	if !user.IsVerified {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Please verify your email before logging in."})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	tokenString, err := auth.GenerateJWT(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}

// ForgotPassword initiates the password reset process.
func ForgotPassword(c *gin.Context) {
	var input struct {
		Email string `json:"email" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resetTokenBytes := make([]byte, 16)
	if _, err := rand.Read(resetTokenBytes); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate reset token"})
		return
	}
	resetToken := hex.EncodeToString(resetTokenBytes)

	var user models.User
	if err := database.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"message": "If an account with that email exists, a password reset link has been sent."})
		return
	}

	user.ResetToken = resetToken
	user.ResetTokenExp = time.Now().Add(time.Hour * 1) // Token valid for 1 hour
	database.DB.Save(&user)

	log.Printf("Password reset token for user %s: %s", user.Email, resetToken)

	c.JSON(http.StatusOK, gin.H{"message": "If an account with that email exists, a password reset link has been sent."})
}

// ResetPassword sets a new password for the user.
func ResetPassword(c *gin.Context) {
	var input struct {
		Token       string `json:"token" binding:"required"`
		NewPassword string `json:"new_password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := database.DB.Where("reset_token = ? AND reset_token_exp > ?", input.Token, time.Now()).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid or expired password reset token"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash new password"})
		return
	}

	user.PasswordHash = string(hashedPassword)
	user.ResetToken = ""
	database.DB.Save(&user)

	c.JSON(http.StatusOK, gin.H{"message": "Password has been reset successfully."})
}
