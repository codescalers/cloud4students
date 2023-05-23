// Package app for c4s backend app
package app

import (
	"context"
	"fmt"
	"io"
	"net/http/httptest"
	"os"
	"path/filepath"

	"testing"

	"github.com/codescalers/cloud4students/internal"
	"github.com/codescalers/cloud4students/middlewares"
	"github.com/codescalers/cloud4students/models"
	"github.com/stretchr/testify/assert"
)

type authHandlerConfig struct {
	unAuthHandlerConfig
	userID string
	token  string
	config internal.Configuration
	db     models.DB
}

type unAuthHandlerConfig struct {
	body        io.Reader
	handlerFunc Handler
	api         string
}

// SetUp sets the needed configuration for testing
func SetUp(t testing.TB) *App {
	dir := t.TempDir()

	configPath := filepath.Join(dir, "config.json")
	dbPath := filepath.Join(dir, "testing.db")

	config := fmt.Sprintf(`
{
	"server": {
		"host": "localhost",
		"port": ":3000",
		"redisHost": "localhost",
		"redisPort": "6379",
		"redisPass": ""		
	},
	"mailSender": {
      "email": "email",
      "sendgrid_key": "my sendgrid_key",
      "timeout": 60 
    },
    "account": {
      "mnemonics": "winner giant reward damage expose pulse recipe manual brand volcano dry avoid",
			"network": "dev"
    },
	"token": {
      "secret": "secret",
      "timeout": 10
    },
	"database": {
      "file": "%s"
    },
	"version": "v1"
}
	`, dbPath)

	err := os.WriteFile(configPath, []byte(config), 0644)
	assert.NoError(t, err)

	app, err := NewApp(context.Background(), configPath)
	assert.NoError(t, err)

	return app
}

func unAuthorizedHandler(req unAuthHandlerConfig) (response *httptest.ResponseRecorder) {
	request := httptest.NewRequest("GET", req.api, req.body)
	response = httptest.NewRecorder()

	WrapFunc(req.handlerFunc).ServeHTTP(response, request)
	return
}

func authorizedNoMiddlewareHandler(req authHandlerConfig) (response *httptest.ResponseRecorder) {
	request := httptest.NewRequest("GET", req.api, req.body)
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", req.token))
	ctx := context.WithValue(request.Context(), middlewares.UserIDKey("UserID"), req.userID)
	newRequest := request.WithContext(ctx)
	response = httptest.NewRecorder()

	WrapFunc(req.handlerFunc).ServeHTTP(response, newRequest)
	return
}

func authorizedHandler(req authHandlerConfig) (response *httptest.ResponseRecorder) {
	request := httptest.NewRequest("GET", req.api, req.body)
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", req.token))
	ctx := context.WithValue(request.Context(), middlewares.UserIDKey("UserID"), req.userID)
	newRequest := request.WithContext(ctx)
	response = httptest.NewRecorder()

	handler := WrapFunc(req.handlerFunc)
	handlerWithAuth := middlewares.Authorization(req.db, req.config.Token.Secret, req.config.Token.Timeout)(handler)
	handlerWithAuth.ServeHTTP(response, newRequest)
	return
}

func adminHandler(req authHandlerConfig) (response *httptest.ResponseRecorder) {
	request := httptest.NewRequest("GET", req.api, req.body)
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", req.token))
	ctx := context.WithValue(request.Context(), middlewares.UserIDKey("UserID"), req.userID)
	newRequest := request.WithContext(ctx)
	response = httptest.NewRecorder()

	handler := WrapFunc(req.handlerFunc)
	handlerWithAdmin := middlewares.AdminAccess(req.db)(handler)
	handlerWithAuth := middlewares.Authorization(req.db, req.config.Token.Secret, req.config.Token.Timeout)(handlerWithAdmin)
	handlerWithAuth.ServeHTTP(response, newRequest)
	return
}
