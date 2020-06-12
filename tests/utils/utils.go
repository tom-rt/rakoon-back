package utils

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"rakoon/rakoon-back/models"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
	"gopkg.in/go-playground/assert.v1"
)

// CreateUser cleans a user created just for the test
func CreateUser(name string, password string, t *testing.T, router *gin.Engine) models.UserCreate {
	var jsonStr = []byte(`{"name":"` + name + `", "password": "` + password + `"}`)
	rec := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/v1/user", bytes.NewBuffer(jsonStr))
	req.Header.Add("Content-Type", "application/json")
	router.ServeHTTP(rec, req)
	assert.Equal(t, 201, rec.Code)

	var user models.UserCreate
	err := json.Unmarshal([]byte(rec.Body.String()), &user)
	if err != nil {
		log.Fatal("Bad output", err.Error())
		t.Fail()
	}

	return user
}

// CleanUser cleans a user created just for the test
func CleanUser(ID int, token string, t *testing.T, router *gin.Engine) {
	rec := httptest.NewRecorder()
	var url string = "/v1/user/" + strconv.Itoa(ID)
	var bearer = "Bearer " + token
	req, _ := http.NewRequest("DELETE", url, nil)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", bearer)
	router.ServeHTTP(rec, req)
	assert.Equal(t, 200, rec.Code)
}

// ConnectUser connects a user
func ConnectUser(name string, password string, t *testing.T, router *gin.Engine) models.UserConnect {
	var jsonStr = []byte(`{"name":"` + name + `", "password": "` + password + `"}`)
	rec := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/v1/user/connect", bytes.NewBuffer(jsonStr))
	req.Header.Add("Content-Type", "application/json")
	router.ServeHTTP(rec, req)
	assert.Equal(t, 200, rec.Code)

	var user models.UserConnect
	err := json.Unmarshal([]byte(rec.Body.String()), &user)
	if err != nil {
		log.Fatal("Bad output", err.Error())
		t.Fail()
	}

	return user
}
