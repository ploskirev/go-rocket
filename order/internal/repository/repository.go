package repository

import (
	"context"

	"github.com/ploskirev/go-rocket/order/internal/model"
)

type OrderRepository interface {
	CreateOrder(ctx context.Context, orderInfo *model.Order) *model.Order
	GetOrder(ctx context.Context, orderUUID string) (*model.Order, error)
	UpdateOrder(ctx context.Context, orderUUID string, orderInfo *model.Order) error
}
