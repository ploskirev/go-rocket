package orderservice

import (
	"context"
	"log"

	"github.com/ploskirev/go-rocket/order/internal/model"
)

func (os *orederService) CancelOrder(ctx context.Context, orderUUID string) error {
	orderInfo, err := os.or.GetOrder(ctx, orderUUID)
	if err != nil {
		log.Printf("Error order not found")
		return err
	}
	if orderInfo.Status == model.PAID {
		log.Printf("Wrong status (conflict)")
		return model.ErrConflict
	}
	if orderInfo.Status != model.PENDING_PAYMENT {
		log.Printf("Wrong status")
		return model.ErrBadRequest
	}

	orderInfo.Status = model.CANCELLED

	if err = os.or.UpdateOrder(ctx, orderUUID, orderInfo); err != nil {
		log.Printf("Error update order")
		return err
	}

	return nil
}
