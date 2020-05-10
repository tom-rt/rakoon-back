package user

import (
	"rakoon/rakoon-back/models"

	"github.com/gin-gonic/gin"
)

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
}
