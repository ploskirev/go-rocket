package orderservice

import (
	"context"
	"log"

	"github.com/ploskirev/go-rocket/order/internal/model"
)

func (os *orederService) PayOrder(ctx context.Context, paymentInfo *model.PaymentInfo) (string, error) {
	orderInfo, err := os.or.GetOrder(ctx, paymentInfo.OrderUUID)
	if err != nil {
		log.Printf("Error order not found")
		return "", err
	}

	transactionUUID, err := os.pc.PayOrder(ctx, &model.PaymentInfo{
		OrderUUID:     paymentInfo.OrderUUID,
		UserUUID:      orderInfo.UserUUID,
		PaymentMethod: paymentInfo.PaymentMethod,
	})
	if err != nil {
		log.Printf("Error pay order")
		return "", err
	}

	if err = os.or.UpdateOrder(ctx, orderInfo.OrderUUID, &model.Order{
		OrderUUID:       orderInfo.OrderUUID,
		UserUUID:        orderInfo.UserUUID,
		PartUUIDs:       orderInfo.PartUUIDs,
		TotalPrice:      orderInfo.TotalPrice,
		TransactionUUID: &transactionUUID,
		PaymentMethod:   &paymentInfo.PaymentMethod,
		Status:          model.PAID,
	}); err != nil {
		log.Printf("Error update order")
		return "", err
	}

	return transactionUUID, nil
}
