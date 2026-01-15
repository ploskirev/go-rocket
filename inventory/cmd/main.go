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

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"

	inventory_v1 "github.com/ploskirev/go-rocket/shared/pkg/proto/inventory/v1"
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
	s.mu.RLock()
	defer s.mu.RUnlock()
	st := *s.storage

	part, ok := st[req.Uuid]
	if !ok {
		log.Printf("Part with UUID %s not found\n", req.Uuid)
		return nil, status.Errorf(codes.NotFound, "part with UUID %s not found", req.Uuid)
	}

	return &inventory_v1.Part{
		Uuid:        part.UUID,
		Name:        part.Name,
		Price:       part.Price,
		Description: part.description,
	}, nil
}

func (s *inventoryService) ListPart(_ context.Context, req *inventory_v1.ListPartsRequest) (*inventory_v1.ListPartsResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	parts := make([]*inventory_v1.Part, 0, len(*s.storage))

	for _, p := range *s.storage {
		parts = append(parts, &inventory_v1.Part{
			Uuid:  p.UUID,
			Name:  p.Name,
			Price: p.Price,
		})
	}

	if req.Filter != nil {
		parts = filterParts(parts, req.Filter)
	}

	return &inventory_v1.ListPartsResponse{
		Parts: parts,
	}, nil
}

type FiltersMaps struct {
	byUUIDs                 map[string]struct{}
	byNames                 map[string]struct{}
	byCategories            map[inventory_v1.Category]struct{}
	byManufactoredCountries map[string]struct{}
	byTagsSlice             []string
}

func convertFiltersToMaps(filters *inventory_v1.PartsFilter) *FiltersMaps {
	filterUUIDsMap := make(map[string]struct{}, len(filters.Uuids))
	for _, f := range filters.Uuids {
		filterUUIDsMap[f] = struct{}{}
	}
	filterNamesMap := make(map[string]struct{}, len(filters.Names))
	for _, f := range filters.Names {
		filterNamesMap[f] = struct{}{}
	}
	filterCategoriesMap := make(map[inventory_v1.Category]struct{}, len(filters.Categories))
	for _, f := range filters.Categories {
		filterCategoriesMap[f] = struct{}{}
	}
	filterManufactoredCountriesMap := make(map[string]struct{}, len(filters.ManufacturerCountries))
	for _, f := range filters.ManufacturerCountries {
		filterManufactoredCountriesMap[f] = struct{}{}
	}

	return &FiltersMaps{
		byUUIDs:                 filterUUIDsMap,
		byNames:                 filterNamesMap,
		byCategories:            filterCategoriesMap,
		byManufactoredCountries: filterManufactoredCountriesMap,
		byTagsSlice:             filters.Tags,
	}
}

func filterParts(parts []*inventory_v1.Part, filters *inventory_v1.PartsFilter) []*inventory_v1.Part {
	filtersMaps := convertFiltersToMaps(filters)
	filteredParts := make([]*inventory_v1.Part, 0, len(parts))

	for _, part := range parts {
		if len(filtersMaps.byUUIDs) > 0 {
			if _, ok := filtersMaps.byUUIDs[part.Uuid]; !ok {
				continue
			}
		}
		if len(filtersMaps.byNames) > 0 {
			if _, ok := filtersMaps.byNames[part.Name]; !ok {
				continue
			}
		}
		if len(filtersMaps.byCategories) > 0 {
			if _, ok := filtersMaps.byCategories[part.Category]; !ok {
				continue
			}
		}
		if len(filtersMaps.byManufactoredCountries) > 0 {
			if _, ok := filtersMaps.byManufactoredCountries[part.Manufacturer.Country]; !ok {
				continue
			}
		}
		if len(filtersMaps.byTagsSlice) > 0 {
			skip := true
			for _, t := range part.Tags {
				if slices.Contains(filtersMaps.byTagsSlice, t) {
					skip = false
					break
				}
			}

			if skip {
				continue
			}
		}

		filteredParts = append(filteredParts, part)
	}

	return filteredParts
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
	store["z9x8c7v6b5"] = Part{
		UUID:  "z9x8c7v6b5",
		Name:  "Second detail",
		Price: 177.5,
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
