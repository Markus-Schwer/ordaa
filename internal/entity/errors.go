package entity

import "errors"

var (
	ErrCannotCreatUUID                 = errors.New("cannot create uuid")
	ErrCannotGetAllOrder               = errors.New("could not get all orders from db")
	ErrOrderNotFound                   = errors.New("order not found")
	ErrGettingOrder                    = errors.New("could not get order from db")
	ErrCannotGetAllOrders              = errors.New("could not get all orders from db")
	ErrCannotGetAllOrderItems          = errors.New("could not get all order items from db")
	ErrGettingOrderItemsOrderAndUser   = errors.New("could not get order items for order and user")
	ErrGettingOrderItem                = errors.New("error getting order item")
	ErrOrderItemNotFound               = errors.New("order item not found")
	ErrCreatingOrderItem               = errors.New("could not create order item")
	ErrOrderNotOpen                    = errors.New("order is not in state open")
	ErrMenuItemUUIDMissing             = errors.New("menu item uuid missing")
	ErrCreatingOrder                   = errors.New("could not create order")
	ErrActiveOrderForMenuAlreadyExists = errors.New("there is already an active order the specified menu")
	ErrUpdatingOrder                   = errors.New("could not update order")
	ErrSugarPersonChangeForbidden      = errors.New("changing sugar person after it has already been set is forbidden")
	ErrSugarPersonNotSet               = errors.New("the sugar person has not been set")
	ErrOrderStateTransitionInvalid     = errors.New("invalid order state transition")
	ErrUpdatingOrderItem               = errors.New("could not update order item")
	ErrOrderUUIDChangeForbidden        = errors.New("changing order uuid is forbidden")
	ErrMenuItemUUIDChangeForbidden     = errors.New("changing menu item uuid is forbidden")
	ErrUserChangeForbidden             = errors.New("changing user is forbidden")
	ErrPaidChangeForbidden             = errors.New("paid status can only be changed by sugar person")
	ErrDeletingOrderItem               = errors.New("could not delete order item")
	ErrCannotGetAllUsers               = errors.New("could not get all users from db")
	ErrGettingUser                     = errors.New("failed to get user")
	ErrUserNotFound                    = errors.New("user not found")
	ErrCreatingUser                    = errors.New("could not create user")
	ErrUpdatingUser                    = errors.New("could not update users")
	ErrDeletingUser                    = errors.New("could not delete user")
	ErrCannotGetAllMenus               = errors.New("could not get all menus from db")
	ErrGettingMenu                     = errors.New("failed to get menu")
	ErrMenuNotFound                    = errors.New("menu not found")
	ErrGettingMenuItem                 = errors.New("failed to get menu item")
	ErrMenuItemNotFound                = errors.New("menu item not found")
	ErrCreatingMenu                    = errors.New("could not create menu")
	ErrUpdatingMenu                    = errors.New("could not update menu")
	ErrCreatingMenuItem                = errors.New("could not create menu item")
	ErrDeletingMenuItem                = errors.New("could not delete menu item")
	ErrDeletingMenu                    = errors.New("could not delete menu")
	ErrMigrationPreparationFailed      = errors.New("database migration preparation failed")
	ErrMigrationFailed                 = errors.New("database migration execution failed")
	ErrSettingPublicKey                = errors.New("setting public key for user")
)
