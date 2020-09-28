package authentication

import (
	"math/rand"
	"os"
	"rakoon/rakoon-back/models"
	"strconv"
	"strings"
	"time"

	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// RefreshToken controller function
func RefreshToken(c *gin.Context) {
	// Check a token is present
	_, checkToken := c.Request.Header["Authorization"]
	if checkToken == false {
		c.JSON(403, gin.H{
			"message": "No token provided",
		})
		c.Abort()
		return
	}

	// Check if the token is formatted properly
	authorization := c.Request.Header["Authorization"][0]
	bearer := strings.Split(authorization, "Bearer ")
	if len(bearer) != 2 {
		c.JSON(403, gin.H{
			"message": "Bad token",
		})
		return
	}
	token := bearer[1]
	splittedToken := strings.Split(token, ".")
	if len(splittedToken) != 3 {
		c.JSON(403, gin.H{
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
		c.JSON(403, gin.H{
			"message": "Bad token",
		})
		return
	}

	// Check signature
	encSignature := GenerateSignature(encHeader, encPayload)
	if encSignature != signature {
		c.JSON(403, gin.H{
			"message": "Bad signature",
		})
		return
	}

	// Check expiration duration
	duration := nowAsUnixMilli() - payload.Iat
	var refreshLimit int
	var envRefreshLimit string = os.Getenv("TOKEN_LIMIT_HOURS")

	if envRefreshLimit != "" {
		refreshLimit, _ = strconv.Atoi(envRefreshLimit)
	} else {
		refreshLimit = 24
	}

	// if duration > hoursToMilliseconds(refreshLimit) {
	if duration > minutesToMilliseconds(refreshLimit) {
		models.SetReauth(payload.ID, true)
		c.JSON(401, gin.H{
			"message": "Token has expired and cannot be refreshed, please reconnect",
		})
		return
	}

	// Check if the user has to re authenticate
	var reauth bool
	reauth, err = GetReauth(payload.ID)
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

	newToken := GenerateToken(payload.ID, payload.IsAdmin)
	c.JSON(200, gin.H{
		"userId":  payload.ID,
		"token":   newToken,
		"isAdmin": payload.IsAdmin,
	})
	return
}

// GenerateToken function
func GenerateToken(id int, isAdmin bool) string {
	var header *models.JwtHeader
	var payload *models.JwtPayload
	const alg = "HS256"
	const typ = "JWT"
	var validityLimit int
	var envValidityLimit string = os.Getenv("TOKEN_VALIDITY_MINUTES")

	if envValidityLimit != "" {
		validityLimit, _ = strconv.Atoi(envValidityLimit)
	} else {
		validityLimit = 15
	}

	// Building and encrypting header
	header = new(models.JwtHeader)
	header.Alg = alg
	header.Typ = typ
	// Error return is ignored here as it cant fail.
	jsonHeader, _ := json.Marshal(header)
	encHeader := base64.RawURLEncoding.EncodeToString([]byte(string(jsonHeader)))

	// Building and encrypting payload
	payload = new(models.JwtPayload)
	payload.ID = id
	payload.IsAdmin = isAdmin
	now := nowAsUnixMilli()
	payload.Iat = now
	payload.Exp = now + minutesToMilliseconds(validityLimit)
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
func VerifyToken(encHeader string, encPayload string, encSignature string) (isValid bool, message string, status int, id int, isAdmin bool) {
	// Decode payload
	decPayloadByte, err := base64.RawURLEncoding.DecodeString(encPayload)
	decPayload := string(decPayloadByte)
	payload := new(models.JwtPayload)
	err = json.Unmarshal([]byte(decPayload), payload)
	if err != nil {
		return false, "Bad token", 403, -1, false
	}

	// Check if the user has to reconnect
	var reauth bool
	reauth, err = GetReauth(payload.ID)
	if reauth {
		return false, "Please reconnect", 401, -1, false
	} else if err != nil {
		return false, "User id in token payload does not exist.", 403, -1, false
	}

	checkSignature := GenerateSignature(encHeader, encPayload)
	if encSignature != checkSignature {
		return false, "Bad signature", 403, -1, false
	}

	// Check token validity date
	now := nowAsUnixMilli()
	if now >= payload.Exp {
		return false, "Token expired.", 401, -1, false
	}

	return true, "Token valid", 200, payload.ID, payload.IsAdmin
}

// UserNameExists function
func UserNameExists(name string) bool {
	_, err := models.GetUserByName(name)
	if err != nil {
		return false
	}
	return true
}

func isAdmin(ID int) bool {
	ret, err := models.IsAdmin(ID)
	if err != nil {
		return false
	}
	return ret
}

func UserIDExists(ID int) bool {
	_, err := models.GetUserByID(ID)
	if err != nil {
		return false
	}
	return true
}

// GetReauth function: fetching in db a user's reauth value
func GetReauth(ID int) (bool, error) {
	user, err := models.GetUserByID(ID)
	return user.Reauth, err
}

// HashPassword function
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// Comparing a salted password and a hash
func CheckPasswordHash(password, hash string) bool {
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

func minutesToMilliseconds(min int) int {
	return min * 60000
}

func hoursToMilliseconds(hours int) int {
	return hours * 3600000
}

func nowAsUnixMilli() int {
	return int(time.Now().UnixNano() / 1e6)
}
