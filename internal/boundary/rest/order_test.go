package rest

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/gofrs/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/Markus-Schwer/ordaa/internal/boundary/utils"
	"github.com/Markus-Schwer/ordaa/internal/entity"
	"gorm.io/gorm"
)

func (s *Suite) TestCreateOrder() {
	menuUuid := uuid.Must(uuid.NewV4())
	menu := &entity.Menu{
		Uuid: &menuUuid,
		Items: []entity.MenuItem{
			{Name: "Butter Chicken", ShortName: "62"},
			{Name: "Chicken Tikka Masala", ShortName: "63"},
		},
	}

	s.repo.GetMenuFunc = func(tx *gorm.DB, uuidMoqParam *uuid.UUID) (*entity.Menu, error) {
		if *uuidMoqParam == menuUuid {
			return menu, nil
		}
		return nil, entity.ErrMenuNotFound
	}

	s.repo.CreateOrderFunc = func(tx *gorm.DB, order *entity.Order) (*entity.Order, error) {
		orderUuid := uuid.Must(uuid.NewV4())
		order.Uuid = &orderUuid
		return order, nil
	}

	orderJson := fmt.Sprintf(`{"menu_uuid": "%s"}`, menuUuid.String())
	req := httptest.NewRequest(http.MethodPost, "/api/orders", strings.NewReader(orderJson))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := s.e.NewContext(req, rec)

	c.Set("user", utils.BuildJwt(s.user))

	// Assertions
	if assert.NoError(s.T(), s.restBoundary.newOrder(c)) {
		assert.Equal(s.T(), http.StatusCreated, rec.Code)

		var createdOrder entity.Order
		json.NewDecoder(rec.Body).Decode(&createdOrder)

		assert.Equal(s.T(), menu.Uuid, createdOrder.MenuUuid)
		assert.Equal(s.T(), s.user.Uuid, createdOrder.Initiator)
	}
}

func (s *Suite) TestCreateOrderItem() {
	menuUuid := uuid.Must(uuid.NewV4())
	menu := &entity.Menu{
		Uuid: &menuUuid,
		Items: []entity.MenuItem{
			{Name: "Butter Chicken", ShortName: "62"},
			{Name: "Chicken Tikka Masala", ShortName: "63"},
		},
	}

	s.repo.GetMenuFunc = func(tx *gorm.DB, uuidMoqParam *uuid.UUID) (*entity.Menu, error) {
		if *uuidMoqParam == menuUuid {
			return menu, nil
		}
		return nil, entity.ErrMenuNotFound
	}

	s.repo.CreateOrderFunc = func(tx *gorm.DB, order *entity.Order) (*entity.Order, error) {
		orderUuid := uuid.Must(uuid.NewV4())
		order.Uuid = &orderUuid
		return order, nil
	}

	orderJson := fmt.Sprintf(`{"menu_uuid": "%s"}`, menu.Uuid.String())
	req := httptest.NewRequest(http.MethodPost, "/api/orders", strings.NewReader(orderJson))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := s.e.NewContext(req, rec)

	c.Set("user", utils.BuildJwt(s.user))

	// Assertions
	if assert.NoError(s.T(), s.restBoundary.newOrder(c)) {
		assert.Equal(s.T(), http.StatusCreated, rec.Code)

		var createdOrder entity.Order
		json.NewDecoder(rec.Body).Decode(&createdOrder)

		assert.Equal(s.T(), menu.Uuid, createdOrder.MenuUuid)
	}
}
