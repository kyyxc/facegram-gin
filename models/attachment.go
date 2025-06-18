package models

import "gorm.io/gorm"

type Attachment struct {
	gorm.Model
	ID          uint   `gorm:"primaryKey"`
	StoragePath string `json:"storage_path"`
	PostID      uint   `json:"post_id"`
	Post        Post   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
