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

func cleanUser(ID int, token string) {
	w := httptest.NewRecorder()
	var url string = "/v1/user/" + strconv.Itoa(ID)
	var bearer = "Bearer " + token
	req, _ := http.NewRequest("DELETE", url, nil)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", bearer)
	router.ServeHTTP(w, req)
}
