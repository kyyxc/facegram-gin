package controllers

import (
	"errors"
	"facegram/database"
	"facegram/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusForbidden, gin.H{"error": "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var userLogged models.User
	if err := database.DB.First(&userLogged, userIDRaw).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusForbidden, gin.H{"error": "User logged not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if username == userLogged.Username {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "You are not allowed to follow yourself"})
		return
	}

	var follow models.Follow
	if err := database.DB.Where("following_id = ?", user.ID).
		Where("follower_id = ?", userLogged.ID).
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
		FollowerID:  userLogged.ID,
		FollowingID: user.ID,
		IsAccepted:  !user.IsPrivate,
	}

	if err := database.DB.Create(&newFollow).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusForbidden, gin.H{"error": "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var userLogged models.User
	if err := database.DB.First(&userLogged, userIDRaw).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusForbidden, gin.H{"error": "User logged not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var follow models.Follow
	if err := database.DB.Where("follower_id = ?", userLogged.ID).
		Where("following_id", user.ID).
		First(&follow).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "You are not following the user"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := database.DB.Delete(&follow).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(204)
}

type FollowResponse struct {
	ID uint `json:"id"`

	Username  string `json:"username"`
	FullName  string `json:"full_name"`
	Bio       string `json:"bio"`
	IsPrivate bool   `json:"is_private"`

	CreatedAt   time.Time `json:"created_at"`
	IsRequested bool      `json:"is_requested"`
}

func GetFollowing(c *gin.Context) {
	username := c.Param("username")
	var user models.User

	if err := database.DB.Where("username", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusForbidden, gin.H{"error": "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var followingIDs []uint

	if err := database.DB.Where("follower_id", user.ID).
		Model(&models.Follow{}).
		Pluck("id", &followingIDs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var followings []models.Follow

	if err := database.DB.Where("id IN ?", followingIDs).
		Preload("Following").
		Find(&followings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var followingsUser []FollowResponse

	for _, follow := range followings {
		user := follow.Following
		followingsUser = append(followingsUser, FollowResponse{
			ID:          user.ID,
			Username:    user.Username,
			FullName:    user.FullName,
			Bio:         user.Bio,
			IsPrivate:   user.IsPrivate,
			CreatedAt:   user.CreatedAt,
			IsRequested: !follow.IsAccepted,
		})
	}

	c.JSON(http.StatusOK, gin.H{"following": followingsUser})
}

func Accept(c *gin.Context) {
	username := c.Param("username")
	userIDRaw, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User id not found"})
	}

	var user models.User
	var userLogged models.User

	if err := database.DB.Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusForbidden, gin.H{"error": "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := database.DB.First(&userLogged, userIDRaw).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusForbidden, gin.H{"error": "User logged not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var follow models.Follow
	if err := database.DB.Where("follower_id", user.ID).Where("following_id", userLogged.ID).First(&follow).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "The user is not following you"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if follow.IsAccepted {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Follow request is already accepted"})
		return
	}

	if err := database.DB.Model(&follow).Update("is_accepted", true).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to accept follower"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"error": "Follow request accepted"})
}

func GetFollower(c *gin.Context) {
	username := c.Param("username")
	var user models.User

	if err := database.DB.Where("username", username).First(&user).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "User not found"})
		return
	}

	var followerIDs []uint

	if err := database.DB.Where("following_id", user.ID).Model(&models.Follow{}).Pluck("id", &followerIDs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var followers []models.Follow

	if err := database.DB.Where("id IN ?", followerIDs).Preload("Follower").Find(&followers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var followerUsers []FollowResponse

	for _, follow := range followers {
		user := follow.Follower
		followerUsers = append(followerUsers, FollowResponse{
			ID:          user.ID,
			Username:    user.Username,
			FullName:    user.FullName,
			Bio:         user.Bio,
			IsPrivate:   user.IsPrivate,
			CreatedAt:   user.CreatedAt,
			IsRequested: !follow.IsAccepted,
		})
	}

	c.JSON(http.StatusOK, gin.H{"following": followerUsers})
}
