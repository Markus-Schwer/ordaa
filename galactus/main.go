package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"gitlab.com/sfz.aalen/hackwerk/dotinder/galactus/orders"
)

func main() {
	galactus := NewGalactus()
	galactus.start()
}

func NewGalactus() Galactus {
	mo := orders.NewMultiOrders()
	rest := NewRestInterface(&mo)
	actionChan := make(chan orders.OrderAction)
	queue := NewQueueClient(actionChan)
	return Galactus{
		mo:         &mo,
		rest:       &rest,
		actionChan: actionChan,
		queue:      &queue,
	}
}

type Galactus struct {
	mo         *orders.MultiOrders
	rest       *RestInterface
	queue      *QueueClient
	actionChan chan orders.OrderAction
}

func (gal *Galactus) start() {
	log.Println("galactus is starting")
	interfaceContext, cancel := context.WithCancel(context.Background())
	go gal.writeActionChanToOrders(interfaceContext)
	go gal.rest.start(interfaceContext)
	go gal.queue.start(interfaceContext)
	// channel to notify when shutdown should happen
	c := make(chan os.Signal, 1)
	// graceful shutdown on SIGINT and SIGTERM
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	// block until shutdown is desired
	<-c
	log.Println("shutting down")
	cancel()
	os.Exit(0)
}

func (gal *Galactus) writeActionChanToOrders(ctx context.Context) {
	done := ctx.Done()
	for {
		select {
		case <-done:
			return
		case action := <-gal.actionChan:
			gal.mo.HandleOrderAction(action)
		}
	}
}
