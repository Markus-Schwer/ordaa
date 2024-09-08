package rest

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"gitlab.com/sfz.aalen/hackwerk/dotinder/boundary/utils"
	"gitlab.com/sfz.aalen/hackwerk/dotinder/entity"
)

func (s *Suite) TestCreateOrder() {
	menu := entity.Menu{
		Items: []entity.MenuItem{
			{Name: "Butter Chicken", ShortName: "62"},
			{Name: "Chicken Tikka Masala", ShortName: "63"},
		},
	}

	_, err := s.repo.CreateMenu(nil, &menu)
	if err != nil {
		s.T().Fatal(err)
	}

	orderJson := fmt.Sprintf(`{"menu_uuid": "%s"}`, menu.Uuid.String())
	req := httptest.NewRequest(http.MethodPost, "/api/orders", strings.NewReader(orderJson))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := s.e.NewContext(req, rec)

	user, err := s.repo.CreateUser(nil, &entity.User{Name: "test"})
	if err != nil {
		s.T().Fatal(err)
	}
	c.Set("user", utils.BuildJwt(user))

	// Assertions
	if assert.NoError(s.T(), s.restBoundary.newOrder(c)) {
		assert.Equal(s.T(), http.StatusCreated, rec.Code)

		var createdOrder entity.Order
		json.NewDecoder(rec.Body).Decode(&createdOrder)

		assert.Equal(s.T(), menu.Uuid, createdOrder.MenuUuid)
		assert.Equal(s.T(), user.Uuid, createdOrder.Initiator)
	}
}

func (s *Suite) TestCreateOrderItem() {
	menu := entity.Menu{
		Items: []entity.MenuItem{
			{Name: "Butter Chicken", ShortName: "62"},
			{Name: "Chicken Tikka Masala", ShortName: "63"},
		},
	}

	_, err := s.repo.CreateMenu(nil, &menu)
	if err != nil {
		s.T().Fatal(err)
	}

	orderJson := fmt.Sprintf(`{"menu_uuid": "%s"}`, menu.Uuid.String())
	req := httptest.NewRequest(http.MethodPost, "/api/orders", strings.NewReader(orderJson))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := s.e.NewContext(req, rec)

	user, err := s.repo.CreateUser(nil, &entity.User{Name: ""})
	if err != nil {
		s.T().Fatal(err)
	}
	c.Set("user", utils.BuildJwt(user))

	// Assertions
	if assert.NoError(s.T(), s.restBoundary.newOrder(c)) {
		assert.Equal(s.T(), http.StatusCreated, rec.Code)

		var createdOrder entity.Order
		json.NewDecoder(rec.Body).Decode(&createdOrder)

		assert.Equal(s.T(), menu.Uuid, createdOrder.MenuUuid)
	}
}
