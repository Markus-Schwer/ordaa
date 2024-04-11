package rest

import (
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog/log"
	"gitlab.com/sfz.aalen/hackwerk/dotinder/boundary/auth"
)

type LoginResponse struct {
	Jwt string `json:"jwt"`
}

func (server *RestBoundary) login(w http.ResponseWriter, r *http.Request) {
	var creds auth.Credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		log.Ctx(server.ctx).Error().Err(err).Msg(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tx := server.repo.Pool.MustBegin()
	token, err := server.authService.Signin(tx, &creds)
	if err != nil {
		log.Ctx(server.ctx).Error().Err(err).Msg(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rawToken, err := auth.SignToken(token)
	if err != nil {
		log.Ctx(server.ctx).Error().Err(err).Msg(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err = tx.Rollback(); err!= nil {
		log.Ctx(server.ctx).Error().Err(err).Msg("error rolling back transaction")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	loginResponse := LoginResponse{Jwt: rawToken}
	w.Header().Add("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(loginResponse)
	if err != nil {
		log.Ctx(server.ctx).Error().Err(err).Msg(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
