package handler

import (
	"fmt"
	"sync"

	"github.com/rsilraf/pos_goexpert/desafios/clean_architecture/pkg/events"
	"github.com/streadway/amqp"
)

type OrderCreatedHandler struct {
	RabbitMQhandler
}

func NewOrderCreatedHandler(rabbitMQChannel *amqp.Channel) *OrderCreatedHandler {
	return &OrderCreatedHandler{
		RabbitMQhandler{rabbitMQChannel},
	}
}

func (h *OrderCreatedHandler) Handle(event events.EventInterface, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Printf("Order created: %v\n", event.GetPayload())

	h.PublishJSON(
		"amq.direct",       // exchange
		"created",          // key
		event.GetPayload(), // message to publish
	)
}
