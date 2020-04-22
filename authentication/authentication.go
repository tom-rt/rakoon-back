package authentication

import (
	"math/rand"
	"net/http"
	"rakoon/user-service/db"
	"rakoon/user-service/models"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/gin-gonic/gin"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func login(c *gin.Context) {

}

// Subscribe function
func Subscribe(c *gin.Context) {

	var subscription models.User
	err := c.BindJSON(&subscription)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Bind json error": err.Error()})
		return
	}

	//check if the user exists
	var user models.User
	errors := db.DB.Where("name = ?", subscription.Name).First(&user).GetErrors()
	if (len(errors) == 0) {
		c.JSON(409, gin.H{
			"message": "Conflict: username already taken.",
		})
		return
	}

	// Salt password
	salt := generateSalt(10)
	saltedPassword := subscription.Password + salt

	// Generate hash
	hash, _ := hashPassword(saltedPassword)

	// Create the user
	subscription.Password = hash
	subscription.Salt = salt
	subscription.LastLogin = time.Now()

	db.DB.NewRecord(subscription)
	db.DB.Create(&subscription)
	// ret = db.DB.NewRecord(subscription)
	c.JSON(200, gin.H{
		"message": "subscribe",
	})

}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func generateSalt(saltLength int) string {
	salt := make([]byte, saltLength)
	for i := range salt {
		salt[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(salt)
}
