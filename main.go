package main

import (
	"fmt"
	"os"

	"rakoon/rakoon-back/db"
	"rakoon/rakoon-back/routes"

	"github.com/tom-rt/goberge"
)

func main() {
	// For security reasons, I check if the secret key is defined, if not I quit the program.
	var secret = os.Getenv("SECRET_KEY")
	if len(secret) <= 0 {
		fmt.Println("ERROR: secret key is not defined.")
		os.Exit(1)
	}

	// Test using an external lib
	goberge.Goberge()

	db.InitDB()
	r := routes.SetupRouter()
	r.Run(":8081")
}
