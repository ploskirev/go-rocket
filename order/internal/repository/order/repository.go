package orderrepo

import (
	"sync"

	"github.com/ploskirev/go-rocket/order/internal/repository/model"
)

type orderRepo struct {
	mu      sync.RWMutex
	storage map[string]*model.Order
}

func NewOrderRepo() *orderRepo {
	return &orderRepo{
		storage: make(map[string]*model.Order),
	}
}
