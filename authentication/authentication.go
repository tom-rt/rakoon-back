package authentication

import (
	"math/rand"
	"net/http"
	"os"
	"rakoon/rakoon-back/authentication/types"
	"rakoon/rakoon-back/db"
	"rakoon/rakoon-back/utils"
	"strings"
	"time"

	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func Connect(c *gin.Context) {
	var connection types.User
	err := c.BindJSON(&connection)

	// Check input formatting
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Incorrect input data": err.Error()})
		return
	}

	// Fetch the user in db
	var user types.User
	errs := db.DB.Where("name = ?", connection.Name).First(&user).GetErrors()
	if len(errs) != 0 {
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
	setReauth(user.Name, false)

	// Generate and return a token
	jwtToken := generateToken(user.Name)
	c.JSON(200, gin.H{
		"token": jwtToken,
	})
	return
}

func LogOut(c *gin.Context) {
	var logout types.Logout
	err := c.BindJSON(&logout)

	// Check input formatting
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Incorrect input data": err.Error()})
		return
	}

	// Check if the user exists
	if !userExists(logout.Name) {
		c.JSON(409, gin.H{
			"message": "User does not exist.",
		})
		return
	}

	// Setting reauth var to true to force the user to reconnect
	setReauth(logout.Name, true)
	c.JSON(http.StatusOK, gin.H{
		"message": "User logged out.",
	})
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
	if userExists(subscription.Name) {
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
	db.DB.NewRecord(subscription)
	db.DB.Create(&subscription)

	// Generate connection token
	token := generateToken(subscription.Name)

	// Subscription success
	c.JSON(200, gin.H{
		"token": token,
	})

}

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
	payload := new(types.JwtPayload)
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
		setReauth(payload.Name, true)
		c.JSON(401, gin.H{
			"message": "Token expired more than a week ago, please reconnect.",
		})
		return
	}

	// Check if the user has to re authenticate
	if GetReauth(payload.Name) {
		c.JSON(401, gin.H{
			"message": "Please reconnect.",
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
	var header *types.JwtHeader
	var payload *types.JwtPayload
	const alg = "HS256"
	const typ = "JWT"

	// Building and encrypting header
	header = new(types.JwtHeader)
	header.Alg = alg
	header.Typ = typ
	jsonHeader, _ := json.Marshal(header)
	encHeader := base64.RawURLEncoding.EncodeToString([]byte(string(jsonHeader)))

	// Building and encrypting payload
	payload = new(types.JwtPayload)
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

func GenerateSignature(encHeader string, encPayload string) string {
	var secret = os.Getenv("SECRET_KEY")
	key := []byte(secret)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(encHeader + "." + encPayload))
	signature := base64.RawURLEncoding.EncodeToString(h.Sum(nil))
	return signature
}

// This function checks if the user has to reconnect and if the token is valid. It is only used in the middleware
func VerifyToken(encHeader string, encPayload string, encSignature string) (bool, string) {
	// Decode payload
	decPayloadByte, err := base64.RawURLEncoding.DecodeString(encPayload)
	decPayload := string(decPayloadByte)
	payload := new(types.JwtPayload)
	err = json.Unmarshal([]byte(decPayload), payload)
	if err != nil {
		return false, "Bad token"
	}

	// Check if the user has to reconnect
	if GetReauth(payload.Name) == true {
		return false, "Please reconnect."
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

// Checking in DB if a given username exists
func userExists(name string) bool {
	var user types.User
	errors := db.DB.Where("name = ?", name).First(&user).GetErrors()
	if len(errors) != 0 {
		return false
	}
	return true
}

// Setting in db a user's reauth value
func setReauth(username string, value bool) {
	var user types.User
	user.Name = username
	db.DB.Model(&user).Where("name = ?", username).Update("reauth", value)
}

// Fetching in db a user's reauth value
func GetReauth(username string) bool {
	var user types.User
	db.DB.Where("name = ?", username).First(&user)
	return user.Reauth
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
