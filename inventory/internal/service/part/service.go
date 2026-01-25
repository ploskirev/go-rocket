package part

import "github.com/ploskirev/go-rocket/inventory/internal/repository"

type partService struct {
	pr repository.PartRepository
}

func NewPartService(pr repository.PartRepository) *partService {
	return &partService{
		pr: pr,
	}
}
