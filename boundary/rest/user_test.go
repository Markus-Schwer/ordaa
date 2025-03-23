package rest

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/gofrs/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"gitlab.com/sfz.aalen/hackwerk/dotinder/crypto"
	"gitlab.com/sfz.aalen/hackwerk/dotinder/entity"
	"gorm.io/gorm"
)

func (s *Suite) TestRegister() {
	s.repo.CreateUserFunc = func(tx *gorm.DB, user *entity.User) (*entity.User, error) {
		userUuid := uuid.Must(uuid.NewV4())
		user.Uuid = &userUuid
		return user, nil
	}

	s.repo.CreatePasswordUserFunc = func(tx *gorm.DB, passwordUser *entity.PasswordUser) (*entity.PasswordUser, error) {
		userUuid := uuid.Must(uuid.NewV4())
		passwordUser.Uuid = &userUuid
		return passwordUser, nil
	}

	userJson := `{"username": "test", "password": "test"}`
	req := httptest.NewRequest(http.MethodPost, "/api/users", strings.NewReader(userJson))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := s.e.NewContext(req, rec)

	// Assertions
	if assert.NoError(s.T(), s.restBoundary.registerUser(c)) {
		assert.Equal(s.T(), http.StatusCreated, rec.Code)

		var createdUser entity.User
		json.NewDecoder(rec.Body).Decode(&createdUser)

		assert.Equal(s.T(), "test", createdUser.Name)
	}
}

func (s *Suite) TestLoginWhenUserDoesNotExist() {
	s.repo.FindPasswordUserFunc = func(tx *gorm.DB, username string) (*entity.PasswordUser, error) {
		return nil, entity.ErrUserNotFound
	}

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
	s.repo.CreateUserFunc = func(tx *gorm.DB, user *entity.User) (*entity.User, error) {
		userUuid := uuid.Must(uuid.NewV4())
		user.Uuid = &userUuid
		return user, nil
	}

	s.repo.CreatePasswordUserFunc = func(tx *gorm.DB, passwordUser *entity.PasswordUser) (*entity.PasswordUser, error) {
		userUuid := uuid.Must(uuid.NewV4())
		passwordUser.Uuid = &userUuid
		return passwordUser, nil
	}

	s.repo.FindPasswordUserFunc = func(tx *gorm.DB, username string) (*entity.PasswordUser, error) {
		userUuid := uuid.Must(uuid.NewV4())
		if username == "test" {
			hash, err := crypto.GeneratePasswordHash("test")
			if err != nil {
				return nil, err
			}

			return &entity.PasswordUser{
				Uuid:     &userUuid,
				UserUuid: &userUuid,
				Username: "test",
				Password: hash,
			}, nil
		}
		return nil, entity.ErrUserNotFound
	}

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
