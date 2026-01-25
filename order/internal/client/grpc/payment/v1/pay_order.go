package paymentv1

import (
	"context"
	"log"

	"github.com/ploskirev/go-rocket/order/internal/client/converter"
	"github.com/ploskirev/go-rocket/order/internal/model"
)

func (c *paymentClient) PayOrder(ctx context.Context, paymentInfo *model.PaymentInfo) (string, error) {
	res, err := c.pc.PayOrder(ctx, converter.PaymentInfoToProto(paymentInfo))
	if err != nil {
		log.Printf("ERROR: Pay order from payment service client: %s", err)
		return "", err
	}

	return res.TransactionUuid, nil
}
