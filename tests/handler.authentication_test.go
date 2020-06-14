package test

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"strings"
	"time"

	"rakoon/rakoon-back/db"
	"rakoon/rakoon-back/models"
	"rakoon/rakoon-back/routes"
	"rakoon/rakoon-back/tests/utils"
	"testing"

	"github.com/gin-gonic/gin"
	"gopkg.in/go-playground/assert.v1"
)

//UserConnect object
type Message struct {
	Message string `json:"message" binding:"required"`
}

// Asserts a connected user can refresh his token.
func TestRefreshToken(t *testing.T) {
	db.InitDB()
	var router *gin.Engine = routes.SetupRouter()

	var user models.UserCreate = utils.CreateUser("Tom", "qwerty1234", t, router)

	var url string = "/v1/refresh/token"
	var bearer = "Bearer " + user.Token

	scndRec := httptest.NewRecorder()
	scndReq, _ := http.NewRequest("POST", url, nil)
	scndReq.Header.Add("Content-Type", "application/json")
	scndReq.Header.Add("Authorization", bearer)

	router.ServeHTTP(scndRec, scndReq)

	var refresh models.UserConnect
	err := json.Unmarshal([]byte(scndRec.Body.String()), &refresh)
	if err != nil {
		log.Fatal("Bad output: ", err.Error())
		t.Fail()
	}

	assert.Equal(t, scndRec.Code, 200)

	utils.CleanUser(user.ID, user.Token, t, router)
	db.CloseDB()
}

// Asserts a token can expirate
func TestTokenExpiracy(t *testing.T) {
	db.InitDB()
	var router *gin.Engine = routes.SetupRouter()

	os.Setenv("TOKEN_VALIDITY_MINUTES", "0")

	var user models.UserCreate = utils.CreateUser("Tom", "qwerty1234", t, router)
	time.Sleep(1 * time.Second)

	var url string = "/v1/refresh/token"
	var bearer = "Bearer " + user.Token

	scndRec := httptest.NewRecorder()
	scndReq, _ := http.NewRequest("POST", url, nil)
	scndReq.Header.Add("Content-Type", "application/json")
	scndReq.Header.Add("Authorization", bearer)

	router.ServeHTTP(scndRec, scndReq)

	var message Message
	err := json.Unmarshal([]byte(scndRec.Body.String()), &message)
	if err != nil {
		log.Fatal("Bad output: ", err.Error())
		t.Fail()
	}

	assert.Equal(t, scndRec.Code, 401)
	assert.Equal(t, message.Message, "Token has expired")

	os.Setenv("TOKEN_VALIDITY_MINUTES", "15")

	var connect models.UserConnect = utils.ConnectUser("Tom", "qwerty1234", t, router)

	utils.CleanUser(user.ID, connect.Token, t, router)
	db.CloseDB()
}

// Asserts a token has a refresh limit
func TestTokenRefreshLimit(t *testing.T) {
	db.InitDB()
	var router *gin.Engine = routes.SetupRouter()

	os.Setenv("TOKEN_LIMIT_HOURS", "0")

	var user models.UserCreate = utils.CreateUser("Tom", "qwerty1234", t, router)
	time.Sleep(1 * time.Second)

	var url string = "/v1/refresh/token"
	var bearer = "Bearer " + user.Token
	scndRec := httptest.NewRecorder()
	scndReq, _ := http.NewRequest("POST", url, nil)
	scndReq.Header.Add("Content-Type", "application/json")
	scndReq.Header.Add("Authorization", bearer)

	router.ServeHTTP(scndRec, scndReq)

	var message Message
	err := json.Unmarshal([]byte(scndRec.Body.String()), &message)
	if err != nil {
		log.Fatal("Bad output: ", err.Error())
		t.Fail()
	}

	assert.Equal(t, scndRec.Code, 401)
	assert.Equal(t, message.Message, "Token has expired and cannot be refreshed, please reconnect")

	os.Setenv("TOKEN_LIMIT_HOURS", "24")

	var connect models.UserConnect = utils.ConnectUser("Tom", "qwerty1234", t, router)

	utils.CleanUser(user.ID, connect.Token, t, router)
	db.CloseDB()
}

// Asserts a token with a modified signature is not valid
func TestTokenSignature(t *testing.T) {
	db.InitDB()
	var router *gin.Engine = routes.SetupRouter()

	var user models.UserCreate = utils.CreateUser("Tom", "qwerty1234", t, router)

	splittedToken := strings.Split(user.Token, ".")
	if len(splittedToken) != 3 {
		log.Fatal("Bad token.")
		t.Fail()
	}

	header := splittedToken[0]
	payload := splittedToken[1]
	signature := splittedToken[2]
	signature = signature + "modif"
	modifiedToken := header + "." + payload + "." + signature

	var url string = "/v1/user/" + strconv.Itoa(user.ID)
	var bearer = "Bearer " + modifiedToken
	scndRec := httptest.NewRecorder()
	scndReq, _ := http.NewRequest("GET", url, nil)
	scndReq.Header.Add("Content-Type", "application/json")
	scndReq.Header.Add("Authorization", bearer)

	router.ServeHTTP(scndRec, scndReq)

	var message Message
	err := json.Unmarshal([]byte(scndRec.Body.String()), &message)
	if err != nil {
		log.Fatal("Bad output: ", err.Error())
		t.Fail()
	}

	assert.Equal(t, scndRec.Code, 401)
	assert.Equal(t, message.Message, "Bad signature")

	var connect models.UserConnect = utils.ConnectUser("Tom", "qwerty1234", t, router)
	utils.CleanUser(user.ID, connect.Token, t, router)
	db.CloseDB()
}

// Asserts a token with a modified payload is not valid
func TestTokenPayload(t *testing.T) {
	db.InitDB()
	var router *gin.Engine = routes.SetupRouter()

	var user models.UserCreate = utils.CreateUser("Tom", "qwerty1234", t, router)

	splittedToken := strings.Split(user.Token, ".")
	if len(splittedToken) != 3 {
		log.Fatal("Bad token.")
		t.Fail()
	}

	header := splittedToken[0]
	payload := splittedToken[1]
	signature := splittedToken[2]

	decPayloadByte, err := base64.RawURLEncoding.DecodeString(payload)

	var payloadObj models.JwtPayload
	err = json.Unmarshal(decPayloadByte, &payloadObj)
	if err != nil {
		log.Fatal("Bad output: ", err.Error())
		t.Fail()
	}
	payloadObj.Iat = 123456
	payloadObj.Exp = 123456
	jsonPayload, _ := json.Marshal(payloadObj)
	newPayload := base64.RawURLEncoding.EncodeToString([]byte(string(jsonPayload)))

	modifiedToken := header + "." + newPayload + "." + signature

	var url string = "/v1/user/" + strconv.Itoa(user.ID)
	var bearer = "Bearer " + modifiedToken
	scndRec := httptest.NewRecorder()
	scndReq, _ := http.NewRequest("GET", url, nil)
	scndReq.Header.Add("Content-Type", "application/json")
	scndReq.Header.Add("Authorization", bearer)

	router.ServeHTTP(scndRec, scndReq)

	var message Message
	err = json.Unmarshal([]byte(scndRec.Body.String()), &message)
	if err != nil {
		log.Fatal("Bad output: ", err.Error())
		t.Fail()
	}

	assert.Equal(t, scndRec.Code, 401)
	assert.Equal(t, message.Message, "Bad signature")

	var connect models.UserConnect = utils.ConnectUser("Tom", "qwerty1234", t, router)
	utils.CleanUser(user.ID, connect.Token, t, router)
	db.CloseDB()
}

// Asserts a user cannot access that does not belong to him
func TestRessourceAccess(t *testing.T) {
	db.InitDB()
	var router *gin.Engine = routes.SetupRouter()

	var user models.UserCreate = utils.CreateUser("Tom", "qwerty1234", t, router)

	var url string = "/v1/user/" + strconv.Itoa(user.ID+1)
	var bearer = "Bearer " + user.Token
	scndRec := httptest.NewRecorder()
	scndReq, _ := http.NewRequest("GET", url, nil)
	scndReq.Header.Add("Content-Type", "application/json")
	scndReq.Header.Add("Authorization", bearer)

	router.ServeHTTP(scndRec, scndReq)

	var message Message
	err := json.Unmarshal([]byte(scndRec.Body.String()), &message)
	if err != nil {
		log.Fatal("Bad output: ", err.Error())
		t.Fail()
	}

	assert.Equal(t, scndRec.Code, 403)
	assert.Equal(t, message.Message, "Forbidden.")

	utils.CleanUser(user.ID, user.Token, t, router)
	db.CloseDB()
}
