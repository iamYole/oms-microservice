package main

import (
	"context"
	"testing"

	"github.com/iamYole/common/api"
	"github.com/iamYole/oms-payments/processor/inmem"
)

func TestPayment(t *testing.T) {
	processor := inmem.NewInmem()
	svc := NewPaymentService(processor)

	t.Run("should create a payment link", func(t *testing.T) {
		link, err := svc.CreatePayment(context.Background(), &api.Order{})
		if err != nil {
			t.Errorf("Create payment error = %v, should be nil", err)
		}
		if link == "" {
			t.Errorf("Create payment link is empty")
		}
	})
}
