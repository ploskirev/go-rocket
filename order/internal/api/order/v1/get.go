package orderapiv1

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/ploskirev/go-rocket/order/internal/api/order/v1/converter"
	"github.com/ploskirev/go-rocket/order/internal/model"
	order_v1 "github.com/ploskirev/go-rocket/shared/pkg/openapi/order/v1"
)

func (a *Api) GetOrder(ctx context.Context, params order_v1.GetOrderParams) (order_v1.GetOrderRes, error) {
	order, err := a.os.GetOrder(ctx, params.OrderUUID)
	if err != nil {
		log.Printf("Error get order: %s", err)
		if errors.Is(err, model.ErrOrderNotFound) {
			return &order_v1.NotFoundError{
				Code:    404,
				Message: fmt.Sprintf("Order not found: %s", err),
			}, nil
		}

		return &order_v1.InternalServerError{
			Code:    500,
			Message: fmt.Sprintf("Error get order: %s", err),
		}, nil
	}

	return converter.OrderToApi(order), nil
}
