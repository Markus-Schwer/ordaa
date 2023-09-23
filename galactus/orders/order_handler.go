package orders

import "fmt"

const (
	// states
	takeOrders = iota
	ordered
	idle
	// transitions
	addItem
	removeItem
	finalize
	cancel
	arrived
)

var transitionTable = map[int]map[int]int{
	takeOrders: {
		addItem:    takeOrders,
		removeItem: takeOrders,
		finalize:   ordered,
		cancel:     idle,
	},
	ordered: {
		cancel:  idle,
		arrived: idle,
	},
	idle: {},
}

// Order map from usernames to item ids
type Order map[string][]string

type orderHandler struct {
	provider     string
	currentState int
	// map of users to their orders
	orders Order
}

func newOrderHandler(provider string) *orderHandler {
	return &orderHandler{
		provider:     provider,
		currentState: takeOrders,
		orders:       make(map[string][]string),
	}
}

func (sm *orderHandler) addItem(user string, item string) error {
	if newState, ok := transitionTable[sm.currentState][addItem]; ok {
		sm.currentState = newState
	} else {
		return fmt.Errorf("cannot add item in state %d", sm.currentState)
	}
	if userOrders, ok := sm.orders[user]; ok {
		sm.orders[user] = append(userOrders, item)
	} else {
		sm.orders[user] = []string{item}
	}
	return nil
}

func (sm *orderHandler) removeItem(user string, item string) error {
	if newState, ok := transitionTable[sm.currentState][removeItem]; ok {
		sm.currentState = newState
	} else {
		return fmt.Errorf("cannot remove item in state %d", sm.currentState)
	}
	if userOrders, ok := sm.orders[user]; ok {
		for i, order := range userOrders {
			if order == item {
				// remove element from the list
				userOrders[i] = userOrders[len(userOrders)-1]
				sm.orders[user] = userOrders[:len(userOrders)-1]
				return nil
			}
		}
		return fmt.Errorf("could not remove %s, user %s did not order it", item, user)
	} else {
		return fmt.Errorf("could not remove %s, user %s did not order yet", item, user)
	}
}

func (sm *orderHandler) finalize() (Order, error) {
	if newState, ok := transitionTable[sm.currentState][finalize]; ok {
		sm.currentState = newState
	} else {
		return nil, fmt.Errorf("cannot finalize in state %d", sm.currentState)
	}
	return sm.orders, nil
}

func (sm *orderHandler) cancel() error {
	if newState, ok := transitionTable[sm.currentState][cancel]; ok {
		sm.currentState = newState
		return nil
	} else {
		return fmt.Errorf("cannot cancel in state %d", sm.currentState)
	}
}

func (sm *orderHandler) arrived() error {
	if newState, ok := transitionTable[sm.currentState][arrived]; ok {
		sm.currentState = newState
		return nil
	} else {
		return fmt.Errorf("order cannot arrive in state %d", sm.currentState)
	}
}

func (sm *orderHandler) getOrders() Order {
	return sm.orders
}

func (sm *orderHandler) getState() string {
	switch sm.currentState {
	case takeOrders:
		return "taking orders"
	case ordered:
		return "ordered. delivery is on its way"
	case idle:
		return "idle. order was cancelled or deliverd"
	default:
		return "error: unknown state, there be dragons"

	}
}
