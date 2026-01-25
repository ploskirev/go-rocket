package paymentv1

import (
	"context"
	"log"

	"github.com/ploskirev/go-rocket/payment/internal/api/payment/v1/converter"
	payment_v1 "github.com/ploskirev/go-rocket/shared/pkg/proto/payment/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (a *api) PayOrder(ctx context.Context, req *payment_v1.PayOrderRequest) (*payment_v1.PayOrderResponse, error) {
	transactionUUID, err := a.ps.PayOrder(ctx, *converter.PaymentInfoToModel(req))
	if err != nil {
		log.Printf("Error while paing order: %s", err)
		return nil, status.Errorf(codes.Internal, "Error while paing order, %s", err)
	}

	return &payment_v1.PayOrderResponse{TransactionUuid: transactionUUID.String()}, nil
}
