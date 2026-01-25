package orderservice

import (
	"context"
	"log"

	"github.com/ploskirev/go-rocket/order/internal/model"
)

func (os *orederService) GetOrder(ctx context.Context, orderUUID string) (*model.Order, error) {
	order, err := os.or.GetOrder(ctx, orderUUID)
	if err != nil {
		log.Printf("Error get order: %s", err)
		return nil, err
	}

	return order, nil
}
