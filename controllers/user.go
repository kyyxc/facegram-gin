package controllers

import (
	"facegram/database"
	"facegram/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func GetUser(c *gin.Context) {
	userIDRaw, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found"})
		return
	}

	var followingsIds []uint
	if err := database.DB.Where("follower_id", userIDRaw).
		Model(models.Follow{}).
		Pluck("following_id", &followingsIds).
		Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	followingsIds = append(followingsIds, userIDRaw.(uint))

	var users []models.User
	if err := database.DB.Where("id NOT IN ?", followingsIds).Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"users": users})
}

type ShowUserResponse struct {
	ID              uint          `json:"id"`
	Username        string        `json:"username"`
	FullName        string        `json:"full_name"`
	Bio             string        `json:"bio"`
	IsPrivate       bool          `json:"is_private"`
	CreatedAt       time.Time     `json:"created_at"`
	UpdatedAt       time.Time     `json:"updated_at"`
	IsYourAccount   bool          `json:"is_your_account"`
	FollowingStatus string        `json:"following_status"`
	PostsCount      int           `json:"posts_count"`
	FollowerCount   int           `json:"followers_count"`
	FollowingCount  int           `json:"following_count"`
	Post            []models.Post `json:"posts"`
}

func ShowUser(c *gin.Context) {
	username := c.Param("username")
	var user models.User
	if err := database.DB.Where("username = ?", username).
		Preload("Posts.Attachments").
		Preload("Followers").
		Preload("Followings").
		First(&user).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{"message": "User not found"})
		return
	}

	userLoggedID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "User logged id not found"})
		return
	}
	isYourAccount := userLoggedID == user.ID
	var followingStatus string

	var follow models.Follow
	if err := database.DB.Where("following_id", user.ID).
		Where("follower_id", userLoggedID).
		First(&follow).
		Error; err != nil {
		followingStatus = "not-following"
	} else {
		if follow.IsAccepted {
			followingStatus = "following"
		} else {
			followingStatus = "requested"
		}
	}
	shouldHidden := user.IsPrivate && !isYourAccount && followingStatus != "following"
	var posts []models.Post
	if !shouldHidden && user.Posts != nil {
		posts = user.Posts
	}

	res := ShowUserResponse{
		ID:              user.ID,
		Username:        user.Username,
		FullName:        user.FullName,
		Bio:             user.Bio,
		IsPrivate:       user.IsPrivate,
		CreatedAt:       user.CreatedAt,
		UpdatedAt:       user.UpdatedAt,
		IsYourAccount:   isYourAccount,
		FollowingStatus: followingStatus,
		PostsCount:      len(user.Posts),
		FollowerCount:   len(user.Followers),
		FollowingCount:  len(user.Followings),
		Post:            posts,
	}
	c.JSON(http.StatusOK, res)
}
