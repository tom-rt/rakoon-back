package db

import (
	"fmt"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// DB Single database connection instance
var DB *gorm.DB

func InitDb() {
	var err error
	var dbHost string = os.Getenv("DB_HOST")
	var dbPort string = os.Getenv("DB_PORT")
	var dbUser string = os.Getenv("DB_USER")
	var dbName string = os.Getenv("DB_NAME")
	
	var dbConnection string = fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbName)
	DB, err = gorm.Open("postgres", dbConnection)
	if err != nil {
		fmt.Println("Error connecting to the database:", err)
	}
}