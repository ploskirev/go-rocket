package service

import (
	"context"

	"github.com/ploskirev/go-rocket/order/internal/model"
)

type OrderService interface {
	CancelOrder(ctx context.Context, orderUUID string) error
	CreateOrder(ctx context.Context, createOrderInfo *model.CreateOrderData) (*model.Order, error)
	GetOrder(ctx context.Context, orderUUID string) (*model.Order, error)
	PayOrder(ctx context.Context, paymentInfo *model.PaymentInfo) (string, error)
}
