package rest

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"gitlab.com/sfz.aalen/hackwerk/dotinder/entity"
)

func (s *Suite) TestRegister() {
	userJson := `{"username": "test", "password": "test"}`
	req := httptest.NewRequest(http.MethodPost, "/api/users", strings.NewReader(userJson))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := s.e.NewContext(req, rec)

	// Assertions
	if assert.NoError(s.T(), s.restBoundary.registerUser(c)) {
		assert.Equal(s.T(), http.StatusCreated, rec.Code)

		var createdUser entity.PasswordUser
		json.NewDecoder(rec.Body).Decode(&createdUser)

		assert.Equal(s.T(), "test", createdUser.Username)
	}
}

func (s *Suite) TestLoginWhenUserDoesNotExist() {
	userJson := `{"username": "doesnotexist", "password": "doesnotexist"}`
	req := httptest.NewRequest(http.MethodPost, "/api/login", strings.NewReader(userJson))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := s.e.NewContext(req, rec)

	// Assertions
	assert.Error(s.T(), s.restBoundary.login(c))
}

func (s *Suite) TestLogin() {
	// Setup
	userJson := `{"username": "test", "password": "test"}`
	registerReq := httptest.NewRequest(http.MethodPost, "/api/users", strings.NewReader(userJson))
	registerReq.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	registerRec := httptest.NewRecorder()
	registerC := s.e.NewContext(registerReq, registerRec)
	assert.NoError(s.T(), s.restBoundary.registerUser(registerC))

	req := httptest.NewRequest(http.MethodPost, "/api/login", strings.NewReader(userJson))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := s.e.NewContext(req, rec)

	// Assertions
	if assert.NoError(s.T(), s.restBoundary.login(c)) {
		assert.Equal(s.T(), http.StatusOK, rec.Code)

		loginResponse := LoginResponse{}
		err := json.NewDecoder(rec.Body).Decode(&loginResponse)
		if err != nil {
			s.T().Fatal(err)
		}

		assert.NotEmpty(s.T(), loginResponse.Jwt)
	}
}

