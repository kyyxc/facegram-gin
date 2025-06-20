package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID uint `json:"id"`

	Username  string `gorm:"unique" json:"username"`
	FullName  string `json:"full_name"`
	Password  string `json:"-"`
	Bio       string `json:"bio"`
	IsPrivate bool   `json:"is_private"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`

	Posts      []Post   `gorm:"foreignKey:UserID" json:"posts,omitempty"`
	Followings []Follow `gorm:"foreignKey:FollowerID" json:"followings,omitempty"`
	Followers  []Follow `gorm:"foreignKey:FollowingID" json:"followers,omitempty"`
}

func (u *User) HashPassword() error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashed)
	return nil
}

func (u *User) CheckPassword(input string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(input)) == nil
}
