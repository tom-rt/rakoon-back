package main

import (
	"rakoon/rakoon-back/db"
	"rakoon/rakoon-back/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	// middleware.InitMiddleware(r)
	db.InitDb()
	routes.InitRoutes(r)

	r.Run(":8081")
}
