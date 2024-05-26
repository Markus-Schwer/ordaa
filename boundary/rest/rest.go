package rest

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"gitlab.com/sfz.aalen/hackwerk/dotinder/boundary/auth"
	"gitlab.com/sfz.aalen/hackwerk/dotinder/boundary/utils"
	"gitlab.com/sfz.aalen/hackwerk/dotinder/entity"
)

type RestBoundary struct {
	ctx         context.Context
	repo        entity.Repository
	authService *auth.AuthService
}

func NewRestBoundary(ctx context.Context, repo entity.Repository, authSerivce *auth.AuthService) *RestBoundary {
	return &RestBoundary{ctx: ctx, repo: repo, authService: authSerivce}
}

func (server *RestBoundary) Start(router *echo.Echo) {
	router.Validator = utils.NewValidator()

	authRouter := router.Group("/api", auth.AuthMiddleware(server.authService, func(c echo.Context, err error) error {
		server.authService.Logout(c)
		return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
	}))

	authRouter.POST("/menus", server.newMenu)
	authRouter.GET("/menus", server.allMenus)
	authRouter.GET("/menus/:uuid", server.getMenu)
	authRouter.PUT("/menus/:uuid", server.updateMenu)
	authRouter.DELETE("/menus/:uuid", server.deleteMenu)

	authRouter.POST("/orders", server.newOrder)
	authRouter.GET("/orders", server.allOrders)
	authRouter.GET("/orders/:uuid", server.getOrder)
	authRouter.PUT("/orders/:uuid", server.updateOrder)
	authRouter.DELETE("/orders/:uuid", server.deleteOrder)

	authRouter.POST("/orders/:order_uuid/items", server.newOrderItem)
	authRouter.GET("/orders/:order_uuid/items", server.allOrderItems)
	authRouter.GET("/orders/:order_uuid/items/:uuid", server.getOrderItem)
	authRouter.PUT("/orders/:order_uuid/items/:uuid", server.updateOrderItem)
	authRouter.DELETE("/orders/:order_uuid/items/:uuid", server.deleteOrderItem)

	router.POST("/api/users", server.registerUser)
	authRouter.GET("/users", server.allUsers)
	authRouter.GET("/users/:uuid", server.getUser)
	authRouter.PUT("/users/:uuid", server.updateUser)
	authRouter.DELETE("/users/:uuid", server.deleteUser)

	router.POST("/api/login", server.login)
}
