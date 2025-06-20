package controllers

import (
	"errors"
	"facegram/database"
	"facegram/models"
	"fmt"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PostInput struct {
	Caption string `json:"caption" form:"caption" binding:"required"`
}

func CreatePost(c *gin.Context) {
	userIDRaw, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found"})
		return
	}

	userID := userIDRaw.(uint)

	var input PostInput

	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid field", "error": err.Error()})
		return
	}

	post := models.Post{
		Caption: input.Caption,
		UserID:  userID,
	}

	if err := database.DB.Create(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid multipart form data"})
		return
	}

	files := form.File["attachments"]
	if len(files) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Attachments is required"})
		return
	}
	for _, file := range files {
		if !isValidImage(file) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": fmt.Sprintf("Invalid file type: %s", file.Filename),
			})
			return
		}

		savePath := filepath.Join("images", file.Filename)
		err := c.SaveUploadedFile(file, savePath)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
			return
		}
		attachments := models.Attachment{
			StoragePath: savePath,
			PostID:      post.ID,
		}

		if err := database.DB.Create(&attachments).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Create post success"})
}

func isValidImage(file *multipart.FileHeader) bool {
	ext := strings.ToLower(filepath.Ext(file.Filename))
	return allowedExts[ext]
}

var allowedExts = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".webp": true,
	".gif":  true,
}

func DeletePost(c *gin.Context) {
	id := c.Param("id")
	userId, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found"})
		return
	}

	var post models.Post

	if err := database.DB.Preload("Attachments").First(&post, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusForbidden, gin.H{"message": "Post not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := database.DB.Where("post_id", post.ID).Delete(&models.Attachment{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if post.UserID != userId {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthenticated"})
		return
	}

	if err := database.DB.Delete(&post, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed delete post"})
		return
	}

	c.Status(http.StatusNoContent)
}

func GetPost(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "0"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))

	if page < 0 {
		page = 0
	}

	if size < 10 || size > 100 {
		size = 10
	}

	var posts []models.Post

	offset := page * size
	if err := database.DB.Offset(offset).Limit(size).Preload("User").Preload("Attachments").First(&posts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"page":  page,
		"size":  size,
		"posts": posts,
	})
}
