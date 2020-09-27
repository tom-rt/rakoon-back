package middleware

import (
	"rakoon/rakoon-back/handlers/authentication"
	"strings"

	"github.com/gin-gonic/gin"
)

//AdminJwtHandling middleware, checks if the token is well formatted and has expired
func AdminJwtHandling(c *gin.Context) {
	var token string

	// Check a token is present
	_, checkToken := c.Request.Header["Authorization"]
	if checkToken == false {
		c.JSON(403, gin.H{
			"message": "No token provided",
		})
		c.Abort()
		return
	}

	// Check the token is formatted correctly
	authorization := c.Request.Header["Authorization"][0]
	bearer := strings.Split(authorization, "Bearer ")
	if len(bearer) != 2 {
		c.JSON(403, gin.H{
			"message": "Bad token",
		})
		c.Abort()
		return
	}
	token = bearer[1]

	splittedToken := strings.Split(token, ".")
	if len(splittedToken) != 3 {
		c.JSON(403, gin.H{
			"message": "Bad token",
		})
		c.Abort()
		return
	}

	// Fetching token data
	header := splittedToken[0]
	payload := splittedToken[1]
	signature := splittedToken[2]

	// Check token validity
	validity, message, status, id, isAdmin := authentication.VerifyToken(string(header), string(payload), string(signature))
	if validity == false {
		c.JSON(status, gin.H{
			"message": message,
		})
		c.Abort()
		return
	}

	if !isAdmin {
		c.JSON(status, gin.H{
			"message": "Admin privileges are required to access this endpoint.",
		})
		c.Abort()
		return
	}

	c.Set("id", id)
	c.Next()
}

//JwtHandling middleware, checks if the token is well formatted and has expired
func JwtHandling(c *gin.Context) {
	var token string

	// Check a token is present
	_, checkToken := c.Request.Header["Authorization"]
	if checkToken == false {
		c.JSON(403, gin.H{
			"message": "No token provided",
		})
		c.Abort()
		return
	}

	// Check the token is formatted correctly
	authorization := c.Request.Header["Authorization"][0]
	bearer := strings.Split(authorization, "Bearer ")
	if len(bearer) != 2 {
		c.JSON(403, gin.H{
			"message": "Bad token",
		})
		c.Abort()
		return
	}
	token = bearer[1]

	splittedToken := strings.Split(token, ".")
	if len(splittedToken) != 3 {
		c.JSON(403, gin.H{
			"message": "Bad token",
		})
		c.Abort()
		return
	}

	// Fetching token data
	header := splittedToken[0]
	payload := splittedToken[1]
	signature := splittedToken[2]

	// Check token validity
	validity, message, status, id, _ := authentication.VerifyToken(string(header), string(payload), string(signature))
	if validity == false {
		c.JSON(status, gin.H{
			"message": message,
		})
		c.Abort()
		return
	}

	c.Set("id", id)
	c.Next()
}
