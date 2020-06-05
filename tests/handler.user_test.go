package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"rakoon/rakoon-back/db"
	"rakoon/rakoon-back/routes"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
	"gopkg.in/go-playground/assert.v1"
)

type UserCreate struct {
	ID    int    `json:"id"`
	Token string `json:"token"`
}

var router *gin.Engine

func TestMain(m *testing.M) {
	db.InitDb()
	router = routes.SetupRouter()
	code := m.Run()
	os.Exit(code)
}

// Asserts user creation works.
func TestCreateUser(t *testing.T) {
	var jsonStr = []byte(`{"name":"TestUser", "password": "qwerty1234"}`)
	rec := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/v1/user", bytes.NewBuffer(jsonStr))
	req.Header.Add("Content-Type", "application/json")
	router.ServeHTTP(rec, req)

	assert.Equal(t, 201, rec.Code)
	var ret UserCreate
	err := json.Unmarshal([]byte(rec.Body.String()), &ret)

	if err != nil {
		log.Fatal("Bad output", err.Error())
		t.Fail()
	}

	cleanUser(ret.ID, ret.Token)
}

// Asserts usernames are unique.
func TestUserDuplicate(t *testing.T) {
	var firstStr = []byte(`{"name":"John", "password": "qwerty1234"}`)
	firstRec := httptest.NewRecorder()
	firstReq, _ := http.NewRequest("POST", "/v1/user", bytes.NewBuffer(firstStr))
	firstReq.Header.Add("Content-Type", "application/json")

	var scndStr = []byte(`{"name":"John", "password": "1234qwerty"}`)
	scndRec := httptest.NewRecorder()
	scndReq, _ := http.NewRequest("POST", "/v1/user", bytes.NewBuffer(scndStr))
	scndReq.Header.Add("Content-Type", "application/json")

	router.ServeHTTP(firstRec, firstReq)
	router.ServeHTTP(scndRec, scndReq)

	fmt.Println(scndRec.Body.String())
	assert.Equal(t, 409, scndRec.Code)

	var ret UserCreate
	err := json.Unmarshal([]byte(firstRec.Body.String()), &ret)

	if err != nil {
		log.Fatal("Bad output", err.Error())
		t.Fail()
	}

	cleanUser(ret.ID, ret.Token)
}

// Asserts a user can connect
func TestUserConnection(t *testing.T) {
	var jsonStr = []byte(`{"name":"Alfred", "password": "notsosafe"}`)

	firstRec := httptest.NewRecorder()
	firstReq, _ := http.NewRequest("POST", "/v1/user", bytes.NewBuffer(jsonStr))
	firstReq.Header.Add("Content-Type", "application/json")

	scndRec := httptest.NewRecorder()
	scndReq, _ := http.NewRequest("POST", "/v1/user/connect", bytes.NewBuffer(jsonStr))
	scndReq.Header.Add("Content-Type", "application/json")

	router.ServeHTTP(firstRec, firstReq)
	router.ServeHTTP(scndRec, scndReq)

	assert.Equal(t, 200, scndRec.Code)

	var ret UserCreate
	err := json.Unmarshal([]byte(firstRec.Body.String()), &ret)
	if err != nil {
		log.Fatal("Bad output", err.Error())
		t.Fail()
	}

	cleanUser(ret.ID, ret.Token)
}

func TestUserPassword(t *testing.T) {
	var createStr = []byte(`{"name":"Tom", "password": "pwdd"}`)
	firstRec := httptest.NewRecorder()
	firstReq, _ := http.NewRequest("POST", "/v1/user", bytes.NewBuffer(createStr))
	firstReq.Header.Add("Content-Type", "application/json")

	var connectStr = []byte(`{"name":"Tom", "password": "not the same password"}`)
	scndRec := httptest.NewRecorder()
	scndReq, _ := http.NewRequest("POST", "/v1/user/connect", bytes.NewBuffer(connectStr))
	scndReq.Header.Add("Content-Type", "application/json")

	router.ServeHTTP(firstRec, firstReq)
	router.ServeHTTP(scndRec, scndReq)

	assert.Equal(t, 404, scndRec.Code)

	var ret UserCreate
	err := json.Unmarshal([]byte(firstRec.Body.String()), &ret)
	if err != nil {
		log.Fatal("Bad output", err.Error())
		t.Fail()
	}

	cleanUser(ret.ID, ret.Token)
}

func TestUserName(t *testing.T) {
	var createStr = []byte(`{"name":"Tom", "password": "pwdd"}`)
	firstRec := httptest.NewRecorder()
	firstReq, _ := http.NewRequest("POST", "/v1/user", bytes.NewBuffer(createStr))
	firstReq.Header.Add("Content-Type", "application/json")

	var connectStr = []byte(`{"name":"DoesNotExist", "password": "pwdd"}`)
	scndRec := httptest.NewRecorder()
	scndReq, _ := http.NewRequest("POST", "/v1/user/connect", bytes.NewBuffer(connectStr))
	scndReq.Header.Add("Content-Type", "application/json")

	router.ServeHTTP(firstRec, firstReq)
	router.ServeHTTP(scndRec, scndReq)

	assert.Equal(t, 404, scndRec.Code)

	var ret UserCreate
	err := json.Unmarshal([]byte(firstRec.Body.String()), &ret)
	if err != nil {
		log.Fatal("Bad output", err.Error())
		t.Fail()
	}

	cleanUser(ret.ID, ret.Token)
}

// Asserts a user can log out
func TestUserLogout(t *testing.T) {
	var jsonStr = []byte(`{"name":"MonsieurMonsieur", "password": "asdfg"}`)

	firstRec := httptest.NewRecorder()
	firstReq, _ := http.NewRequest("POST", "/v1/user", bytes.NewBuffer(jsonStr))
	firstReq.Header.Add("Content-Type", "application/json")

	router.ServeHTTP(firstRec, firstReq)

	assert.Equal(t, 201, firstRec.Code)

	var ret UserCreate
	err := json.Unmarshal([]byte(firstRec.Body.String()), &ret)
	if err != nil {
		log.Fatal("Bad output", err.Error())
		t.Fail()
	}

	var url string = "/v1/user/" + strconv.Itoa(ret.ID) + "/logout"
	fmt.Println(url, ret.Token, ret.ID)
	var bearer = "Bearer " + ret.Token

	scndRec := httptest.NewRecorder()
	scndReq, _ := http.NewRequest("PUT", url, nil)
	scndReq.Header.Add("Content-Type", "application/json")
	scndReq.Header.Add("Authorization", bearer)

	router.ServeHTTP(scndRec, scndReq)
	assert.Equal(t, 200, scndRec.Code)

	cleanUser(ret.ID, ret.Token)
}

// Cleans a user created just for the test
func cleanUser(ID int, token string) {
	w := httptest.NewRecorder()
	var url string = "/v1/user/" + strconv.Itoa(ID)
	var bearer = "Bearer " + token
	req, _ := http.NewRequest("DELETE", url, nil)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", bearer)
	router.ServeHTTP(w, req)
}
