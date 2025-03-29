package rest

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/Markus-Schwer/ordaa/internal/boundary/utils"
	"github.com/Markus-Schwer/ordaa/internal/entity"
	"gorm.io/gorm"
)

func (server *RestBoundary) newMenu(c echo.Context) error {
	var menu entity.Menu
	if err := utils.BindAndValidate(c, &menu); err != nil {
		log.Ctx(server.ctx).Debug().Err(err).Msg("newMenu failed to bind request")
		return utils.WrapBindError(err)
	}

	return server.repo.Transaction(func(tx *gorm.DB) error {
		createdMenu, err := server.repo.CreateMenu(tx, &menu)
		if err != nil {
			log.Ctx(server.ctx).Warn().Err(err).Msg("newMenu error creating menu")
			return utils.NewInternalServerError(err)
		}

		return c.JSON(http.StatusCreated, createdMenu)
	})
}

func (server *RestBoundary) allMenus(c echo.Context) error {
	return server.repo.Transaction(func(tx *gorm.DB) error {
		menus, err := server.repo.GetAllMenus(tx)
		if err != nil {
			log.Ctx(server.ctx).Warn().Err(err).Msg("allMenus error getting menus")
			return utils.NewInternalServerError(err)
		}

		return c.JSON(http.StatusOK, menus)
	})
}

func (server *RestBoundary) getMenu(c echo.Context) error {
	uuid, err := utils.UuidParam(c, "uuid")
	if err != nil {
		log.Ctx(server.ctx).Warn().Err(err).Msg("getMenu error parsing uuid")
		return utils.NewStatusUnprocessableEntity(err.Error())
	}

	return server.repo.Transaction(func(tx *gorm.DB) error {
		menus, err := server.repo.GetMenu(tx, uuid)
		if err != nil {
			log.Ctx(server.ctx).Warn().Err(err).Msg("getMenu error getting menu")
			return utils.NewInternalServerError(err)
		}

		return c.JSON(http.StatusOK, menus)
	})
}

func (server *RestBoundary) updateMenu(c echo.Context) error {
	uuid, err := utils.UuidParam(c, "uuid")
	if err != nil {
		log.Ctx(server.ctx).Warn().Err(err).Msg("updateMenu error parsing uuid")
		return utils.NewStatusUnprocessableEntity(err.Error())
	}

	var menu entity.Menu
	if err = utils.BindAndValidate(c, &menu); err != nil {
		log.Ctx(server.ctx).Debug().Err(err).Msg("updateMenu error binding request")
		return utils.WrapBindError(err)
	}

	return server.repo.Transaction(func(tx *gorm.DB) error {
		createdMenu, err := server.repo.UpdateMenu(tx, uuid, &menu)
		if err != nil {
			log.Ctx(server.ctx).Warn().Err(err).Msg("updateMenu error updating menu")
			return utils.NewInternalServerError(err)
		}

		return c.JSON(http.StatusOK, createdMenu)
	})
}

func (server *RestBoundary) deleteMenu(c echo.Context) error {
	uuid, err := utils.UuidParam(c, "uuid")
	if err != nil {
		log.Ctx(server.ctx).Warn().Err(err).Msg("deleteMenu error parsing uuid")
		return utils.NewStatusUnprocessableEntity(err.Error())
	}

	return server.repo.Transaction(func(tx *gorm.DB) error {
		err = server.repo.DeleteMenu(tx, uuid)
		if err != nil {
			log.Ctx(server.ctx).Warn().Err(err).Msg("deleteMenu error deleting menu")
			return utils.NewInternalServerError(err)
		}

		return c.NoContent(http.StatusOK)
	})
}
