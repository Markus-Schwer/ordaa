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
	"gitlab.com/sfz.aalen/hackwerk/dotinder/entity"
)

func (s *Suite) TestCreateOrder() {
	menuUuid, err := uuid.NewV4()
	if err != nil {
		s.T().Fatal(err)
	}

	menu := entity.Menu{
		Uuid: &menuUuid,
		Items: []entity.MenuItem{},
	}

	_, err = s.repo.CreateMenu(nil, &menu)
	if err != nil {
		s.T().Fatal(err)
	}

	orderJson := fmt.Sprintf(`{"menu_uuid": "%s"}`, menuUuid.String())
	req := httptest.NewRequest(http.MethodPost, "/api/orders", strings.NewReader(orderJson))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := s.e.NewContext(req, rec)

	// Assertions
	if assert.NoError(s.T(), s.restBoundary.newOrder(c)) {
		assert.Equal(s.T(), http.StatusCreated, rec.Code)

		var createdOrder entity.Order
		json.NewDecoder(rec.Body).Decode(&createdOrder)

		assert.Equal(s.T(), menuUuid, *createdOrder.MenuUuid)
	}
}
