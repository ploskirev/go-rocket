package paymentv1

import (
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	payment_v1 "github.com/ploskirev/go-rocket/shared/pkg/proto/payment/v1"
)

const (
	paymentAddress = "localhost:50052"
)

type paymentClient struct {
	pc payment_v1.PaymentServiceClient
}

func NewPaymentClient() (*paymentClient, *grpc.ClientConn, error) {
	paymentConn, err := grpc.NewClient(
		paymentAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Printf("failed to connect inventory service: %v\n", err)
		return nil, nil, err
	}
	pc := payment_v1.NewPaymentServiceClient(paymentConn)

	return &paymentClient{
		pc: pc,
	}, paymentConn, nil
}
