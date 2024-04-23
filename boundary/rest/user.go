package rest

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/gofrs/uuid"
	"github.com/rs/zerolog/log"
	"gitlab.com/sfz.aalen/hackwerk/dotinder/crypto"
	"gitlab.com/sfz.aalen/hackwerk/dotinder/entity"

	_ "github.com/lib/pq"
)

func (server *RestBoundary) registerUser(c echo.Context) error {
	var user entity.NewPasswordUser
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

	tx := server.repo.Pool.MustBegin()
	createdUser, err := server.repo.CreateUser(tx, &entity.NewUser{Name: user.Username})
	if err != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			log.Ctx(server.ctx).Error().Err(err).Msg("error rolling back transaction")
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
		log.Ctx(server.ctx).Error().Err(err).Msg("error creating user")
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	user.UserUuid = createdUser.Uuid

	createdPwUser, err := server.repo.CreatePasswordUser(tx, &user)
	if err != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			log.Ctx(server.ctx).Error().Err(err).Msg("error rolling back transaction")
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
		log.Ctx(server.ctx).Error().Err(err).Msg("error creating password user")
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	err = tx.Commit()
	if err != nil {
		log.Ctx(server.ctx).Error().Err(err).Msg("error committing transaction")
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	c.JSON(http.StatusCreated, createdPwUser)
	return nil
}

func (server *RestBoundary) allUsers(c echo.Context) error {
	tx := server.repo.Pool.MustBegin()
	users, err := server.repo.GetAllUsers(tx)
	if err != nil {
		log.Ctx(server.ctx).Error().Err(err).Msg(err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	c.JSON(http.StatusOK, users)
	return nil
}

func (server *RestBoundary) getUser(c echo.Context) error {
	uuidString := c.Param("uuid")
	uuid, err := uuid.FromString(uuidString)
	if err != nil {
		log.Ctx(server.ctx).Error().Err(err).Msg(err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	tx := server.repo.Pool.MustBegin()
	users, err := server.repo.GetUser(tx, uuid)
	if err != nil {
		log.Ctx(server.ctx).Error().Err(err).Msg(err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	c.JSON(http.StatusOK, users)
	return nil
}

func (server *RestBoundary) updateUser(c echo.Context) error {
	uuidString := c.Param("uuid")
	uuid, err := uuid.FromString(uuidString)
	if err != nil {
		log.Ctx(server.ctx).Error().Err(err).Msg(err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	var user entity.NewUser
	err = c.Bind(&user)
	if err != nil {
		log.Ctx(server.ctx).Error().Err(err).Msg(err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	tx := server.repo.Pool.MustBegin()
	createdUser, err := server.repo.UpdateUser(tx, uuid, &user)
	if err != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
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

	c.JSON(http.StatusOK, createdUser)
	return nil
}

func (server *RestBoundary) deleteUser(c echo.Context) error {
	uuidString := c.Param("uuid")
	uuid, err := uuid.FromString(uuidString)
	if err != nil {
		log.Ctx(server.ctx).Error().Err(err).Msg(err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	tx := server.repo.Pool.MustBegin()
	err = server.repo.DeleteUser(tx, uuid)
	if err != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
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

	c.NoContent(http.StatusNoContent)
	return nil
}
