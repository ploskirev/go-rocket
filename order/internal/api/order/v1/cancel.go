package orderapiv1

import (
	"context"
	"log"

	order_v1 "github.com/ploskirev/go-rocket/shared/pkg/openapi/order/v1"
)

func (a *api) CancelOrder(ctx context.Context, params order_v1.CancelOrderParams) (order_v1.CancelOrderRes, error) {
	if err := a.os.CancelOrder(ctx, params.OrderUUID); err != nil {
		log.Printf("Error cancel order")
		return &order_v1.InternalServerError{
			Code:    500,
			Message: "Error cancel order",
		}, nil
	}

	return &order_v1.CancelOrderNoContent{}, nil
}
