package rest

import (
	"encoding/json"
	"net/http"

	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	"gitlab.com/sfz.aalen/hackwerk/dotinder/entity"
	"gitlab.com/sfz.aalen/hackwerk/dotinder/crypto"

	_ "github.com/lib/pq"
)

func (server *RestBoundary) registerUser(w http.ResponseWriter, r *http.Request) {
	var user entity.NewPasswordUser
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Ctx(server.ctx).Error().Err(err).Msg(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Ctx(server.ctx).Info().Msgf("registering user %s", user.Username)

	user.Password, err = crypto.GeneratePasswordHash(user.Password)
	if err != nil {
		log.Ctx(server.ctx).Error().Err(err).Msg("error generating password hash")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tx := server.repo.Pool.MustBegin()
	createdUser, err := server.repo.CreateUser(tx, &entity.NewUser{Name: user.Username})
	if err != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			log.Ctx(server.ctx).Error().Err(err).Msg("error rolling back transaction")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		log.Ctx(server.ctx).Error().Err(err).Msg("error creating user")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	user.UserUuid = createdUser.Uuid

	createdPwUser, err := server.repo.CreatePasswordUser(tx, &user)
	if err != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			log.Ctx(server.ctx).Error().Err(err).Msg("error rolling back transaction")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		log.Ctx(server.ctx).Error().Err(err).Msg("error creating password user")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = tx.Commit()
	if err != nil {
		log.Ctx(server.ctx).Error().Err(err).Msg("error committing transaction")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(createdPwUser)
	if err != nil {
		log.Ctx(server.ctx).Error().Err(err).Msg(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (server *RestBoundary) allUsers(w http.ResponseWriter, r *http.Request) {
	tx := server.repo.Pool.MustBegin()
	users, err := server.repo.GetAllUsers(tx)
	if err != nil {
		log.Ctx(server.ctx).Error().Err(err).Msg(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(users)
	if err != nil {
		log.Ctx(server.ctx).Error().Err(err).Msg(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (server *RestBoundary) getUser(w http.ResponseWriter, r *http.Request) {
	uuidString := mux.Vars(r)["uuid"]
	uuid, err := uuid.FromString(uuidString)
	if err != nil {
		log.Ctx(server.ctx).Error().Err(err).Msg(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tx := server.repo.Pool.MustBegin()
	users, err := server.repo.GetUser(tx, uuid)
	if err != nil {
		log.Ctx(server.ctx).Error().Err(err).Msg(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(users)
	if err != nil {
		log.Ctx(server.ctx).Error().Err(err).Msg(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (server *RestBoundary) updateUser(w http.ResponseWriter, r *http.Request) {
	uuidString := mux.Vars(r)["uuid"]
	uuid, err := uuid.FromString(uuidString)
	if err != nil {
		log.Ctx(server.ctx).Error().Err(err).Msg(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var user entity.NewUser
	err = json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Ctx(server.ctx).Error().Err(err).Msg(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tx := server.repo.Pool.MustBegin()
	createdUser, err := server.repo.UpdateUser(tx, uuid, &user)
	if err != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			log.Ctx(server.ctx).Error().Err(err).Msg(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		log.Ctx(server.ctx).Error().Err(err).Msg(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = tx.Commit()
	if err != nil {
		log.Ctx(server.ctx).Error().Err(err).Msg(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(createdUser)
	if err != nil {
		log.Ctx(server.ctx).Error().Err(err).Msg(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (server *RestBoundary) deleteUser(w http.ResponseWriter, r *http.Request) {
	uuidString := mux.Vars(r)["uuid"]
	uuid, err := uuid.FromString(uuidString)
	if err != nil {
		log.Ctx(server.ctx).Error().Err(err).Msg(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tx := server.repo.Pool.MustBegin()
	err = server.repo.DeleteUser(tx, uuid)
	if err != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			log.Ctx(server.ctx).Error().Err(err).Msg(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		log.Ctx(server.ctx).Error().Err(err).Msg(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = tx.Commit()
	if err != nil {
		log.Ctx(server.ctx).Error().Err(err).Msg(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
