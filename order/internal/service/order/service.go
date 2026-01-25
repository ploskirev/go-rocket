package orderservice

import (
	grpcClients "github.com/ploskirev/go-rocket/order/internal/client/grpc"
	"github.com/ploskirev/go-rocket/order/internal/repository"
)

type orederService struct {
	ic grpcClients.InventoryClient
	pc grpcClients.PaymentClient
	or repository.OrderRepo
}

func NewOrderService(ic grpcClients.InventoryClient, pc grpcClients.PaymentClient, or repository.OrderRepo) *orederService {
	return &orederService{
		ic: ic,
		pc: pc,
		or: or,
	}
}
