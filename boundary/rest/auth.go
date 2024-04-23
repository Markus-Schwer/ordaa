package rest

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"gitlab.com/sfz.aalen/hackwerk/dotinder/boundary/auth"
)

type LoginResponse struct {
	Jwt string `json:"jwt"`
}

func (server *RestBoundary) login(c echo.Context) error {
	var creds auth.Credentials
	err := c.Bind(&creds)
	if err != nil {
		log.Ctx(server.ctx).Error().Err(err).Msg(err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	tx := server.repo.Pool.MustBegin()
	token, err := server.authService.Signin(tx, &creds)
	if err != nil {
		log.Ctx(server.ctx).Error().Err(err).Msg(err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	rawToken, err := auth.SignToken(token)
	if err != nil {
		log.Ctx(server.ctx).Error().Err(err).Msg(err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	if err = tx.Rollback(); err!= nil {
		log.Ctx(server.ctx).Error().Err(err).Msg("error rolling back transaction")
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	loginResponse := LoginResponse{Jwt: rawToken}
	c.JSON(http.StatusOK, loginResponse)
	return nil
}
