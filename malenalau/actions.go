package main

import (
	"context"
	"fmt"
	"strings"
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

type Orders map[string][]string

type MenuItem struct {
	Id    string
	Name  string
	Price float64
}

type Menu struct {
	Items []MenuItem
}

type ServiceFacade interface {
	GetMenu(provider string) (Menu, error)
	CheckOrderItem(provider string, item string) error
	NewOrder(provider string) (int, error)
	AddOrderItem(orderNo int, user string, item string) error
	RemoveOrderItem(orderNo int, user string, item string) error
	FinalizeOrder(orderNo int) (Orders, error)
	OrderArrived(orderNo int) error
	CancelOrder(orderNo int) error
	GetOrders() ([]OrderMetadata, error)
}

func NewActionRunner(ctx context.Context, services ServiceFacade) *ActionRunner {
	return &ActionRunner{
		ctx:        ctx,
		orders:     make(map[string]int),
		orderMutex: sync.Mutex{},
		services:   services,
	}
}

type ActionRunner struct {
	ctx        context.Context
	orders     map[string]int
	orderMutex sync.Mutex
	services   ServiceFacade
}

func (runner *ActionRunner) runAction(user string, action *ParsedAction) (message string, err error) {
	log.Ctx(runner.ctx).Debug().Msgf("running action %v for user %s", action, user)
	runner.orderMutex.Lock()
	defer runner.orderMutex.Unlock()
	orderNo, hasOrder := runner.orders[action.provider]
	if !hasOrder {
		log.Ctx(runner.ctx).Warn().Msgf("no order for provider %s in cache, will try to find one in galactus", action.provider)
		var ordersMeta []OrderMetadata
		ordersMeta, err = runner.services.GetOrders()
		if err != nil {
			return
		}
		log.Ctx(runner.ctx).Debug().Msgf("loaded active order from galactus: %v", ordersMeta)
		for _, orderMeta := range ordersMeta {
			if orderMeta.Provider == action.provider {
				log.Ctx(runner.ctx).Debug().Msgf("found matching active order: %d", orderMeta.OrderNo)
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
			err = fmt.Errorf("there is already an order for provider '%s'", action.provider)
			return
		}
		orderNo, err = runner.services.NewOrder(action.provider)
		if err != nil {
			return
		}
		runner.orders[action.provider] = orderNo
	case Add:
		err = runner.services.AddOrderItem(orderNo, user, action.item)
		return
	case Remove:
		err = runner.services.RemoveOrderItem(orderNo, user, action.item)
		return
	case Finalize:
		var orders Orders
		orders, err = runner.services.FinalizeOrder(orderNo)
		var menu Menu
		menu, err = runner.services.GetMenu(action.provider)
		message = ordersToTable(orders, menu)
		return
	case Cancel:
		delete(runner.orders, action.provider)
		err = runner.services.CancelOrder(orderNo)
		return
	case Arrived:
		delete(runner.orders, action.provider)
		err = runner.services.OrderArrived(orderNo)
		return
	default:
		err = fmt.Errorf("action verb '%s' not supported", action.verb)
		return
	}
	return
}

func ordersToTable(orders Orders, menu Menu) (table string) {
	for user, o := range orders {
		verboseItems := make([]string, len(o))
		var total float64
		for i, id := range o {
		PerItem:
			for _, item := range menu.Items {
				if id == item.Id {
					verboseItems[i] = fmt.Sprintf("%s: (%s) %f", item.Id, item.Name, item.Price)
					total += float64(item.Price)
					break PerItem
				}
			}
		}
		table += fmt.Sprintf("%s\t\t%s = %f€\n", user, strings.Join(verboseItems, " + "), total)
	}
	return
}
