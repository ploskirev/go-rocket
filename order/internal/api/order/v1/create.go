package orderapiv1

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/ploskirev/go-rocket/order/internal/model"
	order_v1 "github.com/ploskirev/go-rocket/shared/pkg/openapi/order/v1"
)

func (a *api) CreateOrder(ctx context.Context, req *order_v1.CreateOrderRequest) (order_v1.CreateOrderRes, error) {
	orderInfo, err := a.os.CreateOrder(ctx, &model.CreateOrderData{
		UserUUID:  req.UserUUID,
		PartUUIDs: req.PartUuids,
	})
	if err != nil {
		log.Printf("Error create order: %s", err)
		if errors.Is(err, model.ErrPartNotFound) {
			return &order_v1.NotFoundError{
				Code:    404,
				Message: fmt.Sprintf("Error create order: %s", err),
			}, nil
		}

		return &order_v1.InternalServerError{
			Code:    500,
			Message: fmt.Sprintf("Error create order: %s", err),
		}, nil
	}

	return &order_v1.OrderDto{
		OrderUUID:  orderInfo.OrderUUID,
		TotalPrice: float32(orderInfo.TotalPrice),
	}, nil
}
