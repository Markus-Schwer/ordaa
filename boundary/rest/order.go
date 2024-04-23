package rest

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/gofrs/uuid"
	"github.com/rs/zerolog/log"
	"gitlab.com/sfz.aalen/hackwerk/dotinder/entity"

	_ "github.com/lib/pq"
)

func (server *RestBoundary) newOrder(c echo.Context) error {
	var order entity.NewOrder
	err := c.Bind(&order)
	if err != nil {
		log.Ctx(server.ctx).Error().Err(err).Msg(err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	tx := server.repo.Pool.MustBegin()
	createdOrder, err := server.repo.CreateOrder(tx, &order)
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

	c.JSON(http.StatusOK, createdOrder)
	return nil
}

func (server *RestBoundary) allOrders(c echo.Context) error {
	tx := server.repo.Pool.MustBegin()
	orders, err := server.repo.GetAllOrders(tx)
	if err != nil {
		log.Ctx(server.ctx).Error().Err(err).Msg(err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	c.JSON(http.StatusOK, orders)
	return nil
}

func (server *RestBoundary) getOrder(c echo.Context) error {
	uuidString := c.Param("uuid")
	uuid, err := uuid.FromString(uuidString)
	if err != nil {
		log.Ctx(server.ctx).Error().Err(err).Msg("error parsing uuid")
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	tx := server.repo.Pool.MustBegin()
	orders, err := server.repo.GetOrder(tx, uuid)
	if err != nil {
		log.Ctx(server.ctx).Error().Err(err).Msg("error getting order")
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	c.JSON(http.StatusOK, orders)
	return nil
}

func (server *RestBoundary) updateOrder(c echo.Context) error {
	uuidString := c.Param("uuid")
	uuid, err := uuid.FromString(uuidString)
	if err != nil {
		log.Ctx(server.ctx).Error().Err(err).Msg(err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	var order entity.NewOrder
	err = c.Bind(&order)
	if err != nil {
		log.Ctx(server.ctx).Error().Err(err).Msg(err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	tx := server.repo.Pool.MustBegin()
	createdOrder, err := server.repo.UpdateOrder(tx, uuid, &order)
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

	c.JSON(http.StatusOK, createdOrder)
	return nil
}

func (server *RestBoundary) deleteOrder(c echo.Context) error {
	uuidString := c.Param("uuid")
	uuid, err := uuid.FromString(uuidString)
	if err != nil {
		log.Ctx(server.ctx).Error().Err(err).Msg(err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	tx := server.repo.Pool.MustBegin()
	err = server.repo.DeleteOrder(tx, uuid)
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

	c.NoContent(http.StatusOK)
	return nil
}


func (server *RestBoundary) newOrderItem(c echo.Context) error {
	uuidString := c.Param("order_uuid")
	orderUuid, err := uuid.FromString(uuidString)
	if err != nil {
		log.Ctx(server.ctx).Error().Err(err).Msg(err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	var orderItem entity.NewOrderItem
	err = c.Bind(&orderItem)
	if err != nil {
		log.Ctx(server.ctx).Error().Err(err).Msg(err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	tx := server.repo.Pool.MustBegin()
	createdOrderItem, err := server.repo.CreateOrderItem(tx, orderUuid, orderItem)
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

	c.JSON(http.StatusOK, createdOrderItem)
	return nil
}

func (server *RestBoundary) allOrderItems(c echo.Context) error {
	uuidString := c.Param("order_uuid")
	orderUuid, err := uuid.FromString(uuidString)
	if err != nil {
		log.Ctx(server.ctx).Error().Err(err).Msg("could not parse order uuid")
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	tx := server.repo.Pool.MustBegin()
	orderItems, err := server.repo.GetAllOrderItems(tx, orderUuid)
	if err != nil {
		log.Ctx(server.ctx).Error().Err(err).Msg("could not get order items")
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	c.JSON(http.StatusOK, orderItems)
	return nil
}

func (server *RestBoundary) getOrderItem(c echo.Context) error {
	uuidString := c.Param("order_uuid")
	orderUuid, err := uuid.FromString(uuidString)
	if err != nil {
		log.Ctx(server.ctx).Error().Err(err).Msg(err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	uuidString = c.Param("uuid")
	orderItemUuid, err := uuid.FromString(uuidString)
	if err != nil {
		log.Ctx(server.ctx).Error().Err(err).Msg(err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	tx := server.repo.Pool.MustBegin()
	orderItems, err := server.repo.GetOrderItem(tx, orderItemUuid)
	if err != nil {
		log.Ctx(server.ctx).Error().Err(err).Msg(err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	if orderItems.OrderUuid != orderUuid {
		err = errors.New("order item not found")
		log.Ctx(server.ctx).Error().Err(err).Msg(err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	c.JSON(http.StatusOK, orderItems)
	return nil
}

func (server *RestBoundary) updateOrderItem(c echo.Context) error {
	uuidString := c.Param("order_uuid")
	orderUuid, err := uuid.FromString(uuidString)
	if err != nil {
		log.Ctx(server.ctx).Error().Err(err).Msg(err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	uuidString = c.Param("uuid")
	orderItemUuid, err := uuid.FromString(uuidString)
	if err != nil {
		log.Ctx(server.ctx).Error().Err(err).Msg(err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	var orderItem entity.NewOrderItem
	err = c.Bind(&orderItem)
	if err != nil {
		log.Ctx(server.ctx).Error().Err(err).Msg(err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	tx := server.repo.Pool.MustBegin()
	createdOrderItem, err := server.repo.UpdateOrderItem(tx, orderItemUuid, &orderItem)
	if err != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			log.Ctx(server.ctx).Error().Err(err).Msg(err.Error())
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
		log.Ctx(server.ctx).Error().Err(err).Msg(err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	if createdOrderItem.OrderUuid != orderUuid {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			log.Ctx(server.ctx).Error().Err(err).Msg(err.Error())
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
		err = errors.New("order item not found")
		log.Ctx(server.ctx).Error().Err(err).Msg(err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	err = tx.Commit()
	if err != nil {
		log.Ctx(server.ctx).Error().Err(err).Msg(err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	c.JSON(http.StatusOK, createdOrderItem)
	return nil
}

func (server *RestBoundary) deleteOrderItem(c echo.Context) error {
	uuidString := c.Param("order_uuid")
	orderItemUuid, err := uuid.FromString(uuidString)
	if err != nil {
		log.Ctx(server.ctx).Error().Err(err).Msg(err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	tx := server.repo.Pool.MustBegin()
	err = server.repo.DeleteOrderItem(tx, orderItemUuid)
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

	c.NoContent(http.StatusOK)
	return nil
}

