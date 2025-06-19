package database

import (
	"facegram/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

var DB *gorm.DB

func ConnectDB() {
	dsn := "root:@tcp(127.0.0.1:3306)/facegram?charset=utf8mb4&parseTime=True&loc=Local"
	var err error

	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect database", err)
	}

	err = DB.AutoMigrate(&models.User{}, &models.Post{})
	if err != nil {
		log.Fatal("Migration failed", err)
	}
}
