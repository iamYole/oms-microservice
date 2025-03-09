package main

import (
	"context"

	pb "github.com/iamYole/common/api"
)

type PaymentServices interface {
	CreatePayment(context.Context, *pb.Order) (string, error)
}
