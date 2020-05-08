package db

import (
	"fmt"
	"os"

	"github.com/jinzhu/gorm"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// DB Single database connection instance
var DB *gorm.DB

// DB2 golbal sqlx  connection instance
var DB2 *sqlx.DB

//InitDb func
func InitDb() {
	var err error
	var err2 error
	var dbHost string = os.Getenv("DB_HOST")
	var dbPort string = os.Getenv("DB_PORT")
	var dbUser string = os.Getenv("DB_USER")
	var dbName string = os.Getenv("DB_NAME")

	// GORM
	var dbConnection string = fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbName)
	DB, err = gorm.Open("postgres", dbConnection)
	if err != nil {
		fmt.Println("Error connecting to the database:", err)
	}

	//SQLX
	var dbConnection2 string = fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbName)
	DB2, err2 = sqlx.Connect("postgres", dbConnection2)
	if err2 != nil {
		fmt.Println("Error connecting to the database:", err)
	}

}
