package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"slices"
	"sync"
	"syscall"
	"time"

	inventory_v1 "github.com/ploskirev/go-rocket/shared/pkg/proto/inventory/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

const grpcPort = 50051

type Category int32

const (
	UNKNOWN Category = iota
	ENGINE
	FUEL
	PORTHOLE
	WING
)

type Part struct {
	UUID          string
	Name          string
	description   string
	Price         float64
	StockQuantity int64
	Category      Category
	Dimensions    Dimensions
	Manufacturer  Manifacturer
	Tags          []string
	Metadata      map[string]string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type Dimensions struct {
	Length float64
	Width  float64
	Height float64
	Weight float64
}

type Manifacturer struct {
	Name    string
	Country string
	Website string
}

type Storage map[string]Part

type inventoryService struct {
	inventory_v1.UnimplementedInventoryServiceServer
	mu sync.RWMutex

	storage *Storage
}

func NewInventoryService() *inventoryService {
	return &inventoryService{
		storage: &Storage{},
	}
}

func (s *inventoryService) GetPart(_ context.Context, req *inventory_v1.GetPartRequest) (*inventory_v1.Part, error) {
	fmt.Println("UUID: ", req.Uuid)
	s.mu.RLock()
	defer s.mu.RUnlock()
	st := *s.storage

	part, ok := st[req.Uuid]
	if !ok {
		log.Printf("Part with UUID %s not found\n", req.Uuid)
		return nil, status.Errorf(codes.NotFound, "part with UUID %s not found", req.Uuid)
	}

	return &inventory_v1.Part{
		Uuid:  part.UUID,
		Name:  part.Name,
		Price: part.Price,
	}, nil
}

func (s *inventoryService) ListPart(_ context.Context, req *inventory_v1.ListPartsRequest) (*inventory_v1.ListPartsResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	parts := make([]*inventory_v1.Part, len(*s.storage))

	for _, p := range *s.storage {
		parts = append(parts, &inventory_v1.Part{
			Uuid:  p.UUID,
			Name:  p.Name,
			Price: p.Price,
		})
	}

	filteredParts := make([]*inventory_v1.Part, len(parts))
	for _, part := range parts {
		skip := false

		if len(req.Filter.Uuids) > 0 {
			skip = true

			// by UUID
			for _, f := range req.Filter.Uuids {
				if f == part.Uuid {
					skip = false
					break
				}
			}
			if skip == true {
				continue
			}

			// by names
			for _, f := range req.Filter.Names {
				if f == part.Name {
					skip = false
					break
				}
			}
			if skip == true {
				continue
			}

			// by category
			for _, f := range req.Filter.Categories {
				if f == part.Category {
					skip = false
					break
				}
			}
			if skip == true {
				continue
			}

			// by manufacturer country
			for _, f := range req.Filter.ManufacturerCountries {
				if f == part.Manufacturer.Country {
					skip = false
					break
				}
			}
			if skip == true {
				continue
			}

			// by manufacturer country
			for _, f := range req.Filter.Tags {
				if slices.Contains(part.Tags, f) {
					skip = false
					break
				}
			}
			if skip == true {
				continue
			}

			filteredParts = append(filteredParts, part)
		}
	}

	return &inventory_v1.ListPartsResponse{
		Parts: parts,
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

	service := NewInventoryService()

	store := *service.storage
	store["q1w2e3r4t5"] = Part{
		UUID:  "q1w2e3r4t5",
		Name:  "Test detail",
		Price: 535.7,
	}

	inventory_v1.RegisterInventoryServiceServer(s, service)

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
