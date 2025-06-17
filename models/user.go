package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string `gorm:"unique" json:"username"`
	FullName string `json:"full_name"`
	Password string `json:"password"`
	Bio string `json:"bio"`
	IsPrivate bool `json:"is_private"`
}