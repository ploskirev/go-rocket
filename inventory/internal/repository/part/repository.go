package part

import (
	"sync"

	repoModel "github.com/ploskirev/go-rocket/inventory/internal/repository/model"
)

type partRepository struct {
	mu      sync.RWMutex
	storage map[string]*repoModel.Part
}

func NewPartRepository() *partRepository {
	return &partRepository{
		storage: map[string]*repoModel.Part{},
	}
}
