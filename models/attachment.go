package models

import (
	"time"

	"gorm.io/gorm"
)

type Attachment struct {
	ID uint `json:"id"`

	StoragePath string `json:"storage_path"`
	PostID      uint   `json:"post_id"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`

	Post *Post `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"post,omitempty"`
}

func (Attachment) TableName() string {
	return "post_attachments"
}
