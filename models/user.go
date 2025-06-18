package models

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username  string `gorm:"unique" json:"username"`
	FullName  string `json:"full_name"`
	Password  string `json:"password"`
	Bio       string `json:"bio"`
	IsPrivate bool   `json:"is_private"`
	Posts      []Post `gorm:"foreignKey:UserID"`
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
