package db

import (
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// DB golbal sqlx  connection instance
var DB *sqlx.DB

//InitDb func
func InitDb() {
	var err error
	var dbHost string = os.Getenv("DB_HOST")
	var dbPort string = os.Getenv("DB_PORT")
	var dbUser string = os.Getenv("DB_USER")
	var dbName string = os.Getenv("DB_NAME")

	//SQLX
	var dbConnection2 string = fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbName)
	DB, err = sqlx.Connect("postgres", dbConnection2)
	if err != nil {
		fmt.Println("Error connecting to the database:", err)
		os.Exit(1)
	}

}
