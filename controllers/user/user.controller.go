package user

import (
	"fmt"
	"net/http"
	"rakoon/rakoon-back/models"

	"github.com/gin-gonic/gin"
)

func Delete(c *gin.Context) {
	var data models.Username
	err := c.BindJSON(&data)

	// Check input formatting
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Incorrect input data": err.Error()})
		return
	}

	// Delete the user in db
	fmt.Println(data.Name)
	// db.DB.Where("name = ?", data.Name).Delete(&types.User{})

	c.JSON(200, gin.H{
		"message": "User removed",
	})

}

func Get(c *gin.Context) {
	var data models.Username
	err := c.BindJSON(&data)

	// Check input formatting
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Incorrect input data": err.Error()})
		return
	}

	// Check if the user is present in db
	user, err := models.GetUser(data.Name)
	if err != nil {
		c.JSON(404, gin.H{
			"message": "Incorrect user name or password.",
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
