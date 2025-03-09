package stripeprocessor

import (
	"fmt"
	"log"

	"github.com/iamYole/common"
	pb "github.com/iamYole/common/api"
	"github.com/stripe/stripe-go/v81"
	"github.com/stripe/stripe-go/v81/checkout/session"
)

var gatewayHTTPAddr = common.GetString("GATEWAY_HTTPADDR", "http://localhost:8080")

type Stripe struct{}

func NewProcessor() *Stripe {
	return &Stripe{}
}

func (s *Stripe) CreatePaymentLink(order *pb.Order) (string, error) {
	log.Printf("Creating payment link for Customer: %s with order: %s", order.CustomerID, order.ID)
	gateWay_successURL := fmt.Sprintf("%s/success.html?customerID=%s&orderID=%s", gatewayHTTPAddr, order.CustomerID, order.ID)

	items := []*stripe.CheckoutSessionLineItemParams{}
	for _, i := range order.Items {
		items = append(items, &stripe.CheckoutSessionLineItemParams{
			Price:    stripe.String(i.PriceID),
			Quantity: stripe.Int64(int64(i.Quantity)),
		})
	}

	params := &stripe.CheckoutSessionParams{
		LineItems:  items,
		Mode:       stripe.String(string(stripe.CheckoutSessionModePayment)),
		SuccessURL: stripe.String(gateWay_successURL),
	}

	result, err := session.New(params)
	if err != nil {
		return "", err
	}

	return result.URL, nil
}
