package main

import (
	"bytes"
	"context"
	"fmt"
	"sync"
	"text/template"

	"github.com/rodaine/table"
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
	Price int
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

func (runner *ActionRunner) runAction(user string, action *ParsedAction) (messages []string, err error) {
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
		messages = []string{
			ordersToTable(orders, menu),
			ordersToOrderMessage(orders, menu, "Hackwerk 1, 73430 Aalen"),
		}
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

func ordersToTable(orders Orders, menu Menu) string {
	t := table.New("User", "Id", "Name", "Price")
	for user, o := range orders {
		var total int
		for itemIdx, id := range o {
		PerItem:
			for _, item := range menu.Items {
				if id == item.Id {
					userCol := ""
					if itemIdx == 0 {
						userCol = user
					}
					t.AddRow(
						userCol,
						fmt.Sprintf("%3s", item.Id),
						item.Name,
						fmt.Sprintf("%2d.%02d€", item.Price/100, item.Price%100),
					)
					total += item.Price
					break PerItem
				}
			}
		}
		t.AddRow("", "", "total", fmt.Sprintf("%2d.%02d€", total/100, total%100))
	}
	var buf bytes.Buffer
	t.WithWriter(&buf)
	t.Print()
	return buf.String()
}

func ordersToOrderMessage(orders Orders, menu Menu, adress string) string {
	orderCount := make(map[string]int)
	for _, o := range orders {
		for _, id := range o {
			orderCount[id] += 1
		}
	}
	t := table.New("", "")
	for _, item := range menu.Items {
		if orderCount[item.Id] == 0 {
			continue
		}
		t.AddRow(fmt.Sprint(orderCount[item.Id])+"x", fmt.Sprintf("%s (%s)", item.Name, item.Id))
	}
	var tableBuf bytes.Buffer
	t.WithWriter(&tableBuf)
	t.Print()
	tpl, err := template.New("order").Parse("Hallo, ich möchte gerne eine Bestellung aufgeben. Lieferadresse ist {{.adress}}.\n{{.orderString}}\nVielen Dank")
	if err != nil {
		panic(err)
	}
	var templateBuf bytes.Buffer
	err = tpl.Execute(&templateBuf, map[string]interface{}{
		"adress":      adress,
		"orderString": tableBuf.String(),
	})
	if err != nil {
		panic(err)
	}
	return templateBuf.String()
}
