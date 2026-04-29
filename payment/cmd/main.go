package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	paymentApi "github.com/ploskirev/go-rocket/payment/internal/api/payment/v1"
	"github.com/ploskirev/go-rocket/payment/internal/config"
	paymentService "github.com/ploskirev/go-rocket/payment/internal/service/payment"
	payment_v1 "github.com/ploskirev/go-rocket/shared/pkg/proto/payment/v1"
)

const configPath = "./deploy/compose/payment/.env"

func main() {
	err := config.Load(configPath)
	if err != nil {
		panic(fmt.Errorf("failed to load config: %w", err))
	}

	lis, err := net.Listen("tcp", config.AppConfig().PaymentGRPC.Address())
	if err != nil {
		log.Printf("failed to listen: %v\n", err)
	}
	defer func() {
		if cerr := lis.Close(); cerr != nil {
			log.Printf("failed to close listener: %v\n", cerr)
		}
	}()

	s := grpc.NewServer()

	payService := paymentService.NewPaymentService()
	payApi := paymentApi.NewApi(payService)

	payment_v1.RegisterPaymentServiceServer(s, payApi)

	reflection.Register(s)

	go func() {
		log.Printf("🚀 gRPC server listening on %d\n", config.AppConfig().PaymentGRPC.Address())
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
