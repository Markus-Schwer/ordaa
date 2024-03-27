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

func (server *RestBoundary) newOrder(w http.ResponseWriter, r *http.Request) {
	var order entity.NewOrder
	err := json.NewDecoder(r.Body).Decode(&order)
	if err != nil {
		log.Ctx(server.ctx).Error().Err(err).Msg(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tx := server.repo.Pool.MustBegin()
	createdOrder, err := server.repo.CreateOrder(tx, &order)
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
	err = json.NewEncoder(w).Encode(createdOrder)
	if err != nil {
		log.Ctx(server.ctx).Error().Err(err).Msg(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (server *RestBoundary) allOrders(w http.ResponseWriter, r *http.Request) {
	tx := server.repo.Pool.MustBegin()
	orders, err := server.repo.GetAllOrders(tx)
	if err != nil {
		log.Ctx(server.ctx).Error().Err(err).Msg(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(orders)
	if err != nil {
		log.Ctx(server.ctx).Error().Err(err).Msg(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (server *RestBoundary) getOrder(w http.ResponseWriter, r *http.Request) {
	uuidString := mux.Vars(r)["uuid"]
	uuid, err := uuid.FromString(uuidString)
	if err != nil {
		log.Ctx(server.ctx).Error().Err(err).Msg(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tx := server.repo.Pool.MustBegin()
	orders, err := server.repo.GetOrder(tx, uuid)
	if err != nil {
		log.Ctx(server.ctx).Error().Err(err).Msg(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(orders)
	if err != nil {
		log.Ctx(server.ctx).Error().Err(err).Msg(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (server *RestBoundary) updateOrder(w http.ResponseWriter, r *http.Request) {
	uuidString := mux.Vars(r)["uuid"]
	uuid, err := uuid.FromString(uuidString)
	if err != nil {
		log.Ctx(server.ctx).Error().Err(err).Msg(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var order entity.NewOrder
	err = json.NewDecoder(r.Body).Decode(&order)
	if err != nil {
		log.Ctx(server.ctx).Error().Err(err).Msg(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tx := server.repo.Pool.MustBegin()
	createdOrder, err := server.repo.UpdateOrder(tx, uuid, &order)
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
	err = json.NewEncoder(w).Encode(createdOrder)
	if err != nil {
		log.Ctx(server.ctx).Error().Err(err).Msg(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (server *RestBoundary) deleteOrder(w http.ResponseWriter, r *http.Request) {
	uuidString := mux.Vars(r)["uuid"]
	uuid, err := uuid.FromString(uuidString)
	if err != nil {
		log.Ctx(server.ctx).Error().Err(err).Msg(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tx := server.repo.Pool.MustBegin()
	err = server.repo.DeleteOrder(tx, uuid)
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

