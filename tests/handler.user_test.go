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

	// "rakoon/tests/utils/utils"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
	"gopkg.in/go-playground/assert.v1"
)

// Possibility to use only one router for every tests
// var router *gin.Engine
// func TestMain(m *testing.M) {
// 	db.InitDB()
// 	// router = routes.SetupRouter()
// 	code := m.Run()
// 	os.Exit(code)
// }

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
	scndRec := httptest.NewRecorder()
	scndReq, _ := http.NewRequest("POST", "/v1/user", bytes.NewBuffer(scndStr))
	scndReq.Header.Add("Content-Type", "application/json")

	router.ServeHTTP(scndRec, scndReq)
	assert.Equal(t, 409, scndRec.Code)

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
	scndRec := httptest.NewRecorder()
	scndReq, _ := http.NewRequest("POST", "/v1/user/connect", bytes.NewBuffer(connectStr))
	scndReq.Header.Add("Content-Type", "application/json")

	router.ServeHTTP(scndRec, scndReq)

	assert.Equal(t, 404, scndRec.Code)

	utils.CleanUser(user.ID, user.Token, t, router)
	db.CloseDB()
}

// Asserts a non existinig user can't connect
func TestUserNameConnect(t *testing.T) {
	db.InitDB()
	var router *gin.Engine = routes.SetupRouter()

	var user models.UserCreate = utils.CreateUser("Tom", "qwerty1234", t, router)

	var connectStr = []byte(`{"name":"DoesNotExist", "password": "pwdd"}`)
	scndRec := httptest.NewRecorder()
	scndReq, _ := http.NewRequest("POST", "/v1/user/connect", bytes.NewBuffer(connectStr))
	scndReq.Header.Add("Content-Type", "application/json")

	router.ServeHTTP(scndRec, scndReq)

	assert.Equal(t, 404, scndRec.Code)

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

	scndRec := httptest.NewRecorder()
	scndReq, _ := http.NewRequest("PUT", url, nil)
	scndReq.Header.Add("Content-Type", "application/json")
	scndReq.Header.Add("Authorization", bearer)

	router.ServeHTTP(scndRec, scndReq)
	assert.Equal(t, 200, scndRec.Code)

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

	scndRec := httptest.NewRecorder()
	scndReq, _ := http.NewRequest("GET", url, nil)
	scndReq.Header.Add("Content-Type", "application/json")
	scndReq.Header.Add("Authorization", bearer)

	router.ServeHTTP(scndRec, scndReq)

	var get models.UserPublic
	err := json.Unmarshal([]byte(scndRec.Body.String()), &get)
	if err != nil {
		log.Fatal("Bad output", err.Error())
		t.Fail()
	}

	assert.Equal(t, scndRec.Code, 200)
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
	scndRec := httptest.NewRecorder()
	scndReq, _ := http.NewRequest("PUT", url, bytes.NewBuffer(updateJSONStr))
	scndReq.Header.Add("Content-Type", "application/json")
	scndReq.Header.Add("Authorization", bearer)
	router.ServeHTTP(scndRec, scndReq)
	assert.Equal(t, scndRec.Code, 200)

	thirdRec := httptest.NewRecorder()
	url = "/v1/user/" + strconv.Itoa(user.ID)
	thirdReq, _ := http.NewRequest("GET", url, nil)
	thirdReq.Header.Add("Content-Type", "application/json")
	thirdReq.Header.Add("Authorization", bearer)
	router.ServeHTTP(thirdRec, thirdReq)
	var get models.UserPublic
	err := json.Unmarshal([]byte(thirdRec.Body.String()), &get)
	if err != nil {
		log.Fatal("Bad output", err.Error())
		t.Fail()
	}

	assert.Equal(t, thirdRec.Code, 200)
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
	scndRec := httptest.NewRecorder()
	scndReq, _ := http.NewRequest("PUT", url, nil)
	scndReq.Header.Add("Content-Type", "application/json")
	scndReq.Header.Add("Authorization", bearer)
	router.ServeHTTP(scndRec, scndReq)

	assert.Equal(t, scndRec.Code, 200)

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
	scndRec := httptest.NewRecorder()
	scndReq, _ := http.NewRequest("PUT", url, bytes.NewBuffer(passwordJSONStr))
	scndReq.Header.Add("Content-Type", "application/json")
	scndReq.Header.Add("Authorization", bearer)
	router.ServeHTTP(scndRec, scndReq)

	assert.Equal(t, scndRec.Code, 200)

	// A user has to reconnect after a password change
	var uConnect models.UserConnect = utils.ConnectUser("foo", "1234qwerty", t, router)

	utils.CleanUser(user.ID, uConnect.Token, t, router)
	db.CloseDB()
}
