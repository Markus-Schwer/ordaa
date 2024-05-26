package rest

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"gitlab.com/sfz.aalen/hackwerk/dotinder/entity"
)

func (s *Suite) TestGetAllMenus() {
	createMenuItems := []entity.MenuItem{}
	createMenu := &entity.Menu{Name: "testMenu", Url: "https://sangam-aalen.de", Items: createMenuItems}
	_, err := s.repo.CreateMenu(nil, createMenu)
	if err != nil {
		s.T().Fatal(err)
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
	createdMenu, err := s.repo.CreateMenu(nil, &entity.Menu{Name: "testMenu", Url: "https://sangam-aalen.de", Items: []entity.MenuItem{}})
	if err != nil {
		s.T().Fatal(err)
	}

	req := httptest.NewRequest(http.MethodPost, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := s.e.NewContext(req, rec)
	c.SetPath("/api/menus/:uuid")
	c.SetParamNames("uuid")
	c.SetParamValues(createdMenu.Uuid.String())

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
