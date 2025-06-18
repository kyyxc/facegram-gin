package models

import "gorm.io/gorm"

type Post struct {
	gorm.Model
	Caption    string       `json:"caption"`
	UserID     uint         `json:"user_id"`
	User       User         `gorm:"constraint,OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Attachmenta []Attachment `gorm:"foreignKey:PostID"`
}
