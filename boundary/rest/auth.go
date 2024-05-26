package rest

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"gitlab.com/sfz.aalen/hackwerk/dotinder/boundary/auth"
	"gitlab.com/sfz.aalen/hackwerk/dotinder/boundary/utils"
	"gorm.io/gorm"
)

type LoginResponse struct {
	Jwt string `json:"jwt"`
}

func (server *RestBoundary) login(c echo.Context) error {
	var creds auth.Credentials
	if err := utils.BindAndValidate(c, &creds); err != nil {
		log.Ctx(server.ctx).Debug().Err(err).Msg("login failed to bind request")
		return utils.WrapBindError(err)
	}

	return server.repo.Transaction(func(tx *gorm.DB) error {
		token, err := server.authService.Signin(tx, &creds)
		if err != nil {
			log.Ctx(server.ctx).Error().Err(err).Msg("login error signing in")
			return utils.NewInternalServerError(err)
		}
		rawToken, err := auth.SignToken(token)
		if err != nil {
			log.Ctx(server.ctx).Error().Err(err).Msg("login error signing token")
			return utils.NewInternalServerError(err)
		}

		loginResponse := LoginResponse{Jwt: rawToken}
		return c.JSON(http.StatusOK, loginResponse)
	})
}
