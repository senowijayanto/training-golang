package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"training-golang/session-3-unit-test/handler"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetHelloMessage(t *testing.T) {
	t.Run("Positive case - correct message", func(t *testing.T) {
		expectedOutput := "Hello from Gin!"
		actualOutput := handler.GetHelloMessage()
		require.Equal(t, expectedOutput, actualOutput, "The message should be '%s'", expectedOutput)
	})
}

func TestRootHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.Default()
	router.GET("/", handler.RootHandler)

	// Create a new HTTP Request
	req, _ := http.NewRequest("GET", "/", nil)

	// Create a responseRecorder to record the response
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	expectedBody := `{"message":"Hello from Gin!"}`
	assert.JSONEq(t, expectedBody, w.Body.String())
}

type JsonRequest struct {
	Message string `json:"message"`
}

func TestPostHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Setup router
	r := gin.Default()
	r.POST("/", handler.PostHandler)

	t.Run("Positive Case", func(t *testing.T) {
		// Persiapan data JSON
		requestBody := JsonRequest{Message: "Hello from test!"}
		requestBodyBytes, _ := json.Marshal(requestBody)

		// Buat permintaan HTTP POST
		req, _ := http.NewRequest("POST", "/", bytes.NewBuffer(requestBodyBytes))
		req.Header.Set("Content-Type", "application/json")

		// Buat ResponseRecorder untuk merekam response
		w := httptest.NewRecorder()

		// Lakukan permintaan
		r.ServeHTTP(w, req)

		// Periksa status code
		assert.Equal(t, http.StatusOK, w.Code)

		// Periksa body response
		expectedBody := `{"message": "Hello from test!"}`
		assert.JSONEq(t, expectedBody, w.Body.String())
	})

	t.Run("Negative Case - EOF Error", func(t *testing.T) {
		// Persiapan data JSON yang salah
		requestBody := ""
		requestBodyBytes := []byte(requestBody)

		// Buat permintaan HTTP POST
		req, _ := http.NewRequest("POST", "/", bytes.NewBuffer(requestBodyBytes))
		req.Header.Set("Content-Type", "application/json")

		// Buat ResponseRecorder untuk merekam response
		w := httptest.NewRecorder()

		// Lakukan permintaan
		r.ServeHTTP(w, req)

		// Periksa status code
		assert.Equal(t, http.StatusBadRequest, w.Code)

		// Periksa body response
		assert.Contains(t, w.Body.String(), "{\"error\":\"EOF\"}")
	})
}
