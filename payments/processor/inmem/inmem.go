package inmem

import pb "github.com/iamYole/common/api"

type inmem struct {
}

func NewInmem() *inmem {
	return &inmem{}
}

func (i *inmem) CreatePaymentLink(order *pb.Order) (string, error) {
	return "dummy_link", nil
}
