package test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"rakoon/rakoon-back/routes"
	"strings"
	"testing"

	"gopkg.in/go-playground/assert.v1"
)

func TestConnect(t *testing.T) {
	router := routes.SetupRouter()

	body := url.Values{}
	body.Add("name", "tom")
	body.Add("password", "qwerty")
	data := body.Encode()

	req, err := http.NewRequest("POST", "/v1/user", strings.NewReader(data))
	req.Header.Add("Content-Type", "application/json")

	if err != nil {
		t.Fail()
		return
	}

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	fmt.Println(w.Code)
	fmt.Println(w.Body)
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "pong", w.Body.String())
}
