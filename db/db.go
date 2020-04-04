package db

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// DB Single database connecton instance
var DB *gorm.DB

func InitDb() {
	var err error
	DB, err = gorm.Open("postgres", "host=localhost port=5432 user=rakoon dbname=rakoon_user sslmode=disable")
	if err != nil {
		fmt.Println("Error connecting to the database:", err)
	}
}