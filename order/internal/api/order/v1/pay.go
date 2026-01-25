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

func (a *api) PayOrder(ctx context.Context, req *order_v1.PayOrderRequest, params order_v1.PayOrderParams) (order_v1.PayOrderRes, error) {
	transactionUUID, err := a.os.PayOrder(ctx, &model.PaymentInfo{
		OrderUUID:     params.OrderUUID,
		PaymentMethod: converter.PaymentMthodToModel(req.PaymentMethod),
	})
	if err != nil {
		log.Printf("Error pay order: %s", err)

		if errors.Is(err, model.ErrOrderNotFound) {
			return &order_v1.NotFoundError{
				Code:    404,
				Message: fmt.Sprintf("Error pay order (order not found). %s", err),
			}, nil
		}

		return &order_v1.InternalServerError{
			Code:    500,
			Message: fmt.Sprintf("Error pay order: %s", err),
		}, nil
	}

	return &order_v1.PayOrderResponse{TransactionUUID: transactionUUID}, nil
}
