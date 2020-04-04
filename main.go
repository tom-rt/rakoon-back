package main

import (
	"rakoon/user-service/db"
	"rakoon/user-service/middleware"
	"rakoon/user-service/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	middleware.InitMiddleware(r)
	db.InitDb()
	routes.InitRoutes(r)

	r.Run(":8082")
}
