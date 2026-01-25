package grpc

import (
	"context"

	"github.com/ploskirev/go-rocket/order/internal/model"
)

type InventoryClient interface {
	ListParts(ctx context.Context) ([]*model.Part, error)
}

type PaymentClient interface {
	PayOrder(ctx context.Context, paymentInfo *model.PaymentInfo) (string, error)
}
