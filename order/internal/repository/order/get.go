package orderrepo

import (
	"context"
	"log"

	sq "github.com/Masterminds/squirrel"

	"github.com/ploskirev/go-rocket/order/internal/model"
	"github.com/ploskirev/go-rocket/order/internal/repository/converter"
	repoModel "github.com/ploskirev/go-rocket/order/internal/repository/model"
)

func (or *orderRepo) GetOrder(ctx context.Context, orderUUID string) (*model.Order, error) {
	query, args, err := sq.Select("order_uuid", "user_uuid", "part_uuids", "total_price", "transaction_uuid", "payment_method", "order_status", "created_at", "updated_at").
		From("orders").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"order_uuid": orderUUID}).
		ToSql()
	if err != nil {
		log.Printf("failed to build query: %v\n", err)
		return nil, err
	}

	order := repoModel.Order{}

	err = or.pool.QueryRow(ctx, query, args...).Scan(&order.OrderUUID, &order.UserUUID, &order.PartUUIDs, &order.TotalPrice, &order.TransactionUUID, &order.PaymentMethod, &order.Status, &order.CreatedAt, &order.UpdatedAt)
	if err != nil {
		log.Printf("failed to select notes: %v\n", err)
		return nil, err
	}

	return converter.OrderToModel(&order), nil
}
