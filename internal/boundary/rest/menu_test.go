package rest

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/gofrs/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/Markus-Schwer/ordaa/internal/entity"
	"gorm.io/gorm"
)

func (s *Suite) TestGetAllMenus() {
	createMenuItems := []entity.MenuItem{}
	menuUuid := uuid.Must(uuid.NewV4())
	menu := entity.Menu{Uuid: &menuUuid, Name: "testMenu", Url: "https://sangam-aalen.de", Items: createMenuItems}

	s.repo.GetAllMenusFunc = func(tx *gorm.DB) ([]entity.Menu, error) {
		return []entity.Menu{menu}, nil
	}

	req := httptest.NewRequest(http.MethodPost, "/api/menus", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := s.e.NewContext(req, rec)

	// Assertions
	if assert.NoError(s.T(), s.restBoundary.allMenus(c)) {
		menus := []entity.Menu{}
		err := json.NewDecoder(rec.Body).Decode(&menus)
		if err != nil {
			s.T().Fatal(err)
		}

		assert.Equal(s.T(), 1, len(menus))
	}
}

func (s *Suite) TestGetAllMenusEmpty() {
	s.repo.GetAllMenusFunc = func(tx *gorm.DB) ([]entity.Menu, error) {
		return []entity.Menu{}, nil
	}

	req := httptest.NewRequest(http.MethodPost, "/api/menus", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := s.e.NewContext(req, rec)

	// Assertions
	if assert.NoError(s.T(), s.restBoundary.allMenus(c)) {
		menus := []entity.Menu{}
		err := json.NewDecoder(rec.Body).Decode(&menus)
		if err != nil {
			s.T().Fatal(err)
		}

		assert.Equal(s.T(), []entity.Menu{}, menus)
	}
}

func (s *Suite) TestGetMenu() {
	menuUuid := uuid.Must(uuid.NewV4())
	menu := &entity.Menu{Uuid: &menuUuid, Name: "testMenu", Url: "https://sangam-aalen.de", Items: []entity.MenuItem{}}

	s.repo.GetMenuFunc = func(tx *gorm.DB, uuidMoqParam *uuid.UUID) (*entity.Menu, error) {
		if *uuidMoqParam == menuUuid {
			return menu, nil
		}

		return nil, entity.ErrMenuNotFound
	}

	req := httptest.NewRequest(http.MethodPost, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := s.e.NewContext(req, rec)
	c.SetPath("/api/menus/:uuid")
	c.SetParamNames("uuid")
	c.SetParamValues(menuUuid.String())

	// Assertions
	if assert.NoError(s.T(), s.restBoundary.getMenu(c)) {
		var menu entity.Menu
		err := json.NewDecoder(rec.Body).Decode(&menu)
		if err != nil {
			s.T().Fatal(err)
		}

		assert.Equal(s.T(), "testMenu", menu.Name)
		assert.Equal(s.T(), "https://sangam-aalen.de", menu.Url)
		assert.Equal(s.T(), []entity.MenuItem{}, menu.Items)
	}
}
