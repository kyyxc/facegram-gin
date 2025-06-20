package controllers

import (
	"facegram/database"
	"facegram/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Follow(c *gin.Context) {
	username := c.Param("username")
	userIDRaw, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found"})
		return
	}

	var user models.User

	if err := database.DB.Where("username", username).First(&user).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "User not found"})
		return
	}

	var userLogged models.User
	if err := database.DB.First(&userLogged, userIDRaw).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "User logged not found"})
		return
	}

	if username == userLogged.Username {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "You are not allowed to follow yourself"})
		return
	}

	var follow models.Follow
	if err := database.DB.Where("following_id = ?", userLogged.ID).
		Where("follower_id = ?", user.ID).
		First(&follow).Error; err == nil {
		status := "requested"
		if follow.IsAccepted {
			status = "following"
		}
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": "You are already followed",
			"status":  status,
		})
		return
	}

	status := "following"
	if user.IsPrivate {
		status = "requested"
	}

	newFollow := models.Follow{
		FollowerID:  user.ID,
		FollowingID: userLogged.ID,
		IsAccepted:  !user.IsPrivate,
	}

	if err := database.DB.Create(&newFollow).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to follow"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Follow success",
		"status":  status,
	})
}

func Unfollow(c *gin.Context) {
	username := c.Param("username")
	userIDRaw, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found"})
		return
	}

	var user models.User

	if err := database.DB.Where("username", username).First(&user).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "User not found"})
		return
	}

	var userLogged models.User
	if err := database.DB.First(&userLogged, userIDRaw).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "User logged not found"})
		return
	}

	var follow models.Follow
	if err := database.DB.Where("follower_id = ?", user.ID).
		Where("following_id", userLogged.ID).
		First(&follow).Error; err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "You are not following the user"})
		return
	}

	if err := database.DB.Delete(&follow).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unfollow"})
		return
	}

	c.Status(204)
}
