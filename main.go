package main

import (
	"fmt"
	"os"

	"rakoon/rakoon-back/db"
	"rakoon/rakoon-back/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	// For security reasons, I check if the secret key is defined, if not I quit the program.
	var secret = os.Getenv("SECRET_KEY")
	if len(secret) <= 0 {
		fmt.Println("ERROR: secret key is not defined.")
		os.Exit(1)
	}

	r := gin.Default()
	db.InitDb()
	routes.InitRoutes(r)
	r.Run(":8081")
}
