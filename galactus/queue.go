package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"gitlab.com/sfz.aalen/hackwerk/dotinder/galactus/orders"
)

func NewQueueClient(handler chan<- orders.OrderAction) QueueClient {
	return QueueClient{
		handler: handler,
	}
}

type QueueClient struct {
	conn    *amqp.Connection
	handler chan<- orders.OrderAction
}

func (qc *QueueClient) start(ctx context.Context) error {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		return err
	}
	qc.conn = conn
	defer qc.conn.Close()
	ch, err := qc.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()
	q, err := ch.QueueDeclare(
		"order/action", // name
		false,             // durable
		false,             // delete when unused
		false,             // exclusive
		false,             // no-wait
		nil,               // arguments
	)
	if err != nil {
		return err
	}
	msgs, err := ch.Consume(
		q.Name, // queue
		"",              // consumer
		true,            // auto-ack
		false,           // exclusive
		false,           // no-local
		false,           // no-wait
		nil,             // args
	)
	if err != nil {
		return err
	}
	for {
		select {
		case d := <-msgs:
			var action orders.OrderAction
			if err := json.Unmarshal(d.Body, &action); err != nil {
				log.Printf("received invalid message body '%s'", string(d.Body))
				continue
			}
			qc.handler <- action
		case <-ctx.Done():
			return fmt.Errorf("channel closed unexpectedly, retry could help")
		}
		if msgs == nil {
			return fmt.Errorf("channel closed unexpectedly, retry could help")
		}
	}
}

func (qc *QueueClient) writeOrder(orders orders.Order) error {
	ch, err := qc.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()
	q, err := ch.QueueDeclare(
		"order/finalized", // name
		false,             // durable
		false,             // delete when unused
		false,             // exclusive
		false,             // no-wait
		nil,               // arguments
	)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	b, err := json.Marshal(orders)
	if err != nil {
		return err
	}
	err = ch.PublishWithContext(ctx,
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        b,
		})
	return err
}
