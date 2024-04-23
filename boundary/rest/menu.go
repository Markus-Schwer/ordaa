package rest

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/gofrs/uuid"
	"github.com/rs/zerolog/log"
	"gitlab.com/sfz.aalen/hackwerk/dotinder/entity"

	_ "github.com/lib/pq"
)

func (server *RestBoundary) newMenu(c echo.Context) error {
	var menu entity.NewMenu
	err := c.Bind(&menu)
	if err != nil {
		log.Ctx(server.ctx).Error().Err(err).Msg(err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	tx := server.repo.Pool.MustBegin()
	createdMenu, err := server.repo.CreateMenu(tx, &menu)
	if err != nil {
		rollback_err := tx.Rollback()
		if rollback_err != nil {
			log.Ctx(server.ctx).Error().Err(err).Msg(err.Error())
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
		log.Ctx(server.ctx).Error().Err(err).Msg(err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	err = tx.Commit()
	if err != nil {
		log.Ctx(server.ctx).Error().Err(err).Msg(err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	c.JSON(http.StatusOK, createdMenu)
	return nil
}

func (server *RestBoundary) allMenus(c echo.Context) error {
	tx := server.repo.Pool.MustBegin()
	menus, err := server.repo.GetAllMenus(tx)
	if err != nil {
		log.Ctx(server.ctx).Error().Err(err).Msg(err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	c.JSON(http.StatusOK, menus)
	return nil
}

func (server *RestBoundary) getMenu(c echo.Context) error {
	uuid_string := c.Param("uuid")
	uuid, err := uuid.FromString(uuid_string)
	if err != nil {
		log.Ctx(server.ctx).Error().Err(err).Msg(err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	tx := server.repo.Pool.MustBegin()
	menus, err := server.repo.GetMenu(tx, uuid)
	if err != nil {
		log.Ctx(server.ctx).Error().Err(err).Msg(err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	c.JSON(http.StatusOK, menus)
	return nil
}

func (server *RestBoundary) updateMenu(c echo.Context) error {
	uuid_string := c.Param("uuid")
	uuid, err := uuid.FromString(uuid_string)
	if err != nil {
		log.Ctx(server.ctx).Error().Err(err).Msg(err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	var menu entity.NewMenu
	err = c.Bind(&menu)
	if err != nil {
		log.Ctx(server.ctx).Error().Err(err).Msg(err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	tx := server.repo.Pool.MustBegin()
	createdMenu, err := server.repo.UpdateMenu(tx, uuid, &menu)
	if err != nil {
		rollback_err := tx.Rollback()
		if rollback_err != nil {
			log.Ctx(server.ctx).Error().Err(err).Msg(err.Error())
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
		log.Ctx(server.ctx).Error().Err(err).Msg(err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	err = tx.Commit()
	if err != nil {
		log.Ctx(server.ctx).Error().Err(err).Msg(err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	c.JSON(http.StatusOK, createdMenu)
	return nil
}

func (server *RestBoundary) deleteMenu(c echo.Context) error {
	uuid_string := c.Param("uuid")
	uuid, err := uuid.FromString(uuid_string)
	if err != nil {
		log.Ctx(server.ctx).Error().Err(err).Msg(err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	tx := server.repo.Pool.MustBegin()
	err = server.repo.DeleteMenu(tx, uuid)
	if err != nil {
		rollback_err := tx.Rollback()
		if rollback_err != nil {
			log.Ctx(server.ctx).Error().Err(err).Msg(err.Error())
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
		log.Ctx(server.ctx).Error().Err(err).Msg(err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	err = tx.Commit()
	if err != nil {
		log.Ctx(server.ctx).Error().Err(err).Msg(err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	c.NoContent(http.StatusOK)
	return nil
}

