package router_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"training-golang/session-3-unit-test/router"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestSetupRouter_RootHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Inisialisasi router
	r := gin.Default()
	router.SetupRouter(r)

	// Buat permintaan HTTP GET ke root ("/")
	req, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	// Lakukan permintaan
	r.ServeHTTP(w, req)

	// Periksa status code
	assert.Equal(t, http.StatusOK, w.Code)

	// Periksa body response
	expectedBody := `{"message": "Hello from Gin!"}`
	assert.JSONEq(t, expectedBody, w.Body.String())
}

func TestPostHandler_PostiveCase(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Inisialisasi router
	r := gin.Default()
	router.SetupRouter(r)

	// Persiapan data JSON
	requestBody := map[string]string{"message": "Test message"}
	requestBodyBytes, _ := json.Marshal(requestBody)

	// Buat permintaan HTTP POST dengan data JSON yang valid
	req, _ := http.NewRequest("POST", "/private/post", bytes.NewBuffer(requestBodyBytes))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "valid-token")

	w := httptest.NewRecorder()

	// Lakukan permintaan
	r.ServeHTTP(w, req)

	// Periksa status code
	assert.Equal(t, http.StatusOK, w.Code)

	// Periksa body response
	expectedBody := `{"message": "Test message"}`
	assert.JSONEq(t, expectedBody, w.Body.String())
}

func TestPostHandler_NegativeCase_BadRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Inisialisasi router
	r := gin.Default()
	router.SetupRouter(r)

	// Buat permintaan HTTP POST dengan data JSON yang tidak valid
	req, _ := http.NewRequest("POST", "/private/post", bytes.NewBufferString("{Invalid JSON}"))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "valid-token")

	w := httptest.NewRecorder()

	// Lakukan permintaan
	r.ServeHTTP(w, req)

	// Periksa status code
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// Periksa body response
	assert.Contains(t, w.Body.String(), "invalid character")
}

func TestPostHandler_NegativeCase_NoAuthHeader(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Inisialisasi router
	r := gin.Default()
	router.SetupRouter(r)

	// Buat permintaan HTTP POST tanpa header autentifikas
	req, _ := http.NewRequest("POST", "/private/post", nil)

	w := httptest.NewRecorder()

	// Lakukan permintaan
	r.ServeHTTP(w, req)

	// Periksa status code
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	// Periksa body response
	assert.Contains(t, w.Body.String(), "{\"error\":\"Authorization token required\"}")
}
