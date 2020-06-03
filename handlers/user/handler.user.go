package user

import (
	"fmt"
	"net/http"
	"rakoon/rakoon-back/handlers/authentication"
	"rakoon/rakoon-back/models"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// Create a new user
func Create(c *gin.Context) {
	var subscription models.User
	err := c.BindJSON(&subscription)

	// Check formatting
	if err != nil {
		c.JSON(400, gin.H{"Incorrect ay input data": err.Error()})
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
	id := models.CreateUser(subscription)

	// Generate connection token
	token := authentication.GenerateToken(id)

	// Subscription success
	c.JSON(201, gin.H{
		"token": token,
	})

	return
}

// Get user controller function
func Get(c *gin.Context) {
	var ID = c.Param("id")
	var tokenID = fmt.Sprintf("%v", c.MustGet("id"))

	if !matchIDs(c, ID, tokenID) {
		return
	}

	// Already checked in jwt verif ?
	user, _ := models.GetUserPublic(ID)

	c.JSON(200, gin.H{
		"data": user,
	})

	return
}

// Update user fields
func Update(c *gin.Context) {
	var update models.UserUpdate
	var err = c.BindJSON(&update)
	var tokenID = fmt.Sprintf("%v", c.MustGet("id"))

	if err != nil {
		c.JSON(400, gin.H{"Incorrect input data": err.Error()})
		return
	}
	update.ID = c.Param("id")

	if !matchIDs(c, update.ID, tokenID) {
		return
	}

	models.UpdateUser(update)

	c.JSON(200, gin.H{
		"message": "User updated",
	})

	return
}

// UpdatePassword function: update user's password
func UpdatePassword(c *gin.Context) {
	var user models.UserPassword
	var err = c.BindJSON(&user)

	var tokenID = fmt.Sprintf("%v", c.MustGet("id"))

	if err != nil {
		c.JSON(400, gin.H{"Incorrect input data": err.Error()})
		return
	}
	user.ID = c.Param("id")

	if !matchIDs(c, user.ID, tokenID) {
		return
	}

	// Salt password
	salt := authentication.GenerateSalt(10)
	saltedPassword := user.Password + salt

	// Generate hash
	hash, _ := authentication.HashPassword(saltedPassword)
	user.Password = hash

	models.UpdateUserPassword(user, salt)

	c.JSON(200, gin.H{
		"message": "Password updated",
	})

	return
}

// Archive a user (soft delete)
func Archive(c *gin.Context) {
	var ID = c.Param("id")
	var tokenID = fmt.Sprintf("%v", c.MustGet("id"))

	if !matchIDs(c, ID, tokenID) {
		return
	}

	models.ArchiveUser(ID)

	c.JSON(200, gin.H{
		"message": "User archived",
	})

	return
}

// Delete user controller function
func Delete(c *gin.Context) {
	var ID = c.Param("id")
	var tokenID = fmt.Sprintf("%v", c.MustGet("id"))

	if !matchIDs(c, ID, tokenID) {
		return
	}

	models.DeleteUser(ID)

	c.JSON(200, gin.H{
		"message": "User removed",
	})
	return
}

// Connect controller function
func Connect(c *gin.Context) {
	var connection models.User
	err := c.BindJSON(&connection)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Incorrect input data": err.Error()})
		return
	}

	// Fetch the user in db
	var user models.User
	user, err = models.GetUserByName(connection.Name)
	if err != nil {
		c.JSON(404, gin.H{
			"message": "Incorrect user name or password.",
		})
		return
	}

	// Check if the provided password is good
	check := authentication.CheckPasswordHash(connection.Password+user.Salt, user.Password)
	if check == false {
		c.JSON(404, gin.H{
			"message": "Incorrect user name or password.",
		})
		return
	}

	// Setting reauth to false, update last login field
	models.RefreshUserConnection(user.Name, false)

	// Generate and return a token
	jwtToken := authentication.GenerateToken(user.ID)
	c.JSON(200, gin.H{
		"token": jwtToken,
	})
	return
}

// LogOut controller function
func LogOut(c *gin.Context) {
	var paramID = c.Param("id")
	var tokenID = fmt.Sprintf("%v", c.MustGet("id"))

	if !matchIDs(c, paramID, tokenID) {
		return
	}

	// Setting reauth var to true to force the user to reconnect
	ID, _ := strconv.Atoi(paramID)
	models.SetReauth(ID, true)
	c.JSON(http.StatusOK, gin.H{
		"message": "User logged out.",
	})
}

// This function checks if the id present in the token (retrieved by the middleware) matches with the id in the route parameters, or in the route body.
func matchIDs(c *gin.Context, ID string, tokenID string) bool {
	_, err := strconv.Atoi(ID)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "Id not valid",
		})
		return false
	}

	if tokenID != ID {
		c.JSON(403, gin.H{
			"message": "Forbidden.",
		})
		return false
	}
	return true
}
