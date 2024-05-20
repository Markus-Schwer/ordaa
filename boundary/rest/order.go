package rest

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/gofrs/uuid"
	"github.com/rs/zerolog/log"
	"gitlab.com/sfz.aalen/hackwerk/dotinder/entity"
	"gorm.io/gorm"
)

func (server *RestBoundary) newOrder(c echo.Context) error {
	var order entity.Order
	err := c.Bind(&order)
	if err != nil {
		log.Ctx(server.ctx).Error().Err(err).Msg(err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	return server.repo.Transaction(func(tx *gorm.DB) error {
		createdOrder, err := server.repo.CreateOrder(tx, &order)
		if err != nil {
			log.Ctx(server.ctx).Error().Err(err).Msg(err.Error())
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
	
		return c.JSON(http.StatusOK, createdOrder)
	})
}

func (server *RestBoundary) allOrders(c echo.Context) error {
	return server.repo.Transaction(func(tx *gorm.DB) error {
		orders, err := server.repo.GetAllOrders(tx)
		if err != nil {
			log.Ctx(server.ctx).Error().Err(err).Msg(err.Error())
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
	
		return c.JSON(http.StatusOK, orders)
	})
}

func (server *RestBoundary) getOrder(c echo.Context) error {
	uuidString := c.Param("uuid")
	uuid, err := uuid.FromString(uuidString)
	if err != nil {
		log.Ctx(server.ctx).Error().Err(err).Msg("error parsing uuid")
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	return server.repo.Transaction(func(tx *gorm.DB) error {
		orders, err := server.repo.GetOrder(tx, uuid)
		if err != nil {
			log.Ctx(server.ctx).Error().Err(err).Msg("error getting order")
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
	
		return c.JSON(http.StatusOK, orders)
	})
}

func (server *RestBoundary) updateOrder(c echo.Context) error {
	uuidString := c.Param("uuid")
	uuid, err := uuid.FromString(uuidString)
	if err != nil {
		log.Ctx(server.ctx).Error().Err(err).Msg(err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	var order entity.Order
	err = c.Bind(&order)
	if err != nil {
		log.Ctx(server.ctx).Error().Err(err).Msg(err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	return server.repo.Transaction(func(tx *gorm.DB) error {
		createdOrder, err := server.repo.UpdateOrder(tx, uuid, &order)
		if err != nil {
			log.Ctx(server.ctx).Error().Err(err).Msg(err.Error())
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
	
		return c.JSON(http.StatusOK, createdOrder)
	})
}

func (server *RestBoundary) deleteOrder(c echo.Context) error {
	uuidString := c.Param("uuid")
	uuid, err := uuid.FromString(uuidString)
	if err != nil {
		log.Ctx(server.ctx).Error().Err(err).Msg(err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	return server.repo.Transaction(func(tx *gorm.DB) error {
		err = server.repo.DeleteOrder(tx, uuid)
		if err != nil {
			return err
		}
	
		return c.NoContent(http.StatusOK)
	})
}


func (server *RestBoundary) newOrderItem(c echo.Context) error {
	uuidString := c.Param("order_uuid")
	orderUuid, err := uuid.FromString(uuidString)
	if err != nil {
		log.Ctx(server.ctx).Error().Err(err).Msg(err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	var orderItem entity.OrderItem
	err = c.Bind(&orderItem)
	if err != nil {
		log.Ctx(server.ctx).Error().Err(err).Msg(err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	return server.repo.Transaction(func(tx *gorm.DB) error {
		createdOrderItem, err := server.repo.CreateOrderItem(tx, orderUuid, &orderItem)
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, createdOrderItem)
	})
}

func (server *RestBoundary) allOrderItems(c echo.Context) error {
	uuidString := c.Param("order_uuid")
	orderUuid, err := uuid.FromString(uuidString)
	if err != nil {
		log.Ctx(server.ctx).Error().Err(err).Msg("could not parse order uuid")
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	return server.repo.Transaction(func(tx *gorm.DB) error {
		orderItems, err := server.repo.GetAllOrderItems(tx, orderUuid)
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, orderItems)
	})
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

	return server.repo.Transaction(func(tx *gorm.DB) error {
		orderItems, err := server.repo.GetOrderItem(tx, orderItemUuid)
		if err != nil {
			return err
		}

		if orderItems.OrderUuid != orderUuid {
			err = errors.New("order item not found")
			log.Ctx(server.ctx).Error().Err(err).Msg(err.Error())
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
	
		return c.JSON(http.StatusOK, orderItems)
	})
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

	var orderItem entity.OrderItem
	err = c.Bind(&orderItem)
	if err != nil {
		log.Ctx(server.ctx).Error().Err(err).Msg(err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	return server.repo.Transaction(func(tx *gorm.DB) error {
		createdOrderItem, err := server.repo.UpdateOrderItem(tx, orderItemUuid, &orderItem)
		if err != nil {
			return err
		}
	
		if createdOrderItem.OrderUuid != orderUuid {
			err = errors.New("order item not found")
			log.Ctx(server.ctx).Error().Err(err).Msg(err.Error())
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		return c.JSON(http.StatusOK, createdOrderItem)
	})
}

func (server *RestBoundary) deleteOrderItem(c echo.Context) error {
	uuidString := c.Param("order_uuid")
	orderItemUuid, err := uuid.FromString(uuidString)
	if err != nil {
		log.Ctx(server.ctx).Error().Err(err).Msg(err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	return server.repo.Transaction(func(tx *gorm.DB) error {
		err = server.repo.DeleteOrderItem(tx, orderItemUuid)
		if err != nil {
			log.Ctx(server.ctx).Error().Err(err).Msg(err.Error())
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		return c.NoContent(http.StatusOK)
	})
}

