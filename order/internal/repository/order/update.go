package orderrepo

import (
	"context"
	"log"

	"github.com/ploskirev/go-rocket/order/internal/model"
	"github.com/ploskirev/go-rocket/order/internal/repository/converter"
)

func (or *orderRepo) UpdateOrder(ctx context.Context, orderUUID string, orderInfo *model.Order) error {
	or.mu.Lock()
	defer or.mu.Unlock()

	_, ok := or.storage[orderUUID]
	if !ok {
		log.Printf("Error order with uuid (%s) not found", orderUUID)
		return model.ErrOrderNotFound
	}

	or.storage[orderUUID] = converter.OrderToRepoModel(orderInfo)
	return nil
}
