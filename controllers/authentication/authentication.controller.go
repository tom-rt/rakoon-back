package authentication

import (
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"rakoon/rakoon-back/controllers/utils"
	"rakoon/rakoon-back/models"
	"strings"
	"time"

	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// Connect controller function
func Connect(c *gin.Context) {
	var connection models.User
	err := c.BindJSON(&connection)

	// Check input formatting
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
	check := checkPasswordHash(connection.Password+user.Salt, user.Password)
	if check == false {
		c.JSON(404, gin.H{
			"message": "Incorrect user name or password.",
		})
		return
	}

	// Setting reauth to false
	models.SetReauthByName(user.Name, false)

	// Generate and return a token
	jwtToken := generateToken(user.Name)
	c.JSON(200, gin.H{
		"token": jwtToken,
	})
	return
}

// LogOut controller function
func LogOut(c *gin.Context) {
	var logout models.UserID
	err := c.BindJSON(&logout)

	// Check input formatting
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Incorrect input data": err.Error()})
		return
	}

	// Check if the user exists
	if !userIDExists(logout.ID) {
		c.JSON(409, gin.H{
			"message": "User does not exist.",
		})
		return
	}

	// Setting reauth var to true to force the user to reconnect
	models.SetReauthByID(logout.ID, true)
	c.JSON(http.StatusOK, gin.H{
		"message": "User logged out.",
	})
}

// Subscribe a new user
func Subscribe(c *gin.Context) {
	var subscription models.User
	err := c.BindJSON(&subscription)

	// Check formatting
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Incorrect input data": err.Error()})
		return
	}

	// Check if the user name is already taken
	if userNameExists(subscription.Name) {
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

	// Create the user in db
	subscription.Password = hash
	subscription.Salt = salt
	subscription.Reauth = false
	subscription.LastLogin = time.Now()
	models.CreateUser(subscription)

	// Generate connection token
	token := generateToken(subscription.Name)

	// Subscription success
	c.JSON(200, gin.H{
		"token": token,
	})

}

// RefreshToken controller function
func RefreshToken(c *gin.Context) {
	authorization := c.Request.Header["Authorization"][0]
	token := strings.Split(authorization, "Bearer ")[1]
	splittedToken := strings.Split(token, ".")
	if len(splittedToken) != 3 {
		c.JSON(401, gin.H{
			"message": "Bad token",
		})
		return
	}

	// Fetching token data
	encHeader := splittedToken[0]
	encPayload := splittedToken[1]
	signature := splittedToken[2]

	// Decode payload
	decPayloadByte, err := base64.RawURLEncoding.DecodeString(encPayload)
	decPayload := string(decPayloadByte)
	payload := new(models.JwtPayload)
	err = json.Unmarshal([]byte(decPayload), payload)
	if err != nil {
		c.JSON(401, gin.H{
			"message": "Bad token",
		})
		return
	}

	// Check signature
	encSignature := GenerateSignature(encHeader, encPayload)
	if encSignature != signature {
		c.JSON(401, gin.H{
			"message": "Bad signature",
		})
		return
	}

	// Check expiration duration
	duration := utils.NowAsUnixMilli() - payload.Iat
	if duration > utils.HoursToMilliseconds(24) {
		models.SetReauthByName(payload.Name, true)
		c.JSON(401, gin.H{
			"message": "Token expired more than a week ago, please reconnect.",
		})
		return
	}

	// Check if the user has to re authenticate
	var reauth bool
	reauth, err = GetReauth(payload.Name)
	if reauth {
		c.JSON(401, gin.H{
			"message": "Please reconnect.",
		})
		return
	} else if err != nil {
		c.JSON(404, gin.H{
			"message": "User does not exist.",
		})
		return
	}

	newToken := generateToken(payload.Name)
	c.JSON(200, gin.H{
		"message": "Token refreshed.",
		"token":   newToken,
	})
	return
}

func generateToken(name string) string {
	var header *models.JwtHeader
	var payload *models.JwtPayload
	const alg = "HS256"
	const typ = "JWT"

	// Building and encrypting header
	header = new(models.JwtHeader)
	header.Alg = alg
	header.Typ = typ
	jsonHeader, _ := json.Marshal(header)
	encHeader := base64.RawURLEncoding.EncodeToString([]byte(string(jsonHeader)))

	// Building and encrypting payload
	payload = new(models.JwtPayload)
	payload.Name = name
	now := utils.NowAsUnixMilli()
	payload.Iat = now
	payload.Exp = now + utils.MinutesToMilliseconds(15)
	jsonPayload, _ := json.Marshal(payload)
	encPayload := base64.RawURLEncoding.EncodeToString([]byte(string(jsonPayload)))

	// Building signature and token
	signature := GenerateSignature(encHeader, encPayload)
	token := encHeader + "." + encPayload + "." + signature

	return token
}

// GenerateSignature controller function
func GenerateSignature(encHeader string, encPayload string) string {
	var secret = os.Getenv("SECRET_KEY")
	key := []byte(secret)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(encHeader + "." + encPayload))
	signature := base64.RawURLEncoding.EncodeToString(h.Sum(nil))
	return signature
}

// VerifyToken controller: This function checks if the user has to reconnect and if the token is valid. It is only used in the middleware
func VerifyToken(encHeader string, encPayload string, encSignature string) (bool, string) {
	// Decode payload
	decPayloadByte, err := base64.RawURLEncoding.DecodeString(encPayload)
	decPayload := string(decPayloadByte)
	payload := new(models.JwtPayload)
	err = json.Unmarshal([]byte(decPayload), payload)
	if err != nil {
		return false, "Bad token"
	}

	// Check if the user has to reconnect
	var reauth bool
	reauth, err = GetReauth(payload.Name)
	if reauth {
		return false, "Please reconnect"
	} else if err != nil {
		return false, "User does not exist."
	}

	checkSignature := GenerateSignature(encHeader, encPayload)
	if encSignature != checkSignature {
		return false, "Bad signature"
	}

	// Check token validity date
	now := utils.NowAsUnixMilli()
	if now >= payload.Exp {
		return false, "Token has expired"
	}

	return true, "Token valid"
}

func userNameExists(name string) bool {
	_, err := models.GetUserByName(name)
	fmt.Println(err)
	if err != nil {
		return false
	}
	return true
}

func userIDExists(ID string) bool {
	_, err := models.GetUserByID(ID)
	fmt.Println(err)
	if err != nil {
		return false
	}
	return true
}

// GetReauth function: fetching in db a user's reauth value
func GetReauth(username string) (bool, error) {
	user, err := models.GetUserByName(username)
	return user.Reauth, err
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// Comparing a salted password and a hash
func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// Generate a random string to use as a password salt
func generateSalt(saltLength int) string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	salt := make([]byte, saltLength)
	for i := range salt {
		salt[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(salt)
}
