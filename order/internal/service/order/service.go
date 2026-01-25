package orderservice

import (
	grpcClients "github.com/ploskirev/go-rocket/order/internal/client/grpc"
	"github.com/ploskirev/go-rocket/order/internal/repository"
)

type orederService struct {
	ic grpcClients.InventoryClient
	pc grpcClients.PaymentClient
	or repository.OrderRepository
}

func NewOrderService(ic grpcClients.InventoryClient, pc grpcClients.PaymentClient, or repository.OrderRepository) *orederService {
	return &orederService{
		ic: ic,
		pc: pc,
		or: or,
	}
}
