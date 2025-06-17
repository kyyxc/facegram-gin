package database

import (
	"database/sql"
	"fmt"
	"log"
	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func ConnectDB() {
	dsn := "root:@tcp(127.0.0.1:3306)/facegram"
	
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Error opening DB: ", err)
	}
	
	if err = db.Ping(); err != nil {
		log.Fatal("DB connection failed: ", err)
	} 

	fmt.Println("Database connection successfully")
}