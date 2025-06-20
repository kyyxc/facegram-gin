package controllers

import (
	"facegram/database"
	"facegram/models"
	"facegram/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RegisterInput struct {
	Username  string `json:"username" binding:"required"`
	FullName  string `json:"full_name" binding:"required"`
	Password  string `json:"password" binding:"required"`
	Bio       string `json:"bio" binding:"required"`
	IsPrivate bool   `json:"is_private" binding:"required"`
}

type LoginInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func Register(c *gin.Context) {
	var input RegisterInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid field", "error": err.Error()})
		return
	}

	user := models.User{
		Username:  input.Username,
		FullName:  input.FullName,
		Password:  input.Password,
		Bio:       input.Bio,
		IsPrivate: input.IsPrivate,
	}

	var existing models.User
	if err := database.DB.Where("username = ?", input.Username).First(&existing).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username already exists"})
		return
	}

	if err := user.HashPassword(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Could not generate hash password"})
	}

	if err := database.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	token, err := utils.GenerateJWT(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate JWT token"})
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Register success",
		"token":   token,
		"user": gin.H{
			"id":         user.ID,
			"full_name":  user.FullName,
			"username":   user.Username,
			"bio":        user.Bio,
			"is_private": user.IsPrivate,
		},
	})
}

func Login(c *gin.Context) {
	var input LoginInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid field", "error": err.Error()})
	}

	var user models.User
	if err := database.DB.Where("username = ?", input.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Wrong username or password"})
		return
	}

	if !user.CheckPassword(input.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Wrong username or password"})
		return
	}

	token, err := utils.GenerateJWT(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate JWT token"})
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login success",
		"token":   token,
		"user": gin.H{
			"id":         user.ID,
			"full_name":  user.FullName,
			"username":   user.Username,
			"bio":        user.Bio,
			"is_private": user.IsPrivate,
		},
	})
}
