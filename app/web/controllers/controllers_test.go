package controllers

import (
	"bytes"
	"html/template"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIndexController_Get(t *testing.T) {
	// Create a template with name "index.html"
	tmpl, err := template.New("index.html").Parse("Title: {{.Title}}")
	assert.NoError(t, err, "failed to create template")

	// Create a new IndexController
	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{Level: slog.LevelInfo}))
	indexCtrl := NewIndexController(logger)

	// Set up a ServeMux and register the controller router
	mux := http.NewServeMux()
	indexCtrl.Router(mux, tmpl)

	// Create a GET request to the controller route "/"
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	// Process the request
	mux.ServeHTTP(rec, req)

	// Check that the response status is 200 OK
	assert.Equal(t, http.StatusOK, rec.Code, "expected status code 200")

	// Check that the response body contains the expected title
	expected := "Title: NewIndexController"
	assert.Contains(t, rec.Body.String(), expected, "response body should contain expected title")
}

func TestIndexController_MethodNotAllowed(t *testing.T) {
	// Create a template with name "index.html"
	tmpl, err := template.New("index.html").Parse("Title: {{.Title}}")
	assert.NoError(t, err, "failed to create template")

	// Create a new IndexController
	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{Level: slog.LevelInfo}))
	indexCtrl := NewIndexController(logger)

	// Set up a ServeMux and register the controller router
	mux := http.NewServeMux()
	indexCtrl.Router(mux, tmpl)

	// Create a POST request to the controller route "/"
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	rec := httptest.NewRecorder()

	// Process the request
	mux.ServeHTTP(rec, req)

	// Check that the response status is 405 Method Not Allowed
	assert.Equal(t, http.StatusMethodNotAllowed, rec.Code, "expected status code 405")

	// Check that the response body contains the method not allowed message
	expected := "Method not allowed"
	assert.Contains(t, rec.Body.String(), expected, "response body should contain method not allowed message")
}
