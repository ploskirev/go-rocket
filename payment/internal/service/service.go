package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/ploskirev/go-rocket/payment/internal/model"
)

type PaymentService interface {
	PayOrder(ctx context.Context, _ model.PaymentInfo) (*uuid.UUID, error)
}
