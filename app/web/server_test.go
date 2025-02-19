package web

import (
	"bytes"
	"html/template"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/meesooqa/ttag/app/proc/mocks"
	"github.com/meesooqa/ttag/app/web/controllers"
)

func TestServer_ControllerRoute(t *testing.T) {
	// Create a template for index
	tmpl, err := template.New("index.html").Parse("Title: {{.Title}}")
	assert.NoError(t, err, "failed to create template")

	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{Level: slog.LevelInfo}))

	// Create an IndexController
	repo := &mocks.RepositoryMock{}
	indexCtrl := controllers.NewIndexController(logger, repo)

	// Initialize the server with the controller
	srv := NewServer(logger, []controllers.Controller{indexCtrl})
	// In tests, manually assign our template
	srv.templates = tmpl

	// Get the server router
	handler := srv.router()

	// Create a GET request to the controller route "/"
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	// Process the request
	handler.ServeHTTP(rec, req)

	// Check that the response status is 200 OK
	assert.Equal(t, http.StatusOK, rec.Code, "expected status code 200")

	// Check that the response body contains the expected title
	expected := "Title: NewIndexController"
	assert.Contains(t, rec.Body.String(), expected, "response body should contain expected title")
}

func TestServer_StaticFiles(t *testing.T) {
	// Create a temporary directory for static files
	staticDir := t.TempDir()

	// Create a test file in the directory
	fileContent := "Hello, static file!"
	fileName := "test.txt"
	err := os.WriteFile(staticDir+"/"+fileName, []byte(fileContent), 0644)
	assert.NoError(t, err, "failed to create static file")

	// Initialize the server without controllers
	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{Level: slog.LevelInfo}))
	srv := NewServer(logger, []controllers.Controller{})
	// Set the static files path
	srv.tplStaticLocation = staticDir

	// Get the server router
	handler := srv.router()

	// Create a GET request to the static file route
	req := httptest.NewRequest(http.MethodGet, "/static/"+fileName, nil)
	rec := httptest.NewRecorder()

	// Process the request
	handler.ServeHTTP(rec, req)

	// Check that the response status is 200 OK
	assert.Equal(t, http.StatusOK, rec.Code, "expected status code 200")

	// Read the response body using io.ReadAll
	body, err := io.ReadAll(rec.Body)
	assert.NoError(t, err, "failed to read response body")

	// Check that the response body contains the file content
	assert.Contains(t, string(body), fileContent, "response body should contain file content")
}
