package controllers

import (
	"facegram/database"
	"facegram/models"
	"facegram/utils"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

func Register(c *gin.Context) {
	var input struct {
		Username  string `json:"username" binding:"required"`
		FullName  string `json:"full_name" binding:"required"`
		Password  string `json:"password" binding:"required"`
		Bio       string `json:"bio" binding:"required"`
		IsPrivate bool `json:"is_private" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var existing models.User
	if err := database.DB.Where("username = ?", input.Username).First(&existing).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username already exists"})
		return
	}

	hashed, _ := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)

	user := models.User{
		Username:  input.Username,
		FullName:  input.FullName,
		Password:  string(hashed),
		Bio:       input.Bio,
		IsPrivate: input.IsPrivate,
	}

	if err := database.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	token, err := utils.GenerateJWT(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate JWT token"})
	}

	c.JSON(http.StatusOK, gin.H{
		"user": gin.H{
			"id":         user.ID,
			"username":   user.Username,
			"full_name":  user.FullName,
			"is_private": user.IsPrivate,
			"bio":        user.Bio,
		},
		"token": token,
	})
}
