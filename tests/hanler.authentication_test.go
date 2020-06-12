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
	"testing"

	"github.com/gin-gonic/gin"
	"gopkg.in/go-playground/assert.v1"
)

// Asserts a connected user can refresh his token.
func TestRefreshToken(t *testing.T) {
	db.InitDB()
	var router *gin.Engine = routes.SetupRouter()

	var jsonStr = []byte(`{"name":"TestUser", "password": "qwerty1234"}`)
	rec := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/v1/user", bytes.NewBuffer(jsonStr))
	req.Header.Add("Content-Type", "application/json")
	router.ServeHTTP(rec, req)

	assert.Equal(t, 201, rec.Code)
	var ret models.UserCreate
	err := json.Unmarshal([]byte(rec.Body.String()), &ret)

	if err != nil {
		log.Fatal("Bad output", err.Error())
		t.Fail()
	}

	// CleanUser(ret.ID, ret.Token, router)
	db.CloseDB()
}

// Asserts a token can expirate
func TestTokenExpiracy(t *testing.T) {
	db.InitDB()
	var router *gin.Engine = routes.SetupRouter()

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

	assert.Equal(t, 409, scndRec.Code)

	var ret models.UserCreate
	err := json.Unmarshal([]byte(firstRec.Body.String()), &ret)

	if err != nil {
		log.Fatal("Bad output", err.Error())
		t.Fail()
	}

	// CleanUser(ret.ID, ret.Token, router)
	db.CloseDB()
}
