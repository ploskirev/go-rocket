package orderrepo

import (
	"context"
	"log"

	"github.com/ploskirev/go-rocket/order/internal/model"
	"github.com/ploskirev/go-rocket/order/internal/repository/converter"
)

func (or *orderRepo) GetOrder(ctx context.Context, orderUUID string) (*model.Order, error) {
	or.mu.RLock()
	defer or.mu.RUnlock()

	orderInfo, ok := or.storage[orderUUID]
	if !ok {
		log.Printf("Error order with uuid (%s) not found", orderUUID)
		return nil, model.ErrOrderNotFound
	}
	return converter.OrderToModel(orderInfo), nil
}
