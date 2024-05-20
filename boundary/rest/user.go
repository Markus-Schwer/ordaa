package rest

import (
	"net/http"

	"github.com/gofrs/uuid"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"gitlab.com/sfz.aalen/hackwerk/dotinder/crypto"
	"gitlab.com/sfz.aalen/hackwerk/dotinder/entity"
	"gorm.io/gorm"
)

func (server *RestBoundary) registerUser(c echo.Context) error {
	var user entity.PasswordUser
	err := c.Bind(&user)
	if err != nil {
		log.Ctx(server.ctx).Error().Err(err).Msg(err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	log.Ctx(server.ctx).Info().Msgf("registering user %s", user.Username)

	user.Password, err = crypto.GeneratePasswordHash(user.Password)
	if err != nil {
		log.Ctx(server.ctx).Error().Err(err).Msg("error generating password hash")
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return server.repo.Transaction(func(tx *gorm.DB) error {
		createdUser, err := server.repo.CreateUser(tx, &entity.User{Name: user.Username})
		if err != nil {
			log.Ctx(server.ctx).Error().Err(err).Msg("error creating user")
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
		user.UserUuid = *createdUser.Uuid

		log.Ctx(server.ctx).Info().Msgf("created pw user %s, %s", user.Username, user.Password)
		createdPwUser, err := server.repo.CreatePasswordUser(tx, &user)
		if err != nil {
			log.Ctx(server.ctx).Error().Err(err).Msg("error creating password user")
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		return c.JSON(http.StatusCreated, createdPwUser)
	})
}

func (server *RestBoundary) allUsers(c echo.Context) error {
	return server.repo.Transaction(func(tx *gorm.DB) error {
		users, err := server.repo.GetAllUsers(tx)
		if err != nil {
			log.Ctx(server.ctx).Error().Err(err).Msg(err.Error())
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		return c.JSON(http.StatusOK, users)
	})
}

func (server *RestBoundary) getUser(c echo.Context) error {
	uuidString := c.Param("uuid")
	uuid, err := uuid.FromString(uuidString)
	if err != nil {
		log.Ctx(server.ctx).Error().Err(err).Msg(err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	return server.repo.Transaction(func(tx *gorm.DB) error {
		users, err := server.repo.GetUser(tx, uuid)
		if err != nil {
			log.Ctx(server.ctx).Error().Err(err).Msg(err.Error())
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		return c.JSON(http.StatusOK, users)
	})
}

func (server *RestBoundary) updateUser(c echo.Context) error {
	uuidString := c.Param("uuid")
	uuid, err := uuid.FromString(uuidString)
	if err != nil {
		log.Ctx(server.ctx).Error().Err(err).Msg(err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	var user entity.User
	err = c.Bind(&user)
	if err != nil {
		log.Ctx(server.ctx).Error().Err(err).Msg(err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	return server.repo.Transaction(func(tx *gorm.DB) error {
		createdUser, err := server.repo.UpdateUser(tx, uuid, &user)
		if err != nil {
			log.Ctx(server.ctx).Error().Err(err).Msg(err.Error())
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		return c.JSON(http.StatusOK, createdUser)
	})
}

func (server *RestBoundary) deleteUser(c echo.Context) error {
	uuidString := c.Param("uuid")
	uuid, err := uuid.FromString(uuidString)
	if err != nil {
		log.Ctx(server.ctx).Error().Err(err).Msg(err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	return server.repo.Transaction(func(tx *gorm.DB) error {
		err = server.repo.DeleteUser(tx, uuid)
		if err != nil {
			log.Ctx(server.ctx).Error().Err(err).Msg(err.Error())
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		return c.NoContent(http.StatusNoContent)
	})
}
