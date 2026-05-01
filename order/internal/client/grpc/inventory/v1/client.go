package inventoryv1

import (
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/ploskirev/go-rocket/order/internal/config"
	inventory_v1 "github.com/ploskirev/go-rocket/shared/pkg/proto/inventory/v1"
)

type InventoryClient struct {
	ic inventory_v1.InventoryServiceClient
}

func NewInventoryClient() (*InventoryClient, *grpc.ClientConn, error) {
	inventoryAddress := config.AppConfig().Inventory.Address()
	inventoryConn, err := grpc.NewClient(
		inventoryAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Printf("failed to connect inventory service: %v\n", err)
		return nil, nil, err
	}
	ic := inventory_v1.NewInventoryServiceClient(inventoryConn)

	return &InventoryClient{
		ic: ic,
	}, inventoryConn, nil
}
