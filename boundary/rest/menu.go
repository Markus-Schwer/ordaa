package rest

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/gofrs/uuid"
	"github.com/rs/zerolog/log"
	"gitlab.com/sfz.aalen/hackwerk/dotinder/entity"
	"gorm.io/gorm"
)

func (server *RestBoundary) newMenu(c echo.Context) error {
	var menu entity.Menu
	err := c.Bind(&menu)
	if err != nil {
		log.Ctx(server.ctx).Error().Err(err).Msg(err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	return server.repo.Transaction(func(tx *gorm.DB) error {
		createdMenu, err := server.repo.CreateMenu(tx, &menu)
		if err != nil {
			log.Ctx(server.ctx).Error().Err(err).Msg(err.Error())
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
	
		return c.JSON(http.StatusOK, createdMenu)
	})
}

func (server *RestBoundary) allMenus(c echo.Context) error {
	return server.repo.Transaction(func(tx *gorm.DB) error {
		menus, err := server.repo.GetAllMenus(tx)
		if err != nil {
			log.Ctx(server.ctx).Error().Err(err).Msg(err.Error())
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
	
		c.JSON(http.StatusOK, menus)
		return nil
	})
}

func (server *RestBoundary) getMenu(c echo.Context) error {
	uuid_string := c.Param("uuid")
	uuid, err := uuid.FromString(uuid_string)
	if err != nil {
		log.Ctx(server.ctx).Error().Err(err).Msg(err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	return server.repo.Transaction(func(tx *gorm.DB) error {
		menus, err := server.repo.GetMenu(tx, uuid)
		if err != nil {
			log.Ctx(server.ctx).Error().Err(err).Msg(err.Error())
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
	
		c.JSON(http.StatusOK, menus)
		return nil
	})
}

func (server *RestBoundary) updateMenu(c echo.Context) error {
	uuid_string := c.Param("uuid")
	uuid, err := uuid.FromString(uuid_string)
	if err != nil {
		log.Ctx(server.ctx).Error().Err(err).Msg(err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	var menu entity.Menu
	err = c.Bind(&menu)
	if err != nil {
		log.Ctx(server.ctx).Error().Err(err).Msg(err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	return server.repo.Transaction(func(tx *gorm.DB) error {
		createdMenu, err := server.repo.UpdateMenu(tx, uuid, &menu)
		if err != nil {
			tx.Rollback()
			log.Ctx(server.ctx).Error().Err(err).Msg(err.Error())
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
	
		c.JSON(http.StatusOK, createdMenu)
		return nil
	})
}

func (server *RestBoundary) deleteMenu(c echo.Context) error {
	uuid_string := c.Param("uuid")
	uuid, err := uuid.FromString(uuid_string)
	if err != nil {
		log.Ctx(server.ctx).Error().Err(err).Msg(err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	return server.repo.Transaction(func(tx *gorm.DB) error {
		err = server.repo.DeleteMenu(tx, uuid)
		if err != nil {
			tx.Rollback()
			log.Ctx(server.ctx).Error().Err(err).Msg(err.Error())
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
	
		c.NoContent(http.StatusOK)
		return nil
	})
}

