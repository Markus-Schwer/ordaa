package rest

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"gitlab.com/sfz.aalen/hackwerk/ordaa/boundary/utils"
	"gitlab.com/sfz.aalen/hackwerk/ordaa/entity"
	"gorm.io/gorm"
)

func (server *RestBoundary) newOrder(c echo.Context) error {
	var order entity.Order
	if err := utils.BindAndValidate(c, &order); err != nil {
		log.Ctx(server.ctx).Debug().Err(err).Msg("newOrder failed to bind request")
		return utils.WrapBindError(err)
	}

	return server.repo.Transaction(func(tx *gorm.DB) error {
		user, err := utils.CurrentUser(c, server.repo, tx)
		if err != nil {
			log.Ctx(server.ctx).Warn().Err(err).Msg("newOrder error getting current user")
			return utils.NewInternalServerError(err)
		}

		order.Initiator = user.Uuid

		createdOrder, err := server.repo.CreateOrder(tx, &order)
		if errors.Is(err, entity.ErrActiveOrderForMenuAlreadyExists) {
			log.Ctx(server.ctx).Warn().Err(err).Msg("newOrder error creating order")
			return utils.NewError(http.StatusBadRequest, err.Error())
		} else if err != nil {
			log.Ctx(server.ctx).Warn().Err(err).Msg("newOrder error creating order")
			return utils.NewInternalServerError(err)
		}

		return c.JSON(http.StatusCreated, createdOrder)
	})
}

func (server *RestBoundary) allOrders(c echo.Context) error {
	return server.repo.Transaction(func(tx *gorm.DB) error {
		orders, err := server.repo.GetAllOrders(tx)
		if err != nil {
			log.Ctx(server.ctx).Warn().Err(err).Msg("allOrders error getting orders")
			return utils.NewInternalServerError(err)
		}

		return c.JSON(http.StatusOK, orders)
	})
}

func (server *RestBoundary) getOrder(c echo.Context) error {
	uuid, err := utils.UuidParam(c, "uuid")
	if err != nil {
		log.Ctx(server.ctx).Warn().Err(err).Msg("getOrder error parsing uuid")
		return utils.NewStatusUnprocessableEntity(err.Error())
	}

	return server.repo.Transaction(func(tx *gorm.DB) error {
		orders, err := server.repo.GetOrder(tx, uuid)
		if err != nil {
			log.Ctx(server.ctx).Warn().Err(err).Msg("getOrder error getting order")
			return utils.NewInternalServerError(err)
		}

		return c.JSON(http.StatusOK, orders)
	})
}

func (server *RestBoundary) updateOrder(c echo.Context) error {
	uuid, err := utils.UuidParam(c, "uuid")
	if err != nil {
		log.Ctx(server.ctx).Warn().Err(err).Msg("updateOrder error parsing uuid")
		return utils.NewStatusUnprocessableEntity(err.Error())
	}

	var order entity.Order
	if err = utils.BindAndValidate(c, &order); err != nil {
		log.Ctx(server.ctx).Debug().Err(err).Msg("updateOrder error binding request")
		return utils.WrapBindError(err)
	}

	return server.repo.Transaction(func(tx *gorm.DB) error {
		_, err := server.repo.GetOrder(tx, uuid)
		if err != nil {
			log.Ctx(server.ctx).Warn().Err(err).Msg("updateOrder error getting order")
			return utils.NewInternalServerError(err)
		}

		user, err := utils.CurrentUser(c, server.repo, tx)
		if err != nil {
			log.Ctx(server.ctx).Warn().Err(err).Msg("updateOrder error getting current user")
			return utils.NewInternalServerError(err)
		}

		createdOrder, err := server.repo.UpdateOrder(tx, uuid, user.Uuid, &order)
		if err != nil {
			log.Ctx(server.ctx).Warn().Err(err).Msg("updateOrder error updating order")
			return utils.NewInternalServerError(err)
		}

		return c.JSON(http.StatusOK, createdOrder)
	})
}

func (server *RestBoundary) deleteOrder(c echo.Context) error {
	uuid, err := utils.UuidParam(c, "uuid")
	if err != nil {
		log.Ctx(server.ctx).Warn().Err(err).Msg("deleteOrder error parsing uuid")
		return utils.NewStatusUnprocessableEntity(err.Error())
	}

	return server.repo.Transaction(func(tx *gorm.DB) error {
		err = server.repo.DeleteOrder(tx, uuid)
		if err != nil {
			log.Ctx(server.ctx).Warn().Err(err).Msg("deleteOrder error deleting order")
			return utils.NewInternalServerError(err)
		}

		return c.NoContent(http.StatusOK)
	})
}

func (server *RestBoundary) newOrderItem(c echo.Context) error {
	orderUuid, err := utils.UuidParam(c, "order_uuid")
	if err != nil {
		log.Ctx(server.ctx).Warn().Err(err).Msg("newOrderItem error parsing order uuid")
		return utils.NewStatusUnprocessableEntity(err.Error())
	}

	var orderItem entity.OrderItem
	if err = utils.BindAndValidate(c, &orderItem); err != nil {
		log.Ctx(server.ctx).Debug().Err(err).Msg("newOrderItem error binding request")
		return utils.WrapBindError(err)
	}

	orderItem.OrderUuid = orderUuid

	return server.repo.Transaction(func(tx *gorm.DB) error {
		createdOrderItem, err := server.repo.CreateOrderItem(tx, orderUuid, &orderItem)
		if err != nil {
			log.Ctx(server.ctx).Warn().Err(err).Msg("newOrderItem error creating order item")
			return utils.NewInternalServerError(err)
		}

		return c.JSON(http.StatusCreated, createdOrderItem)
	})
}

func (server *RestBoundary) allOrderItems(c echo.Context) error {
	orderUuid, err := utils.UuidParam(c, "order_uuid")
	if err != nil {
		log.Ctx(server.ctx).Warn().Err(err).Msg("allOrderItems could not parse order uuid")
		return utils.NewStatusUnprocessableEntity(err.Error())
	}

	return server.repo.Transaction(func(tx *gorm.DB) error {
		orderItems, err := server.repo.GetAllOrderItems(tx, orderUuid)
		if err != nil {
			log.Ctx(server.ctx).Warn().Err(err).Msg("allOrderItems error getting order items")
			return utils.NewInternalServerError(err)
		}

		return c.JSON(http.StatusOK, orderItems)
	})
}

func (server *RestBoundary) getOrderItem(c echo.Context) error {
	orderUuid, err := utils.UuidParam(c, "order_uuid")
	if err != nil {
		log.Ctx(server.ctx).Warn().Err(err).Msg("getOrderItem error parsing order uuid")
		return utils.NewStatusUnprocessableEntity(err.Error())
	}

	orderItemUuid, err := utils.UuidParam(c, "uuid")
	if err != nil {
		log.Ctx(server.ctx).Warn().Err(err).Msg("getOrderItem error parsing order item uuid")
		return utils.NewStatusUnprocessableEntity(err.Error())
	}

	return server.repo.Transaction(func(tx *gorm.DB) error {
		orderItems, err := server.repo.GetOrderItem(tx, orderItemUuid)
		if err != nil {
			log.Ctx(server.ctx).Warn().Err(err).Msg("getOrderItem error getting order item")
			return utils.NewInternalServerError(err)
		}

		if orderItems.OrderUuid != orderUuid {
			err = errors.New("order item not found")
			log.Ctx(server.ctx).Warn().Err(err).Msg("getOrderItem error order item not found")
			return utils.NewInternalServerError(err)
		}

		return c.JSON(http.StatusOK, orderItems)
	})
}

func (server *RestBoundary) updateOrderItem(c echo.Context) error {
	orderUuid, err := utils.UuidParam(c, "order_uuid")
	if err != nil {
		log.Ctx(server.ctx).Warn().Err(err).Msg("updateOrderItem error parsing order uuid")
		return utils.NewStatusUnprocessableEntity(err.Error())
	}

	orderItemUuid, err := utils.UuidParam(c, "uuid")
	if err != nil {
		log.Ctx(server.ctx).Warn().Err(err).Msg("updateOrderItem error parsing order item uuid")
		return utils.NewStatusUnprocessableEntity(err.Error())
	}

	var orderItem entity.OrderItem
	if err = utils.BindAndValidate(c, &orderItem); err != nil {
		log.Ctx(server.ctx).Debug().Err(err).Msg("updateOrderItem error binding request")
		return utils.WrapBindError(err)
	}

	return server.repo.Transaction(func(tx *gorm.DB) error {
		user, err := utils.CurrentUser(c, server.repo, tx)
		if err != nil {
			log.Ctx(server.ctx).Warn().Err(err).Msg("updateOrderItem error getting current user")
			return utils.NewInternalServerError(err)
		}

		existingOrderItem, err := server.repo.GetOrderItem(tx, orderItemUuid)
		if err != nil {
			log.Ctx(server.ctx).Warn().Err(err).Msg("updateOrderItem error getting order item")
			return utils.NewInternalServerError(err)
		}

		if orderItem.OrderUuid == nil {
			orderItem.OrderUuid = orderUuid
		}

		if orderItem.MenuItemUuid == nil {
			orderItem.MenuItemUuid = existingOrderItem.MenuItemUuid
		}

		if orderItem.User == nil {
			orderItem.User = existingOrderItem.User
		}

		updatedOrderItem, err := server.repo.UpdateOrderItem(tx, orderItemUuid, user.Uuid, &orderItem)
		if err != nil {
			log.Ctx(server.ctx).Warn().Err(err).Msg("updateOrderItem error updating order item")
			return utils.NewInternalServerError(err)
		}

		if *updatedOrderItem.OrderUuid != *orderUuid {
			err = errors.New("order item not found")
			log.Ctx(server.ctx).Warn().Err(err).Msg("updateOrderItem error order item not found. The order item uuid is correct, but the order uuid is not")
			return utils.NewInternalServerError(err)
		}

		return c.JSON(http.StatusOK, updatedOrderItem)
	})
}

func (server *RestBoundary) deleteOrderItem(c echo.Context) error {
	orderItemUuid, err := utils.UuidParam(c, "order_uuid")
	if err != nil {
		log.Ctx(server.ctx).Warn().Err(err).Msg("deleteOrderItem error parsing order item uuid")
		return utils.NewStatusUnprocessableEntity(err.Error())
	}

	return server.repo.Transaction(func(tx *gorm.DB) error {
		err = server.repo.DeleteOrderItem(tx, orderItemUuid)
		if err != nil {
			log.Ctx(server.ctx).Warn().Err(err).Msg("deleteOrderItem error deleting order item")
			return utils.NewInternalServerError(err)
		}

		return c.NoContent(http.StatusOK)
	})
}
