package main

import (
	"rakoon/rakoon-back/db"
	"rakoon/rakoon-back/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	db.InitDb()
	routes.InitRoutes(r)

	r.Run(":8081")
}
