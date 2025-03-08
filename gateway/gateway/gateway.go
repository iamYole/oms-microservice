package gateway

import (
	"context"

	pb "github.com/iamYole/common/api"
)

type OrdersGateway interface {
	CreateOrder(ctx context.Context, p *pb.CreateOrderRequest) (*pb.Order, error)
}
