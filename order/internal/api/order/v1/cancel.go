package orderapiv1

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/ploskirev/go-rocket/order/internal/model"
	order_v1 "github.com/ploskirev/go-rocket/shared/pkg/openapi/order/v1"
)

func (a *Api) CancelOrder(ctx context.Context, params order_v1.CancelOrderParams) (order_v1.CancelOrderRes, error) {
	if err := a.os.CancelOrder(ctx, params.OrderUUID); err != nil {
		log.Printf("Error cancel order")

		if errors.Is(err, model.ErrOrderNotFound) {
			return &order_v1.NotFoundError{
				Code:    404,
				Message: fmt.Sprintf("Error cancel order (order not found). %s", err),
			}, nil
		}
		if errors.Is(err, model.ErrConflict) {
			return &order_v1.ConflictError{
				Code:    409,
				Message: fmt.Sprintf("Error cancel order (conflict order status). %s", err),
			}, nil
		}
		if errors.Is(err, model.ErrConflict) {
			return &order_v1.BadRequestError{
				Code:    400,
				Message: fmt.Sprintf("Error cancel order (wrong order status). %s", err),
			}, nil
		}

		return &order_v1.InternalServerError{
			Code:    500,
			Message: fmt.Sprintf("Error cancel order: %s", err),
		}, nil
	}

	return &order_v1.CancelOrderNoContent{}, nil
}
