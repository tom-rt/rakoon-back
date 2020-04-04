package authentication

import (
	"fmt"
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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//check if the user exists
	var user models.User
	db.DB.First(&user)
	fmt.Println(user)

	// Salt password
	salt := generateSalt(10)
	saltedPassword := subscription.Password + salt

	// Generate hash
	hash, _ := hashPassword(saltedPassword)

	// Create the user
	subscription.Password = hash
	subscription.Salt = salt
	subscription.LastLogin = time.Now()

	ret := db.DB.NewRecord(subscription)
	fmt.Println(ret)
	var reg = db.DB.Create(&subscription)
	fmt.Printf("%+v\n", reg.Error)
	ret = db.DB.NewRecord(subscription)
	fmt.Println(ret)
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
