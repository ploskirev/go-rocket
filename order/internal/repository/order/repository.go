package orderrepo

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type orderRepo struct {
	pool *pgxpool.Pool
}

func NewOrderRepo(pool *pgxpool.Pool) *orderRepo {
	return &orderRepo{
		pool: pool,
	}
}
