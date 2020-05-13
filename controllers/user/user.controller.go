package user

import (
	"rakoon/rakoon-back/controllers/authentication"
	"rakoon/rakoon-back/models"
	"time"

	"github.com/gin-gonic/gin"
)

// Create a new user
func Create(c *gin.Context) {
	var subscription models.User
	err := c.BindJSON(&subscription)

	// Check formatting
	if err != nil {
		c.JSON(400, gin.H{"Incorrect input data": err.Error()})
		return
	}

	// Check if the user name is already taken
	if authentication.UserNameExists(subscription.Name) {
		c.JSON(409, gin.H{
			"message": "Conflict: username already taken.",
		})
		return
	}

	// Salt password
	salt := authentication.GenerateSalt(10)
	saltedPassword := subscription.Password + salt

	// Generate hash
	hash, _ := authentication.HashPassword(saltedPassword)

	// Create the user in db
	subscription.Password = hash
	subscription.Salt = salt
	subscription.Reauth = false
	subscription.LastLogin = time.Now()
	models.CreateUser(subscription)

	// Generate connection token
	token := authentication.GenerateToken(subscription.Name)

	// Subscription success
	c.JSON(200, gin.H{
		"token": token,
	})

}

// Get user controller function
func Get(c *gin.Context) {
	var userID = c.Param("id")

	user, err := models.GetUserPublicByID(userID)
	if err != nil {
		c.JSON(404, gin.H{
			"message": "Bad user id.",
		})
		return
	}

	c.JSON(200, gin.H{
		"data": user,
	})
	return
}

func Update(c *gin.Context) {
	var update models.UserUpdate
	var err = c.BindJSON(&update)

	if err != nil {
		c.JSON(400, gin.H{"Incorrect input data": err.Error()})
		return
	}

	if !authentication.UserIDExists(update.ID) {
		c.JSON(404, gin.H{
			"message": "Bad user Id.",
		})
		return
	}

	models.UpdateUser(update)

	newToken := authentication.GenerateToken(update.Name)
	c.JSON(200, gin.H{
		"message": "User updated",
		"token":   newToken,
	})
	return

}

// Delete user controller function
func Delete(c *gin.Context) {
	var userID = c.Param("id")

	_, err := models.GetUserByID(userID)
	if err != nil {
		c.JSON(404, gin.H{
			"message": "User does not exist",
		})
		return
	}

	models.DeleteUser(userID)

	c.JSON(200, gin.H{
		"message": "User removed",
	})

}
