package orderrepo

import (
	"sync"

	"github.com/ploskirev/go-rocket/order/internal/repository/model"
)

type orederRepo struct {
	mu      sync.RWMutex
	storage map[string]*model.Order
}

func NewOrderRepo() *orederRepo {
	return &orederRepo{
		storage: make(map[string]*model.Order),
	}
}
