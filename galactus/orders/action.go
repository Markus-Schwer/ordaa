package orders

import "github.com/google/uuid"

type OrderAction struct {
	Action   string    `json:"action"`
	User     string    `json:"user"`
	Item     string    `json:"item"`
	OrderNo  int       `json:"orderNo"`
	Provider string    `json:"provider"`
	Uuid     uuid.UUID `json:"uuid"`
}

type ActiveOrder struct {
	Provider string `json:"provider"`
	OrderNo  int    `json:"orderNo"`
}

func (action *OrderAction) respondMissingOrderNo() OrderActionResponse {
	return &NoActiveOrder{uuid: action.Uuid, orderNo: action.OrderNo}
}

func (action *OrderAction) respondOkWithOrderNo(orderNo int) OrderActionResponse {
	return &OkWithOrderNo{uuid: action.Uuid, OrderNo: orderNo}
}

func (action *OrderAction) respondOkWithOrder(order *Order) OrderActionResponse {
	return &OkWithOrder{uuid: action.Uuid, Order: order}
}

func (action *OrderAction) respondWithPossibleError(err error) (res OrderActionResponse) {
	if err != nil {
		res = &GenericError{uuid: action.Uuid, err: err}
	} else {
		res = &Ok{uuid: action.Uuid}
	}
	return
}

func (action *OrderAction) respondGenericError(err error) OrderActionResponse {
	return &GenericError{uuid: action.Uuid, err: err}
}

func (action *OrderAction) respondWithActiveOrders(activeOrders []ActiveOrder) OrderActionResponse {
	return &OkWithActiveOrders{uuid: action.Uuid, ActiveOrders: activeOrders}
}
