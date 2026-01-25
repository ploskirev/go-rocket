package inventoryv1

import (
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	inventory_v1 "github.com/ploskirev/go-rocket/shared/pkg/proto/inventory/v1"
)

const (
	inventoryAddress = "localhost:50051"
)

type inventoryClient struct {
	ic inventory_v1.InventoryServiceClient
}

func NewInventoryClient() (*inventoryClient, *grpc.ClientConn, error) {
	inventoryConn, err := grpc.NewClient(
		inventoryAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Printf("failed to connect inventory service: %v\n", err)
		return nil, nil, err
	}
	ic := inventory_v1.NewInventoryServiceClient(inventoryConn)

	return &inventoryClient{
		ic: ic,
	}, inventoryConn, nil
}
