package main

import (
	"context"

	"github.com/iamYole/common"
	pb "github.com/iamYole/common/api"
)

type service struct {
	store OrderStore
}

func NewService(store OrderStore) *service {
	return &service{store}
}

func (s *service) GetOrder(ctx context.Context, o *pb.GetOrderRequest) (*pb.Order, error) {
	return s.store.Get(ctx, o.OrderID, o.CustomerID)
}

func (s *service) CreateOrder(ctx context.Context, p *pb.CreateOrderRequest, items []*pb.Item) (*pb.Order, error) {
	id, err := s.store.Create(ctx, p, items)
	if err != nil {
		return nil, err
	}

	order := &pb.Order{
		ID:        id,
		CustomerID: p.CustomerID,
		Status:     "pendind",
		Items:      items,
	}

	return order, nil
}

func (s *service) ValidateOrder(ctx context.Context, p *pb.CreateOrderRequest) ([]*pb.Item, error) {
	if len(p.Items) == 0 {
		return nil, common.ErrNoItems
	}

	mergedItems := mergeItemsQuantities(p.Items)
	//log.Print(mergedItems)

	//temp
	var itemsWithPrice []*pb.Item
	for _, i := range mergedItems {
		itemsWithPrice = append(itemsWithPrice, &pb.Item{
			PriceID:  "price_1R0Yd4KkSopfABcRq7qflVBQ",
			ID:       i.ID,
			Name:     "Jollof Rice",
			Quantity: i.Quantity,
		})
	}

	//validate with the stock service

	return itemsWithPrice, nil
}

func mergeItemsQuantities(items []*pb.ItemsWithQuantity) []*pb.ItemsWithQuantity {
	merged := make([]*pb.ItemsWithQuantity, 0)

	for _, item := range items {
		found := false

		for _, finalitem := range merged {
			if finalitem.ID == item.ID {
				finalitem.Quantity += item.Quantity
				found = true
				break
			}
		}

		if !found {
			merged = append(merged, item)
		}
	}

	return merged
}
