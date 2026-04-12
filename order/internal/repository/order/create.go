package orderrepo

import (
	"context"
	"log"
	"strings"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/ploskirev/go-rocket/order/internal/model"
)

func (or *orderRepo) CreateOrder(ctx context.Context, orderInfo *model.Order) (*model.Order, error) {
	query, args, err := sq.Insert("orders").
		PlaceholderFormat(sq.Dollar).
		Columns("order_uuid", "user_uuid", "part_uuids", "total_price", "order_status", "created_at").
		Values(orderInfo.OrderUUID, orderInfo.UserUUID, strings.Join(orderInfo.PartUUIDs, ","), orderInfo.TotalPrice, orderInfo.Status, time.Now()).
		ToSql()
	if err != nil {
		log.Printf("failed to build query: %v\n", err)
		return nil, err
	}

	_, err = or.pool.Exec(ctx, query, args...)
	if err != nil {
		log.Printf("failed to exec query: %v\n", err)
		return nil, err
	}

	return orderInfo, nil
}
