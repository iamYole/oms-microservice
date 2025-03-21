package main

import (
	"context"
	"encoding/json"
	"log"

	pb "github.com/iamYole/common/api"
	"github.com/iamYole/common/broker"
	amqp "github.com/rabbitmq/amqp091-go"
	"google.golang.org/grpc"
)

type grpcHandler struct {
	pb.UnimplementedOrderServiceServer
	service OrderService
	channel *amqp.Channel
}

func NewGRPCHandler(grpcServer *grpc.Server, service OrderService, channel *amqp.Channel) {
	handler := &grpcHandler{
		service: service,
		channel: channel,
	}
	pb.RegisterOrderServiceServer(grpcServer, handler)
}

func (h *grpcHandler) GetOrder(ctx context.Context, o *pb.GetOrderRequest) (*pb.Order, error) {

	order, err := h.service.GetOrder(ctx, o)
	if err != nil {
		return nil, err
	}

	return order, nil
}

func (h *grpcHandler) CreateOrder(ctx context.Context, p *pb.CreateOrderRequest) (*pb.Order, error) {
	log.Printf("New Order Received\n %v", p)

	items, err := h.service.ValidateOrder(ctx, p)
	if err != nil {
		return nil, err
	}

	order, err := h.service.CreateOrder(ctx, p,items)
	if err != nil {
		return nil, err
	}

	marshalledOrder, err := json.Marshal(order)
	if err != nil {
		return nil, err
	}

	//Send message to Queue
	q, err := h.channel.QueueDeclare(broker.OrderCreatedEvent, true, false, false, false, nil)
	if err != nil {
		log.Fatal(err)
	}

	h.channel.PublishWithContext(ctx, "", q.Name, false, false, amqp.Publishing{
		ContentType:  "application/json",
		Body:         marshalledOrder,
		DeliveryMode: amqp.Persistent,
	})

	return order, nil
}
