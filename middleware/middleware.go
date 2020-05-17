package middleware

import (
	"rakoon/rakoon-back/controllers/authentication"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

//CorsHandling middleware, allows cross origins
func CorsHandling(c *gin.Context) {

	cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"POST", "OPTIONS", "GET", "PUT", "DELETE"},
		AllowHeaders:     []string{"Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization", "accept", "origin", "Cache-Control", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	})
}

//JwtHandling middleware, checks if the token is well formatted and has expired
func JwtHandling(c *gin.Context) {
	var token string

	// Check a token is present
	_, checkToken := c.Request.Header["Authorization"]
	if checkToken == false {
		c.JSON(401, gin.H{
			"message": "No token provided",
		})
		c.Abort()
		return
	}

	// Check the token is formatted correctly
	authorization := c.Request.Header["Authorization"][0]
	bearer := strings.Split(authorization, "Bearer ")

	if len(bearer) != 2 {
		c.JSON(401, gin.H{
			"message": "Bad token",
		})
		c.Abort()
		return
	} else {
		token = bearer[1]
	}

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

	// Check token authenticity
	authenticity, message := authentication.VerifyToken(string(header), string(payload), string(signature))
	if authenticity == false {
		c.JSON(401, gin.H{
			"message": message,
		})
		c.Abort()
		return
	}

	return
}
