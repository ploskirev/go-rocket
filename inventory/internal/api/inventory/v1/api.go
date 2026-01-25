package apiinventoryv1

import (
	"github.com/ploskirev/go-rocket/inventory/internal/service"
	inventory_v1 "github.com/ploskirev/go-rocket/shared/pkg/proto/inventory/v1"
)

type api struct {
	inventory_v1.UnimplementedInventoryServiceServer

	ps service.PartService
}

func NewApi(ps service.PartService) *api {
	return &api{
		ps: ps,
	}
}
