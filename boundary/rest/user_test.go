package rest

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"gitlab.com/sfz.aalen/hackwerk/dotinder/boundary/auth"
	"gitlab.com/sfz.aalen/hackwerk/dotinder/entity"
)

func TestRegister(t *testing.T) {
	// Setup
	userJson := `{"username": "test", "password": "test"}`
	ctx := context.Background()
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/users", strings.NewReader(userJson))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	repo := entity.NewMockRepository()
	authService := auth.NewAuthService(ctx, repo)
	restBoundary := NewRestBoundary(ctx, repo, authService)

	// Assertions
	if assert.NoError(t, restBoundary.registerUser(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)

		var createdUser entity.PasswordUser
		json.NewDecoder(rec.Body).Decode(&createdUser)

		assert.Equal(t, "test", createdUser.Username)
	}
}

func TestLoginWhenUserDoesNotExist(t *testing.T) {
	// Setup
	userJson := `{"username": "doesnotexist", "password": "doesnotexist"}`
	ctx := context.Background()
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/login", strings.NewReader(userJson))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	repo := entity.NewMockRepository()
	authService := auth.NewAuthService(ctx, repo)
	restBoundary := NewRestBoundary(ctx, repo, authService)

	// Assertions
	assert.Error(t, restBoundary.login(c))
}

func TestLogin(t *testing.T) {
	// Setup
	userJson := `{"username": "test", "password": "test"}`
	ctx := context.Background()
	e := echo.New()

	repo := entity.NewMockRepository()
	authService := auth.NewAuthService(ctx, repo)
	restBoundary := NewRestBoundary(ctx, repo, authService)

	registerReq := httptest.NewRequest(http.MethodPost, "/api/users", strings.NewReader(userJson))
	registerReq.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	registerRec := httptest.NewRecorder()
	registerC := e.NewContext(registerReq, registerRec)
	assert.NoError(t, restBoundary.registerUser(registerC))

	req := httptest.NewRequest(http.MethodPost, "/api/login", strings.NewReader(userJson))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Assertions
	if assert.NoError(t, restBoundary.login(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)

		loginResponse := LoginResponse{}
		err := json.NewDecoder(rec.Body).Decode(&loginResponse)
		if err != nil {
			t.Fatal(err)
		}

		assert.NotEmpty(t, loginResponse.Jwt)
	}
}

