package authentication

import (
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"rakoon/rakoon-back/controllers/utils"
	"rakoon/rakoon-back/models"
	"strings"

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

	// Setting reauth to false, updated last login field
	models.RefreshUserConnection(user.Name, false)

	// Generate and return a token
	input := new(models.JwtInput)
	input.Name = user.Name
	input.IsAdmin = nil
	jwtToken := GenerateToken(input)
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
	if !UserIDExists(logout.ID) {
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

	input := new(models.JwtInput)
	input.Name = payload.Name
	input.IsAdmin = nil
	newToken := GenerateToken(input)
	c.JSON(200, gin.H{
		"message": "Token refreshed.",
		"token":   newToken,
	})
	return
}

// GenerateToken function
func GenerateToken(input *models.JwtInput) string {
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
	payload.Name = input.Name
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

// UserNameExists function
func UserNameExists(name string) bool {
	_, err := models.GetUserByName(name)
	if err != nil {
		return false
	}
	return true
}

func UserIDExists(ID string) bool {
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

// HashPassword function
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// Comparing a salted password and a hash
func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// GenerateSalt a random string to use as a password salt
func GenerateSalt(saltLength int) string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	salt := make([]byte, saltLength)
	for i := range salt {
		salt[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(salt)
}
