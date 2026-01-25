package payment

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/ploskirev/go-rocket/payment/internal/model"
)

func (ps *paymentService) PayOrder(ctx context.Context, _ model.PaymentInfo) (*uuid.UUID, error) {
	transactionUUID, err := uuid.NewV7()
	if err != nil {
		log.Printf("Ошибка генерации transaction uuid %s", err)
		return nil, fmt.Errorf("%s: %s", model.ErrGenerateUUID, err)
	}

	log.Printf("Оплата прошла успешно, transaction_uuid: %s", transactionUUID)

	return &transactionUUID, nil
}
