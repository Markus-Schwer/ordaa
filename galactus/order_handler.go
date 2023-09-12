package main

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

type OrderHandler struct {
	currentState int
	// map of users to their orders
	orders map[string][]string
}

func NewOrderHandler() *OrderHandler {
	return &OrderHandler{
		currentState: takeOrders,
		orders:       make(map[string][]string),
	}
}

func (sm *OrderHandler) addItem(user string, item string) error {
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

func (sm *OrderHandler) removeItem(user string, item string) error {
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

func (sm *OrderHandler) finalize() (map[string][]string, error) {
	if newState, ok := transitionTable[sm.currentState][finalize]; ok {
		sm.currentState = newState
	} else {
		return nil, fmt.Errorf("cannot finalize in state %d", sm.currentState)
	}
	return sm.orders, nil
}

func (sm *OrderHandler) cancel() error {
	if newState, ok := transitionTable[sm.currentState][cancel]; ok {
		sm.currentState = newState
		return nil
	} else {
		return fmt.Errorf("cannot cancel in state %d", sm.currentState)
	}
}

func (sm *OrderHandler) arrived() error {
	if newState, ok := transitionTable[sm.currentState][arrived]; ok {
		sm.currentState = newState
		return nil
	} else {
		return fmt.Errorf("order cannot arrive in state %d", sm.currentState)
	}
}

func (sm *OrderHandler) getOrders() map[string][]string {
	return sm.orders
}

func (sm *OrderHandler) getState() string {
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
