package orderrepo

import (
	"context"

	"github.com/ploskirev/go-rocket/order/internal/model"
	"github.com/ploskirev/go-rocket/order/internal/repository/converter"
)

func (or *orederRepo) CreateOrder(ctx context.Context, orderInfo *model.Order) *model.Order {
	or.mu.Lock()
	defer or.mu.Unlock()

	or.storage[orderInfo.OrderUUID] = converter.OrderToRepoModel(orderInfo)

	return orderInfo
}
