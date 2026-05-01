package paymentv1

import (
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/ploskirev/go-rocket/order/internal/config"
	payment_v1 "github.com/ploskirev/go-rocket/shared/pkg/proto/payment/v1"
)

type PaymentClient struct {
	pc payment_v1.PaymentServiceClient
}

func NewPaymentClient() (*PaymentClient, *grpc.ClientConn, error) {
	paymentAddress := config.AppConfig().Payment.Address()
	paymentConn, err := grpc.NewClient(
		paymentAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Printf("failed to connect inventory service: %v\n", err)
		return nil, nil, err
	}
	pc := payment_v1.NewPaymentServiceClient(paymentConn)

	return &PaymentClient{
		pc: pc,
	}, paymentConn, nil
}
