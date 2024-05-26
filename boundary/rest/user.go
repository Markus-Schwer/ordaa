package rest

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"gitlab.com/sfz.aalen/hackwerk/dotinder/crypto"
	"gitlab.com/sfz.aalen/hackwerk/dotinder/entity"
	"gitlab.com/sfz.aalen/hackwerk/dotinder/boundary/utils"
	"gorm.io/gorm"
)

func (server *RestBoundary) registerUser(c echo.Context) error {
	var user entity.PasswordUser
	var err error

	// Bind request
	if err = c.Bind(&user); err != nil {
		log.Ctx(server.ctx).Error().Err(err).Msg("registerUser failed to bind request")
		return utils.WrapBindError(err)
	}

	user.Password, err = crypto.GeneratePasswordHash(user.Password)
	if err != nil {
		log.Ctx(server.ctx).Error().Err(err).Msg("registerUser error generating password hash")
		return utils.NewInternalServerError(err)
	}

	return server.repo.Transaction(func(tx *gorm.DB) error {
		createdUser, err := server.repo.CreateUser(tx, &entity.User{Name: user.Username})
		if err != nil {
			log.Ctx(server.ctx).Error().Err(err).Msg("registerUser error creating user")
			return utils.NewInternalServerError(err)
		}
		user.UserUuid = createdUser.Uuid

		_, err = server.repo.CreatePasswordUser(tx, &user)
		if err != nil {
			log.Ctx(server.ctx).Error().Err(err).Msg("registerUser error creating password user")
			return utils.NewInternalServerError(err)
		}

		return c.JSON(http.StatusCreated, createdUser)
	})
}

func (server *RestBoundary) allUsers(c echo.Context) error {
	return server.repo.Transaction(func(tx *gorm.DB) error {
		users, err := server.repo.GetAllUsers(tx)
		if err != nil {
			log.Ctx(server.ctx).Error().Err(err).Msg("allUsers error getting users")
			return utils.NewInternalServerError(err)
		}

		return c.JSON(http.StatusOK, users)
	})
}

func (server *RestBoundary) getUser(c echo.Context) error {
	uuid, err := utils.UuidParam(c, "uuid")
	if err != nil {
		log.Ctx(server.ctx).Error().Err(err).Msg("getUser error parsing uuid")
		return utils.NewStatusUnprocessableEntity(err.Error())
	}

	return server.repo.Transaction(func(tx *gorm.DB) error {
		users, err := server.repo.GetUser(tx, uuid)
		if err != nil {
			log.Ctx(server.ctx).Error().Err(err).Msg("getUser error getting user")
			return utils.NewInternalServerError(err)
		}

		return c.JSON(http.StatusOK, users)
	})
}

func (server *RestBoundary) updateUser(c echo.Context) error {
	uuid, err := utils.UuidParam(c, "uuid")
	if err != nil {
		log.Ctx(server.ctx).Error().Err(err).Msg("updateUser error parsing uuid")
		return utils.NewStatusUnprocessableEntity(err.Error())
	}

	var user entity.User
	if err = c.Bind(&user); err != nil {
		log.Ctx(server.ctx).Error().Err(err).Msg("updateUser error binding request")
		return utils.WrapBindError(err)
	}

	return server.repo.Transaction(func(tx *gorm.DB) error {
		createdUser, err := server.repo.UpdateUser(tx, uuid, &user)
		if err != nil {
			log.Ctx(server.ctx).Error().Err(err).Msg("updateUser error updating user")
			return utils.NewInternalServerError(err)
		}

		return c.JSON(http.StatusOK, createdUser)
	})
}

func (server *RestBoundary) deleteUser(c echo.Context) error {
	uuid, err := utils.UuidParam(c, "uuid")
	if err != nil {
		log.Ctx(server.ctx).Error().Err(err).Msg("deleteUser error parsing uuid")
		return utils.NewStatusUnprocessableEntity(err.Error())
	}

	return server.repo.Transaction(func(tx *gorm.DB) error {
		err = server.repo.DeleteUser(tx, uuid)
		if err != nil {
			log.Ctx(server.ctx).Error().Err(err).Msg("deleteUser error deleting user")
			return utils.NewInternalServerError(err)
		}

		return c.NoContent(http.StatusNoContent)
	})
}
