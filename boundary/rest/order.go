package rest

import (
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog/log"

	_ "github.com/lib/pq"
)

//func (server *RestBoundary) newOrder(w http.ResponseWriter, r *http.Request) {
//	var menu entity.NewOrder
//	err := json.NewDecoder(r.Body).Decode(&menu)
//	if err != nil {
//		log.Ctx(server.ctx).Error().Err(err).Msg(err.Error())
//		http.Error(w, err.Error(), http.StatusBadRequest)
//		return
//	}
//
//	tx := server.repo.Pool.MustBegin()
//	createdOrder, err := server.repo.CreateOrder(tx, &menu)
//	if err != nil {
//		rollbackErr := tx.Rollback()
//		if rollbackErr != nil {
//			log.Ctx(server.ctx).Error().Err(err).Msg(err.Error())
//			http.Error(w, err.Error(), http.StatusInternalServerError)
//			return
//		}
//		log.Ctx(server.ctx).Error().Err(err).Msg(err.Error())
//		http.Error(w, err.Error(), http.StatusInternalServerError)
//		return
//	}
//	err = tx.Commit()
//	if err != nil {
//		log.Ctx(server.ctx).Error().Err(err).Msg(err.Error())
//		http.Error(w, err.Error(), http.StatusInternalServerError)
//		return
//	}
//
//	w.Header().Add("Content-Type", "application/json")
//	err = json.NewEncoder(w).Encode(createdOrder)
//	if err != nil {
//		log.Ctx(server.ctx).Error().Err(err).Msg(err.Error())
//		http.Error(w, err.Error(), http.StatusInternalServerError)
//		return
//	}
//}

func (server *RestBoundary) allOrders(w http.ResponseWriter, r *http.Request) {
	tx := server.repo.Pool.MustBegin()
	menus, err := server.repo.GetAllOrders(tx)
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

//func (server *RestBoundary) getOrder(w http.ResponseWriter, r *http.Request) {
//	uuidString := mux.Vars(r)["uuid"]
//	uuid, err := uuid.FromString(uuidString)
//	if err != nil {
//		log.Ctx(server.ctx).Error().Err(err).Msg(err.Error())
//		http.Error(w, err.Error(), http.StatusBadRequest)
//		return
//	}
//
//	tx := server.repo.Pool.MustBegin()
//	menus, err := server.repo.GetOrder(tx, uuid)
//	if err != nil {
//		log.Ctx(server.ctx).Error().Err(err).Msg(err.Error())
//		http.Error(w, err.Error(), http.StatusInternalServerError)
//		return
//	}
//
//	w.Header().Add("Content-Type", "application/json")
//	err = json.NewEncoder(w).Encode(menus)
//	if err != nil {
//		log.Ctx(server.ctx).Error().Err(err).Msg(err.Error())
//		http.Error(w, err.Error(), http.StatusInternalServerError)
//		return
//	}
//}
//
//func (server *RestBoundary) updateOrder(w http.ResponseWriter, r *http.Request) {
//	uuidString := mux.Vars(r)["uuid"]
//	uuid, err := uuid.FromString(uuidString)
//	if err != nil {
//		log.Ctx(server.ctx).Error().Err(err).Msg(err.Error())
//		http.Error(w, err.Error(), http.StatusBadRequest)
//		return
//	}
//
//	var menu entity.NewOrder
//	err = json.NewDecoder(r.Body).Decode(&menu)
//	if err != nil {
//		log.Ctx(server.ctx).Error().Err(err).Msg(err.Error())
//		http.Error(w, err.Error(), http.StatusBadRequest)
//		return
//	}
//
//	tx := server.repo.Pool.MustBegin()
//	createdOrder, err := server.repo.UpdateOrder(tx, uuid, &menu)
//	if err != nil {
//		rollbackErr := tx.Rollback()
//		if rollbackErr != nil {
//			log.Ctx(server.ctx).Error().Err(err).Msg(err.Error())
//			http.Error(w, err.Error(), http.StatusInternalServerError)
//			return
//		}
//		log.Ctx(server.ctx).Error().Err(err).Msg(err.Error())
//		http.Error(w, err.Error(), http.StatusInternalServerError)
//		return
//	}
//	err = tx.Commit()
//	if err != nil {
//		log.Ctx(server.ctx).Error().Err(err).Msg(err.Error())
//		http.Error(w, err.Error(), http.StatusInternalServerError)
//		return
//	}
//
//	w.Header().Add("Content-Type", "application/json")
//	err = json.NewEncoder(w).Encode(createdOrder)
//	if err != nil {
//		log.Ctx(server.ctx).Error().Err(err).Msg(err.Error())
//		http.Error(w, err.Error(), http.StatusInternalServerError)
//		return
//	}
//}
//
//func (server *RestBoundary) deleteOrder(w http.ResponseWriter, r *http.Request) {
//	uuidString := mux.Vars(r)["uuid"]
//	uuid, err := uuid.FromString(uuidString)
//	if err != nil {
//		log.Ctx(server.ctx).Error().Err(err).Msg(err.Error())
//		http.Error(w, err.Error(), http.StatusBadRequest)
//		return
//	}
//
//	tx := server.repo.Pool.MustBegin()
//	err = server.repo.DeleteOrder(tx, uuid)
//	if err != nil {
//		rollbackErr := tx.Rollback()
//		if rollbackErr != nil {
//			log.Ctx(server.ctx).Error().Err(err).Msg(err.Error())
//			http.Error(w, err.Error(), http.StatusInternalServerError)
//		}
//		log.Ctx(server.ctx).Error().Err(err).Msg(err.Error())
//		http.Error(w, err.Error(), http.StatusInternalServerError)
//		return
//	}
//	err = tx.Commit()
//	if err != nil {
//		log.Ctx(server.ctx).Error().Err(err).Msg(err.Error())
//		http.Error(w, err.Error(), http.StatusInternalServerError)
//		return
//	}
//
//	w.WriteHeader(http.StatusOK)
//}

