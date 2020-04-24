package authentication

import (
	"math/rand"
	"net/http"
	"rakoon/rakoon-back/db"
	"rakoon/rakoon-back/authentication/types"
	"time"
	"os"

	"encoding/base64"
	"encoding/json"
    "crypto/hmac"
	"crypto/sha256"
	"golang.org/x/crypto/bcrypt"

	"github.com/gin-gonic/gin"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func Connect(c *gin.Context) {
	var jwtToken string
	var connection types.User
	err := c.BindJSON(&connection)

	// Check formatting
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Incorrect input data": err.Error()})
		return
	}

	// Fetch the user
	var user types.User
	errs := db.DB.Where("name = ?", connection.Name).First(&user).GetErrors()

	if (len(errs) != 0) {
		c.JSON(404, gin.H{
			"message": "Incorrect user name or password.",
		})
		return
	}

	check := checkPasswordHash(connection.Password + user.Salt, user.Password)

	if (check == false) {
		c.JSON(404, gin.H{
			"message": "Incorrect user name or password.",
		})
		return
	}

	jwtToken = generateToken(user.Name)

	c.JSON(200, gin.H{
		"token": jwtToken,
	})
	return
}

// Subscribe a new user
func Subscribe(c *gin.Context) {
	var subscription types.User
	err := c.BindJSON(&subscription)

	// Check formatting
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Incorrect input data": err.Error()})
		return
	}

	// Check if the user name is already taken
	var user types.User
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

	// Subscription success
	c.JSON(200, gin.H{
		"message": "subscribe",
	})

}

func generateToken(name string) (string) {
	var token string
	var header *types.JwtHeader
	var payload *types.JwtPayload
	var signature string
	var secret = os.Getenv("SECRET_KEY")
	var alg = os.Getenv("ALG")
	var typ = os.Getenv("TYP")

	header = new(types.JwtHeader)
	header.Alg = alg
	header.Typ = typ
	jsonHeader, _ := json.Marshal(header)
	encHeader := base64.RawURLEncoding.EncodeToString([]byte(string(jsonHeader)))

	payload = new(types.JwtPayload)
	payload.Name = name
	jsonPayload, _ := json.Marshal(payload)
	encPayload := base64.RawURLEncoding.EncodeToString([]byte(string(jsonPayload)))

	toEncrypt := encHeader + "." + encPayload
	
	key := []byte(secret)
    h := hmac.New(sha256.New, key)
    h.Write([]byte(toEncrypt))
    signature = base64.RawURLEncoding.EncodeToString(h.Sum(nil))

    token = encHeader + "." + encPayload + "." + signature

	return token
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
