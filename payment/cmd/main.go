package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"

	payment_v1 "github.com/ploskirev/go-rocket/shared/pkg/proto/payment/v1"
)

const grpcPort = 50052

type paymentService struct {
	payment_v1.UnimplementedPaymentServiceServer
}

func (s *paymentService) PayOrder(_ context.Context, req *payment_v1.PayOrderRequest) (*payment_v1.PayOrderResponse, error) {
	// fmt.Println("PAY ORDER!")

	transactionUUID, err := uuid.NewV6()
	if err != nil {
		log.Printf("Ошибка генерации transaction uuid %s", err)
		return nil, status.Errorf(codes.NotFound, "Ошибка генерации transaction uuid %s", err)
	}

	log.Printf("Оплата прошла успешно, transaction_uuid: %s", transactionUUID)

	return &payment_v1.PayOrderResponse{
		TransactionUuid: transactionUUID.String(),
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Printf("failed to listen: %v\n", err)
	}
	defer func() {
		if cerr := lis.Close(); cerr != nil {
			log.Printf("failed to close listener: %v\n", cerr)
		}
	}()

	s := grpc.NewServer()

	service := &paymentService{}

	payment_v1.RegisterPaymentServiceServer(s, service)

	reflection.Register(s)

	go func() {
		log.Printf("🚀 gRPC server listening on %d\n", grpcPort)
		err = s.Serve(lis)
		if err != nil {
			log.Printf("failed to serve: %v\n", err)
			return
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("🛑 Shutting down gRPC server...")
	s.GracefulStop()
	log.Println("✅ Server stopped")
}
