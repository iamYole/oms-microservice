package main

import (
	"context"

	pb "github.com/iamYole/common/api"
	"github.com/iamYole/oms-payments/processor"
)

type service struct {
	processor processor.PaymentProcessor
}

func NewPaymentService(processor processor.PaymentProcessor) *service {
	return &service{processor: processor}
}

func (s *service) CreatePayment(ctx context.Context, order *pb.Order) (string, error) {
	//connect to payment processor
	link, err := s.processor.CreatePaymentLink(order)
	if err != nil {
		return "", err
	}

	return link, nil
}
