package main

import (
	"fmt"
	"sync"

	"github.com/rs/zerolog/log"
)

const (
	New      = "new"
	Add      = "add"
	Remove   = "remove"
	Finalize = "finalize"
	Cancel   = "cancel"
	Arrived  = "arrived"
)

func getSupportedActionVerbs() []string {
	return []string{New, Add, Remove, Finalize, Cancel, Arrived}
}

type Action struct {
	orderNo int
	user    string
	verb    string
	item    string
}

type OrderMetadata struct {
	OrderNo  int    `json:"orderNo"`
	Provider string `json:"provider"`
}

type ServiceFacade interface {
	CheckOrderItem(provider string, item string) error
	NewOrder(provider string) (int, error)
	AddOrderItem(orderNo int, user string, item string) error
	RemoveOrderItem(orderNo int, user string, item string) error
	FinalizeOrder(orderNo int) error
	OrderArrived(orderNo int) error
	CancelOrder(orderNo int) error
	GetOrders() ([]OrderMetadata, error)
}

func NewActionRunner(services ServiceFacade) *ActionRunner {
	return &ActionRunner{
		orders:     make(map[string]int),
		orderMutex: sync.Mutex{},
		services:   services,
	}
}

type ActionRunner struct {
	orders     map[string]int
	orderMutex sync.Mutex
	services   ServiceFacade
}

func (runner *ActionRunner) runAction(user string, action *ParsedAction) error {
	log.Debug().Msgf("running action %v for user %s", action, user)
	runner.orderMutex.Lock()
	defer runner.orderMutex.Unlock()
	orderNo, hasOrder := runner.orders[action.provider]
	if !hasOrder {
		log.Warn().Msgf("no order for provider %s in cache, will try to find one in galactus", action.provider)
		ordersMeta, err := runner.services.GetOrders()
		if err != nil {
			return err
		}
		log.Debug().Msgf("loaded active order from galactus: %v", ordersMeta)
		for _, orderMeta := range ordersMeta {
			if orderMeta.Provider == action.provider {
				log.Debug().Msgf("found matching active order: %d", orderMeta.OrderNo)
				runner.orders[orderMeta.Provider] = orderMeta.OrderNo
				orderNo = orderMeta.OrderNo
				hasOrder = true
				break
			}
		}
	}

	switch action.verb {
	case New:
		if hasOrder {
			return fmt.Errorf("there is already an order for provider '%s'", action.provider)
		}
		orderNo, err := runner.services.NewOrder(action.provider)
		if err != nil {
			return err
		}
		runner.orders[action.provider] = orderNo
	case Add:
		return runner.services.AddOrderItem(orderNo, user, action.item)
	case Remove:
		return runner.services.RemoveOrderItem(orderNo, user, action.item)
	case Finalize:
		return runner.services.FinalizeOrder(orderNo)
	case Cancel:
		delete(runner.orders, action.provider)
		return runner.services.CancelOrder(orderNo)
	case Arrived:
		delete(runner.orders, action.provider)
		return runner.services.OrderArrived(orderNo)
	default:
		return fmt.Errorf("action verb '%s' not supported", action.verb)
	}
	return nil
}
