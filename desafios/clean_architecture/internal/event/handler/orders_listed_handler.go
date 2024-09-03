package handler

import (
	"fmt"
	"sync"

	"github.com/rsilraf/pos_goexpert/desafios/clean_architecture/pkg/events"
	"github.com/streadway/amqp"
)

type OrdersListedHandler struct {
	RabbitMQhandler
}

func NewOrdersListedHandler(rabbitMQChannel *amqp.Channel) *OrdersListedHandler {
	return &OrdersListedHandler{
		RabbitMQhandler{rabbitMQChannel},
	}
}

func (h *OrdersListedHandler) Handle(event events.EventInterface, wg *sync.WaitGroup) {
	defer wg.Done()
	// fmt.Printf("Orders listed: %v", event.GetPayload())
	fmt.Println("Orders listed")

	h.PublishJSON(
		"amq.direct",       // exchange
		"listed",           // key
		event.GetPayload(), // message to publish
	)
}
