package middleware

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// InitMiddleware function that inits all the middleware
func InitMiddleware(r *gin.Engine) {
	r.Use(CorsHandling)
	r.Use(JwtHandling)
}

//CorsHandling middleware
func CorsHandling(c *gin.Context) {
	cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"POST", "OPTIONS", "GET", "PUT", "DELETE"},
		AllowHeaders:     []string{"Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization", "accept", "origin", "Cache-Control", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	})
}

//JwtHandling middleware
func JwtHandling(c *gin.Context) {
}
