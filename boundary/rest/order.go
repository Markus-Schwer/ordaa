package rest

import (
	"encoding/json"
	"errors"
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


func (server *RestBoundary) newOrderItem(w http.ResponseWriter, r *http.Request) {
	uuidString := mux.Vars(r)["order_uuid"]
	orderUuid, err := uuid.FromString(uuidString)
	if err != nil {
		log.Ctx(server.ctx).Error().Err(err).Msg(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var orderItem entity.NewOrderItem
	err = json.NewDecoder(r.Body).Decode(&orderItem)
	if err != nil {
		log.Ctx(server.ctx).Error().Err(err).Msg(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tx := server.repo.Pool.MustBegin()
	createdOrderItem, err := server.repo.CreateOrderItem(tx, orderUuid, orderItem)
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
	err = json.NewEncoder(w).Encode(createdOrderItem)
	if err != nil {
		log.Ctx(server.ctx).Error().Err(err).Msg(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (server *RestBoundary) allOrderItems(w http.ResponseWriter, r *http.Request) {
	uuidString := mux.Vars(r)["order_uuid"]
	orderUuid, err := uuid.FromString(uuidString)
	if err != nil {
		log.Ctx(server.ctx).Error().Err(err).Msg(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tx := server.repo.Pool.MustBegin()
	orderItems, err := server.repo.GetAllOrderItems(tx, orderUuid)
	if err != nil {
		log.Ctx(server.ctx).Error().Err(err).Msg(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(orderItems)
	if err != nil {
		log.Ctx(server.ctx).Error().Err(err).Msg(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (server *RestBoundary) getOrderItem(w http.ResponseWriter, r *http.Request) {
	uuidString := mux.Vars(r)["order_uuid"]
	orderUuid, err := uuid.FromString(uuidString)
	if err != nil {
		log.Ctx(server.ctx).Error().Err(err).Msg(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	uuidString = mux.Vars(r)["uuid"]
	orderItemUuid, err := uuid.FromString(uuidString)
	if err != nil {
		log.Ctx(server.ctx).Error().Err(err).Msg(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tx := server.repo.Pool.MustBegin()
	orderItems, err := server.repo.GetOrderItem(tx, orderItemUuid)
	if err != nil {
		log.Ctx(server.ctx).Error().Err(err).Msg(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if orderItems.OrderUuid != orderUuid {
		err = errors.New("order item not found")
		log.Ctx(server.ctx).Error().Err(err).Msg(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(orderItems)
	if err != nil {
		log.Ctx(server.ctx).Error().Err(err).Msg(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (server *RestBoundary) updateOrderItem(w http.ResponseWriter, r *http.Request) {
	uuidString := mux.Vars(r)["order_uuid"]
	orderUuid, err := uuid.FromString(uuidString)
	if err != nil {
		log.Ctx(server.ctx).Error().Err(err).Msg(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	uuidString = mux.Vars(r)["uuid"]
	orderItemUuid, err := uuid.FromString(uuidString)
	if err != nil {
		log.Ctx(server.ctx).Error().Err(err).Msg(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var orderItem entity.NewOrderItem
	err = json.NewDecoder(r.Body).Decode(&orderItem)
	if err != nil {
		log.Ctx(server.ctx).Error().Err(err).Msg(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tx := server.repo.Pool.MustBegin()
	createdOrderItem, err := server.repo.UpdateOrderItem(tx, orderItemUuid, &orderItem)
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

	if createdOrderItem.OrderUuid != orderUuid {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			log.Ctx(server.ctx).Error().Err(err).Msg(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = errors.New("order item not found")
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
	err = json.NewEncoder(w).Encode(createdOrderItem)
	if err != nil {
		log.Ctx(server.ctx).Error().Err(err).Msg(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (server *RestBoundary) deleteOrderItem(w http.ResponseWriter, r *http.Request) {
	uuidString := mux.Vars(r)["uuid"]
	orderItemUuid, err := uuid.FromString(uuidString)
	if err != nil {
		log.Ctx(server.ctx).Error().Err(err).Msg(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tx := server.repo.Pool.MustBegin()
	err = server.repo.DeleteOrderItem(tx, orderItemUuid)
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

