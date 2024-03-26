package rest

import (
	"encoding/json"
	"net/http"

	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	"gitlab.com/sfz.aalen/hackwerk/dotinder/entity"

	_ "github.com/lib/pq"
)

func (server *RestBoundary) newMenu(w http.ResponseWriter, r *http.Request) {
	var menu entity.NewMenu
	err := json.NewDecoder(r.Body).Decode(&menu)
	if err != nil {
		log.Ctx(server.ctx).Error().Err(err).Msg(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tx := server.repo.Pool.MustBegin()
	createdMenu, err := server.repo.CreateMenu(tx, &menu)
	if err != nil {
		rollback_err := tx.Rollback()
		if rollback_err != nil {
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
	err = json.NewEncoder(w).Encode(createdMenu)
	if err != nil {
		log.Ctx(server.ctx).Error().Err(err).Msg(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (server *RestBoundary) allMenus(w http.ResponseWriter, r *http.Request) {
	tx := server.repo.Pool.MustBegin()
	menus, err := server.repo.GetAllMenus(tx)
	if err != nil {
		log.Ctx(server.ctx).Error().Err(err).Msg(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(menus)
	if err != nil {
		log.Ctx(server.ctx).Error().Err(err).Msg(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (server *RestBoundary) getMenu(w http.ResponseWriter, r *http.Request) {
	uuid_string := mux.Vars(r)["uuid"]
	uuid, err := uuid.FromString(uuid_string)
	if err != nil {
		log.Ctx(server.ctx).Error().Err(err).Msg(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tx := server.repo.Pool.MustBegin()
	menus, err := server.repo.GetMenu(tx, uuid)
	if err != nil {
		log.Ctx(server.ctx).Error().Err(err).Msg(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(menus)
	if err != nil {
		log.Ctx(server.ctx).Error().Err(err).Msg(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (server *RestBoundary) updateMenu(w http.ResponseWriter, r *http.Request) {
	uuid_string := mux.Vars(r)["uuid"]
	uuid, err := uuid.FromString(uuid_string)
	if err != nil {
		log.Ctx(server.ctx).Error().Err(err).Msg(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var menu entity.NewMenu
	err = json.NewDecoder(r.Body).Decode(&menu)
	if err != nil {
		log.Ctx(server.ctx).Error().Err(err).Msg(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tx := server.repo.Pool.MustBegin()
	createdMenu, err := server.repo.UpdateMenu(tx, uuid, &menu)
	if err != nil {
		rollback_err := tx.Rollback()
		if rollback_err != nil {
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
	err = json.NewEncoder(w).Encode(createdMenu)
	if err != nil {
		log.Ctx(server.ctx).Error().Err(err).Msg(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (server *RestBoundary) deleteMenu(w http.ResponseWriter, r *http.Request) {
	uuid_string := mux.Vars(r)["uuid"]
	uuid, err := uuid.FromString(uuid_string)
	if err != nil {
		log.Ctx(server.ctx).Error().Err(err).Msg(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tx := server.repo.Pool.MustBegin()
	err = server.repo.DeleteMenu(tx, uuid)
	if err != nil {
		rollback_err := tx.Rollback()
		if rollback_err != nil {
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

