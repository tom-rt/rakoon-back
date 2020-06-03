package test

import (
	"net/http"
	"net/http/httptest"
	"rakoon/rakoon-back/routes"
	"testing"

	"gopkg.in/go-playground/assert.v1"
)

func UserHandlerTest(t *testing.T) {
	router := routes.SetupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "pong", w.Body.String())
}
