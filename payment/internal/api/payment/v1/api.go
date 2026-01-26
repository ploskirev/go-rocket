package paymentv1

import (
	"github.com/ploskirev/go-rocket/payment/internal/service"
	payment_v1 "github.com/ploskirev/go-rocket/shared/pkg/proto/payment/v1"
)

type api struct {
	payment_v1.UnimplementedPaymentServiceServer

	ps service.PaymentService
}

func NewApi(ps service.PaymentService) *api {
	return &api{
		ps: ps,
	}
}
