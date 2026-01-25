package part

import (
	"github.com/ploskirev/go-rocket/inventory/internal/repository/model"
)

func (pr *partRepository) InitStorage() {
	pr.storage["q1w2e3r4t5"] = &model.Part{
		UUID:  "q1w2e3r4t5",
		Name:  "Test detail",
		Price: 535.7,
	}
	pr.storage["z9x8c7v6b5"] = &model.Part{
		UUID:  "z9x8c7v6b5",
		Name:  "Second detail",
		Price: 177.5,
	}
}
