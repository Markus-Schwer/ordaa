package orders

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
)

type OrderAction struct {
	Action  string `json:"action"`
	User    string `json:"user"`
	Item    string `json:"item"`
	OrderNo int    `json:"orderNo"`
}

// thread safe object to manage orders
type MultiOrders struct {
	activeOrders map[int]*orderHandler
	nextId       int
	mu           sync.Mutex
}

func NewMultiOrders() MultiOrders {
	return MultiOrders{
		activeOrders: make(map[int]*orderHandler),
		nextId:       1,
		mu:           sync.Mutex{},
	}
}

func (moo *MultiOrders) CreateNewOrder() int {
	moo.mu.Lock()
	defer moo.mu.Unlock()
	id := moo.nextId
	moo.nextId += 1
	moo.activeOrders[id] = newOrderHandler()
	return id
}

func (moo *MultiOrders) GetOrder(orderNo int) (Order, error) {
	moo.mu.Lock()
	defer moo.mu.Unlock()
	oh, ok := moo.activeOrders[orderNo]
	if !ok {
		return nil, nil
	}
	return oh.getOrders(), nil
}

// the first return value is an optional orders parameter
func (moo *MultiOrders) HandleOrderAction(ctx context.Context, action OrderAction) (Order, error) {
	if action.Item != "" {
		err := checkItem(fmt.Sprintf("%s/sangam/check", ctx.Value("OMEGA_STAR_URL").(string)), action.Item)
		if err != nil {
			return nil, err
		}
	}
	moo.mu.Lock()
	defer moo.mu.Unlock()
	oh, ok := moo.activeOrders[action.OrderNo]
	if !ok {
		return nil, fmt.Errorf("no order with no %d found", action.OrderNo)
	}
	switch action.Action {
	case "add":
		return nil, oh.addItem(action.User, action.Item)
	case "remove":
		return nil, oh.removeItem(action.User, action.Item)
	case "finalize":
		return oh.finalize()
	case "arrived":
		return nil, oh.arrived()
	case "cancel":
		return nil, oh.cancel()
	case "":
		return nil, fmt.Errorf("invalid empty action")
	default:
		return nil, fmt.Errorf("unknown action: %s", action.Action)
	}
}

func checkItem(url string, item string) error {
	b, err := json.Marshal([]string{item})
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(b))
	if err != nil {
		return err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	b, err = io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	var data []string
	err = json.Unmarshal(b, &data)
	if err != nil {
		return err
	}
	if len(data) == 0 {
		return nil
	}
	return fmt.Errorf("invalid order item: %s", string(b))
}
