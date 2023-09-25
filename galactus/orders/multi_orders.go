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

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// thread safe object to manage orders
type MultiOrders struct {
	ctx          context.Context
	activeOrders map[int]*orderHandler
	nextId       int
	mu           sync.RWMutex
	out          chan<- OrderActionResponse
}

func NewMultiOrders(ctx context.Context) (*MultiOrders, chan OrderActionResponse) {
	out := make(chan OrderActionResponse, 10)
	mo := &MultiOrders{
		ctx:          ctx,
		activeOrders: make(map[int]*orderHandler),
		nextId:       1,
		mu:           sync.RWMutex{},
		out:          out,
	}
	return mo, out
}

func (moo *MultiOrders) Start(in <-chan OrderAction) {
	log.Ctx(moo.ctx).Info().Msg("Orders handler starts listening")
	for {
		select {
		case <-moo.ctx.Done():
			log.Ctx(moo.ctx).Info().Msg("Orders handler is shutting down")
			return
		case action := <-in:
			go moo.handleOrderAction(action)
		}
	}
}

func (moo *MultiOrders) handleOrderAction(action OrderAction) {
	log.Ctx(moo.ctx).Debug().Dict("action", zerolog.Dict().
		Str("verb", action.Action).
		Str("user", action.User)).
		Msg("handling action")
	if action.Action == "new" {
		// TODO: validate provider
		moo.mu.Lock()
		log.Ctx(moo.ctx).Debug().Msg("acquired lock")
		action.OrderNo = moo.nextId
		moo.nextId += 1
		moo.activeOrders[action.OrderNo] = newOrderHandler(action.Provider)
		moo.mu.Unlock()
		log.Ctx(moo.ctx).Debug().Msg("released lock")
		moo.out <- action.respondOkWithOrderNo(action.OrderNo)
		return
	}
	if action.Action == "active" {
		moo.mu.RLock()
		log.Ctx(moo.ctx).Debug().Msg("acquired lock")
		activeOrders := make([]ActiveOrder, 0)
		for orderNo, oh := range moo.activeOrders {
			activeOrders = append(activeOrders, ActiveOrder{
				OrderNo:  orderNo,
				Provider: oh.provider,
			})
		}
		moo.mu.RUnlock()
		log.Ctx(moo.ctx).Debug().Msg("released lock")
		moo.out <- action.respondWithActiveOrders(activeOrders)
		return
	}
	oh, ok := moo.activeOrders[action.OrderNo]
	if !ok {
		moo.out <- action.respondMissingOrderNo()
		return
	}
	action.Provider = oh.provider
	switch action.Action {
	case "status":
		moo.mu.RLock()
		moo.out <- action.respondOkWithOrder(&oh.orders)
		moo.mu.RUnlock()
	case "add":
		moo.mu.Lock()
		moo.out <- action.respondWithPossibleError(oh.addItem(action.User, action.Item))
		moo.mu.Unlock()
	case "remove":
		moo.mu.Lock()
		moo.out <- action.respondWithPossibleError(oh.removeItem(action.User, action.Item))
		moo.mu.Unlock()
	case "finalize":
		moo.mu.Lock()
		order, err := oh.finalize()
		moo.mu.Unlock()
		if err != nil {
			moo.out <- action.respondGenericError(err)
		} else {
			moo.out <- action.respondOkWithOrder(&order)
		}
	case "arrived":
		moo.mu.Lock()
		delete(moo.activeOrders, action.OrderNo)
		moo.mu.Unlock()
		moo.out <- action.respondWithPossibleError(oh.arrived())
	case "cancel":
		moo.mu.Lock()
		log.Ctx(moo.ctx).Debug().Msg("acquired lock")
		delete(moo.activeOrders, action.OrderNo)
		moo.mu.Unlock()
		log.Ctx(moo.ctx).Debug().Msg("released lock")
		moo.out <- action.respondWithPossibleError(oh.cancel())
	case "":
		moo.out <- action.respondGenericError(fmt.Errorf("invalid empty action"))
	default:
		moo.out <- action.respondGenericError(fmt.Errorf("unknown action: %s", action.Action))
	}
}

func (moo *MultiOrders) checkItems(action *OrderAction) error {
	if action.Item == "" {
		return nil
	}
	url := fmt.Sprintf("%s/%s/check", moo.ctx.Value("OmegaStarURL").(string), action.Provider)
	b, err := json.Marshal([]string{action.Item})
	if err != nil {
		return err
	}
	withTo, cancel := context.WithTimeout(moo.ctx, time.Second*5)
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
	var invalidItems []string
	err = json.Unmarshal(b, &invalidItems)
	if err != nil {
		return err
	}
	if len(invalidItems) == 0 {
		return nil
	}
	return fmt.Errorf("invalid order item(s): %v", invalidItems)
}
