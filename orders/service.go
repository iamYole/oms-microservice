package main

import (
	"context"
	"log"

	"github.com/iamYole/common"
	pb "github.com/iamYole/common/api"
)

type service struct {
	store OrderStore
}

func NewService(store OrderStore) *service {
	return &service{store}
}

func (s *service) CreateOrder(context.Context) error {
	return nil
}

func (s *service) ValidateOrder(ctx context.Context, p *pb.CreateOrderRequest) error {
	if len(p.Items) == 0 {
		return common.ErrNoItems
	}

	mergedItema := mergeItemQuantities(p.Items)
	log.Println(mergedItema)

	//validate with the stock service

	return nil
}

func mergeItemQuantities(items []*pb.ItemsWithQuantity) []*pb.ItemsWithQuantity {
	merged := make([]*pb.ItemsWithQuantity, 0)

	for _, item := range items {
		found := false

		for _, finalitem := range merged {
			if finalitem.ID == item.ID {
				finalitem.Quantity += item.Quantity
				found = true
				break
			}

			if !found {
				merged = append(merged, item)
			}
		}

	}

	return merged
}
