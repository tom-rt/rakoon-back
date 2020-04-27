package authentication

import (
	"math/rand"
	"net/http"
	"os"
	"rakoon/rakoon-back/authentication/types"
	"rakoon/rakoon-back/db"
	"time"

	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"

	"golang.org/x/crypto/bcrypt"

	"github.com/gin-gonic/gin"
)

func Connect(c *gin.Context) {
	var jwtToken string
	var connection types.User
	err := c.BindJSON(&connection)

	// Check input formatting
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Incorrect input data": err.Error()})
		return
	}

	// Fetch the user
	var user types.User
	errs := db.DB.Where("name = ?", connection.Name).First(&user).GetErrors()

	if len(errs) != 0 {
		c.JSON(404, gin.H{
			"message": "Incorrect user name or password.",
		})
		return
	}

	check := checkPasswordHash(connection.Password+user.Salt, user.Password)

	if check == false {
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
	if len(errors) == 0 {
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

func generateToken(name string) string {
	var token string
	var signature string
	var header *types.JwtHeader
	var payload *types.JwtPayload
	var alg = "HS256"
	var typ = "JWT"
	var now int64

	header = new(types.JwtHeader)
	header.Alg = alg
	header.Typ = typ
	jsonHeader, _ := json.Marshal(header)
	encHeader := base64.RawURLEncoding.EncodeToString([]byte(string(jsonHeader)))

	payload = new(types.JwtPayload)
	payload.Name = name
	now = NowAsUnixMilli()
	payload.Iat = now
	payload.Exp = now + 60000
	jsonPayload, _ := json.Marshal(payload)
	encPayload := base64.RawURLEncoding.EncodeToString([]byte(string(jsonPayload)))

	signature = GenerateSignature(encHeader, encPayload)

	token = encHeader + "." + encPayload + "." + signature

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

func VerifyToken(encHeader string, encPayload string, encSignature string) (bool, string) {

	// Decode payload
	decPayloadByte, err := base64.RawURLEncoding.DecodeString(encPayload)
	decPayload := string(decPayloadByte)
	payload := new(types.JwtPayload)
	err = json.Unmarshal([]byte(decPayload), payload)

	if err != nil {
		return false, "Bad token"
	}

	checkSignature := GenerateSignature(encHeader, encPayload)

	if encSignature != checkSignature {
		return false, "Bad signature"
	}

	// Check token validity date
	now := NowAsUnixMilli()

	if now >= payload.Exp {
		return false, "Token has expired"
	}

	return true, "Token valid"
}

func NowAsUnixMilli() int64 {
	return time.Now().UnixNano() / 1e6
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
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	salt := make([]byte, saltLength)
	for i := range salt {
		salt[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(salt)
}
