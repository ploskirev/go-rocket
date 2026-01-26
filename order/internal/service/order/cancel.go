package orderservice

import (
	"context"
	"fmt"
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
		return fmt.Errorf("%w: Order has status %s", model.ErrConflict, orderInfo.Status)
	}
	if orderInfo.Status != model.PENDING_PAYMENT {
		log.Printf("Wrong status")
		return fmt.Errorf("%w: Order has status %s", model.ErrBadRequest, orderInfo.Status)
	}

	orderInfo.Status = model.CANCELLED

	if err = os.or.UpdateOrder(ctx, orderUUID, orderInfo); err != nil {
		log.Printf("Error update order")
		return err
	}

	return nil
}
