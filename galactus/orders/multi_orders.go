package orders

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
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

func (moo *MultiOrders) CreateNewOrder(provider string) int {
	moo.mu.Lock()
	defer moo.mu.Unlock()
	id := moo.nextId
	moo.nextId += 1
	moo.activeOrders[id] = newOrderHandler(provider)
	return id
}

func (moo *MultiOrders) GetOrder(orderNo int) (Order, string, error) {
	moo.mu.Lock()
	defer moo.mu.Unlock()
	oh, ok := moo.activeOrders[orderNo]
	if !ok {
		return nil, "", fmt.Errorf("there is no active order with no %d", orderNo)
	}
	return oh.getOrders(), oh.provider, nil
}

func (moo *MultiOrders) GetOrders() []map[string]interface{} {
	moo.mu.Lock()
	defer moo.mu.Unlock()
	orderMeta := make([]map[string]interface{}, 0)
	for orderNo, oh := range moo.activeOrders {
		orderMeta = append(orderMeta, map[string]interface{}{
			"orderNo":  orderNo,
			"provider": oh.provider,
		})
	}
	return orderMeta
}

// the first return value is an optional orders parameter
func (moo *MultiOrders) HandleOrderAction(ctx context.Context, action OrderAction) (Order, error) {
	moo.mu.Lock()
	defer moo.mu.Unlock()
	oh, ok := moo.activeOrders[action.OrderNo]
	if !ok {
		return nil, fmt.Errorf("no order with no %d found", action.OrderNo)
	}
	if action.Item != "" {
		err := checkItem(ctx, fmt.Sprintf("%s/%s/check", ctx.Value("OMEGA_STAR_URL").(string), oh.provider), action.Item)
		if err != nil {
			return nil, err
		}
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

func checkItem(ctx context.Context, url string, item string) error {
	b, err := json.Marshal([]string{item})
	if err != nil {
		return err
	}
	withTo, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()
	req, err := http.NewRequestWithContext(withTo, http.MethodPost, url, bytes.NewReader(b))
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
