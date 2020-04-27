package middleware

import (
	"rakoon/rakoon-back/authentication"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

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

	// Check a token is present
	_, checkToken := c.Request.Header["Authorization"]
	if checkToken == false {
		c.JSON(401, gin.H{
			"message": "No token found provided",
		})
		c.Abort()
		return
	}

	// Check the token is formatted correctly
	authorization := c.Request.Header["Authorization"][0]
	token := strings.Split(authorization, "Bearer ")[1]
	splittedToken := strings.Split(token, ".")
	if len(splittedToken) != 3 {
		c.JSON(401, gin.H{
			"message": "Bad token",
		})
		c.Abort()
		return
	}

	// Fetching token data
	header := splittedToken[0]
	payload := splittedToken[1]
	signature := splittedToken[2]

	authenticity, message := authentication.VerifyToken(string(header), string(payload), string(signature))

	if authenticity == false {
		c.JSON(401, gin.H{
			"message": message,
		})
		c.Abort()
		return
	}
}
