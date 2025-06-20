package models

import (
	"time"

	"gorm.io/gorm"
)

type Post struct {
	ID      uint   `json:"id"`

	Caption string `json:"caption"`
	UserID  uint   `json:"user_id"`
	User    *User   `gorm:"constraint,OnUpdate:CASCADE,OnDelete:SET NULL;" json:"user,omitempty"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`

	Attachments []Attachment `gorm:"foreignKey:PostID" json:"attachments,omitempty"`
}
