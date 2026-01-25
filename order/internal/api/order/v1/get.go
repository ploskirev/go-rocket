package orderapiv1

import (
	"context"
	"fmt"
	"log"

	"github.com/ploskirev/go-rocket/order/internal/api/order/v1/converter"
	order_v1 "github.com/ploskirev/go-rocket/shared/pkg/openapi/order/v1"
)

func (a *api) GetOrder(ctx context.Context, params order_v1.GetOrderParams) (order_v1.GetOrderRes, error) {
	order, err := a.os.GetOrder(ctx, params.OrderUUID)
	if err != nil {
		log.Printf("Error get order: %s", err)
		return &order_v1.InternalServerError{
			Code:    500,
			Message: fmt.Sprintf("Error get order: %s", err),
		}, nil
	}

	return converter.OrderToApi(order), nil
}
