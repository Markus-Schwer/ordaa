package rest

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"gitlab.com/sfz.aalen/hackwerk/dotinder/boundary/auth"
	"gitlab.com/sfz.aalen/hackwerk/dotinder/entity"
)

func TestGetAllMenus(t *testing.T) {
	// Setup
	ctx := context.Background()
	e := echo.New()

	repo := entity.NewMockRepository()
	authService := auth.NewAuthService(ctx, repo)
	restBoundary := NewRestBoundary(ctx, repo, authService)

	createMenuItems := []entity.MenuItem{}
	createMenu := &entity.Menu{Name: "testMenu", Url: "https://sangam-aalen.de", Items: createMenuItems}
	_, err := repo.CreateMenu(nil, createMenu)
	if err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest(http.MethodPost, "/api/menus", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Assertions
	if assert.NoError(t, restBoundary.allMenus(c)) {
		menus := []entity.Menu{}
		err := json.NewDecoder(rec.Body).Decode(&menus)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, 1, len(menus))
	}
}

func TestGetAllMenusEmpty(t *testing.T) {
	// Setup
	ctx := context.Background()
	e := echo.New()

	repo := entity.NewMockRepository()
	authService := auth.NewAuthService(ctx, repo)
	restBoundary := NewRestBoundary(ctx, repo, authService)

	req := httptest.NewRequest(http.MethodPost, "/api/menus", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Assertions
	if assert.NoError(t, restBoundary.allMenus(c)) {
		menus := []entity.Menu{}
		err := json.NewDecoder(rec.Body).Decode(&menus)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, []entity.Menu{}, menus)
	}
}

func TestGetMenu(t *testing.T) {
	// Setup
	ctx := context.Background()
	e := echo.New()

	repo := entity.NewMockRepository()
	authService := auth.NewAuthService(ctx, repo)
	restBoundary := NewRestBoundary(ctx, repo, authService)

	createdMenu, err := repo.CreateMenu(nil, &entity.Menu{Name: "testMenu", Url: "https://sangam-aalen.de", Items: []entity.MenuItem{}})
	if err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest(http.MethodPost, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/api/menus/:uuid")
	c.SetParamNames("uuid")
	c.SetParamValues(createdMenu.Uuid.String())

	// Assertions
	if assert.NoError(t, restBoundary.getMenu(c)) {
		var menu entity.Menu
		err := json.NewDecoder(rec.Body).Decode(&menu)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, "testMenu", menu.Name)
		assert.Equal(t, "https://sangam-aalen.de", menu.Url)
		assert.Equal(t, []entity.MenuItem{}, menu.Items)
	}
}
