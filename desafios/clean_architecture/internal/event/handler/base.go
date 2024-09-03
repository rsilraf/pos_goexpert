package handler

import (
	"encoding/json"
	"log"

	"github.com/streadway/amqp"
)

type RabbitMQhandler struct {
	Channel *amqp.Channel
}

func (r *RabbitMQhandler) PublishJSON(exchange string, key string, payload interface{}) {

	jsonOutput, _ := json.Marshal(payload)

	msg := amqp.Publishing{
		ContentType: "application/json",
		Body:        jsonOutput,
	}
	log.Printf("Publishing message to exchange %s with key %s\n", exchange, key)
	r.Channel.Publish(
		exchange, // exchange
		key,      // key name
		false,    // mandatory
		false,    // immediate
		msg,      // message
	)
}
