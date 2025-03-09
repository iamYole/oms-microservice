package main

import (
	"context"
	"encoding/json"
	"log"

	pb "github.com/iamYole/common/api"
	"github.com/iamYole/common/broker"
	amqp "github.com/rabbitmq/amqp091-go"
)

type consumer struct {
	service PaymentServices
}

func NewConsumer(service PaymentServices) *consumer {
	return &consumer{service: service}
}

func (c *consumer) Listen(ch *amqp.Channel) {
	q, err := ch.QueueDeclare(broker.OrderCreatedEvent, true, false, false, false, nil)
	if err != nil {
		log.Fatal(err)
	}

	msgs, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		log.Fatal(err)
	}

	var forever chan struct{}
	go func() {
		for d := range msgs {
			log.Printf("Received message %s", d.Body)

			order := &pb.Order{}
			if err := json.Unmarshal(d.Body, order); err != nil {
				log.Printf("Failed to unmarshal order: %v", err)
				continue
			}

			paymentLink, err := c.service.CreatePayment(context.Background(), order)
			if err != nil {
				log.Printf("failed to create payment link: %v", err)
				continue
			}

			log.Printf("Payment link created %s", paymentLink)
		}
	}()
	<-forever
}
