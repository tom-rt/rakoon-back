package test

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"

	"rakoon/rakoon-back/db"
	"rakoon/rakoon-back/models"
	"rakoon/rakoon-back/routes"
	"rakoon/rakoon-back/tests/utils"

	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
	"gopkg.in/go-playground/assert.v1"
)

// Asserts user creation works.
func TestCreateUser(t *testing.T) {
	db.InitDB()
	var router *gin.Engine = routes.SetupRouter()

	var user models.UserCreate = utils.CreateUser("TestUser", "qwerty1234", t, router)

	utils.CleanUser(user.ID, user.Token, t, router)
	db.CloseDB()
}

// Asserts usernames are unique.
func TestUserDuplicate(t *testing.T) {
	db.InitDB()
	var router *gin.Engine = routes.SetupRouter()

	var user models.UserCreate = utils.CreateUser("John", "qwerty1234", t, router)

	var scndStr = []byte(`{"name":"John", "password": "1234qwerty"}`)
	record := httptest.NewRecorder()
	request, _ := http.NewRequest("POST", "/v1/user", bytes.NewBuffer(scndStr))
	request.Header.Add("Content-Type", "application/json")

	router.ServeHTTP(record, request)
	assert.Equal(t, 409, record.Code)

	utils.CleanUser(user.ID, user.Token, t, router)
	db.CloseDB()
}

// Asserts a user can connect
func TestUserConnection(t *testing.T) {
	db.InitDB()
	var router *gin.Engine = routes.SetupRouter()

	var uCreate models.UserCreate = utils.CreateUser("Alf", "qwerty1234", t, router)

	var uConnect models.UserConnect = utils.ConnectUser("Alf", "qwerty1234", t, router)

	utils.CleanUser(uCreate.ID, uConnect.Token, t, router)
	db.CloseDB()
}

// Asserts the user password is properly checked
func TestUserPasswordConnect(t *testing.T) {
	db.InitDB()
	var router *gin.Engine = routes.SetupRouter()

	var user models.UserCreate = utils.CreateUser("Tom", "qwerty1234", t, router)

	var connectStr = []byte(`{"name":"Tom", "password": "not the same password"}`)
	record := httptest.NewRecorder()
	request, _ := http.NewRequest("POST", "/v1/user/connect", bytes.NewBuffer(connectStr))
	request.Header.Add("Content-Type", "application/json")

	router.ServeHTTP(record, request)

	assert.Equal(t, 404, record.Code)

	utils.CleanUser(user.ID, user.Token, t, router)
	db.CloseDB()
}

// Asserts a non existinig user can't connect
func TestUserNameConnect(t *testing.T) {
	db.InitDB()
	var router *gin.Engine = routes.SetupRouter()

	var user models.UserCreate = utils.CreateUser("Tom", "qwerty1234", t, router)

	var connectStr = []byte(`{"name":"DoesNotExist", "password": "pwdd"}`)
	record := httptest.NewRecorder()
	request, _ := http.NewRequest("POST", "/v1/user/connect", bytes.NewBuffer(connectStr))
	request.Header.Add("Content-Type", "application/json")

	router.ServeHTTP(record, request)

	assert.Equal(t, 404, record.Code)

	utils.CleanUser(user.ID, user.Token, t, router)
	db.CloseDB()
}

// Asserts a user can log out
func TestUserLogout(t *testing.T) {
	db.InitDB()
	var router *gin.Engine = routes.SetupRouter()

	var user models.UserCreate = utils.CreateUser("Jon", "qwerty1234", t, router)

	var url string = "/v1/user/" + strconv.Itoa(user.ID) + "/logout"
	var bearer = "Bearer " + user.Token

	record := httptest.NewRecorder()
	request, _ := http.NewRequest("PUT", url, nil)
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Authorization", bearer)

	router.ServeHTTP(record, request)
	assert.Equal(t, 200, record.Code)

	var uConnect models.UserConnect = utils.ConnectUser("Jon", "qwerty1234", t, router)

	utils.CleanUser(user.ID, uConnect.Token, t, router)
	db.CloseDB()
}

// Asserts you can get a user's data
func TestUserGet(t *testing.T) {
	db.InitDB()
	var router *gin.Engine = routes.SetupRouter()

	var user models.UserCreate = utils.CreateUser("Jean", "qwerty1234", t, router)

	var url string = "/v1/user/" + strconv.Itoa(user.ID)
	var bearer = "Bearer " + user.Token
	record := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", url, nil)
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Authorization", bearer)

	router.ServeHTTP(record, request)

	var get models.UserPublic
	err := json.Unmarshal([]byte(record.Body.String()), &get)
	if err != nil {
		log.Fatal("Bad output", err.Error())
		t.Fail()
	}

	assert.Equal(t, record.Code, 200)
	assert.Equal(t, get.Name, "Jean")
	assert.Equal(t, get.Reauth, false)

	utils.CleanUser(user.ID, user.Token, t, router)
	db.CloseDB()
}

// Asserts you can update a user's data
func TestUserUpdate(t *testing.T) {
	db.InitDB()
	var router *gin.Engine = routes.SetupRouter()

	var user models.UserCreate = utils.CreateUser("Tom", "qwerty1234", t, router)

	var updateJSONStr = []byte(`{"name":"Mike"}`)
	var url string = "/v1/user/" + strconv.Itoa(user.ID)
	var bearer = "Bearer " + user.Token
	record := httptest.NewRecorder()
	request, _ := http.NewRequest("PUT", url, bytes.NewBuffer(updateJSONStr))
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Authorization", bearer)
	router.ServeHTTP(record, request)
	assert.Equal(t, record.Code, 200)

	getRecorder := httptest.NewRecorder()
	url = "/v1/user/" + strconv.Itoa(user.ID)
	getRequest, _ := http.NewRequest("GET", url, nil)
	getRequest.Header.Add("Content-Type", "application/json")
	getRequest.Header.Add("Authorization", bearer)
	router.ServeHTTP(getRecorder, getRequest)
	var get models.UserPublic
	err := json.Unmarshal([]byte(getRecorder.Body.String()), &get)
	if err != nil {
		log.Fatal("Bad output", err.Error())
		t.Fail()
	}

	assert.Equal(t, getRecorder.Code, 200)
	assert.Equal(t, get.ID, user.ID)
	assert.Equal(t, get.Name, "Mike")
	assert.Equal(t, get.Reauth, false)

	utils.CleanUser(user.ID, user.Token, t, router)
	db.CloseDB()
}

// Asserts you can archive a user'
func TestUserArchive(t *testing.T) {
	db.InitDB()
	var router *gin.Engine = routes.SetupRouter()

	var user models.UserCreate = utils.CreateUser("Jean", "qwerty1234", t, router)

	var url string = "/v1/user/" + strconv.Itoa(user.ID) + "/archive"
	var bearer = "Bearer " + user.Token
	record := httptest.NewRecorder()
	request, _ := http.NewRequest("PUT", url, nil)
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Authorization", bearer)
	router.ServeHTTP(record, request)

	assert.Equal(t, record.Code, 200)

	// After an archive, the user must reconnect to be able to be deleted
	var uConnect models.UserConnect = utils.ConnectUser("Jean", "qwerty1234", t, router)

	utils.CleanUser(user.ID, uConnect.Token, t, router)
	db.CloseDB()
}

// Asserts you can change a user's password
func TestUserPasswordChange(t *testing.T) {
	db.InitDB()
	var router *gin.Engine = routes.SetupRouter()

	var user models.UserCreate = utils.CreateUser("foo", "qwerty1234", t, router)

	var passwordJSONStr = []byte(`{"password": "1234qwerty"}`)
	var url string = "/v1/user/" + strconv.Itoa(user.ID) + "/password"
	var bearer = "Bearer " + user.Token
	record := httptest.NewRecorder()
	request, _ := http.NewRequest("PUT", url, bytes.NewBuffer(passwordJSONStr))
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Authorization", bearer)
	router.ServeHTTP(record, request)

	assert.Equal(t, record.Code, 200)

	// A user has to reconnect after a password change
	var uConnect models.UserConnect = utils.ConnectUser("foo", "1234qwerty", t, router)

	utils.CleanUser(user.ID, uConnect.Token, t, router)
	db.CloseDB()
}
