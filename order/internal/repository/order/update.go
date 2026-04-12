package orderrepo

import (
	"context"
	"log"
	"strings"

	sq "github.com/Masterminds/squirrel"
	"github.com/ploskirev/go-rocket/order/internal/model"
)

func (or *orderRepo) UpdateOrder(ctx context.Context, orderUUID string, orderInfo *model.Order) error {
	query, args, err := sq.Update("orders").
		PlaceholderFormat(sq.Dollar).
		Set("user_uuid", orderInfo.UserUUID).
		Set("part_uuids", strings.Join(orderInfo.PartUUIDs, ",")).
		Set("total_price", orderInfo.TotalPrice).
		Set("transaction_uuid", orderInfo.TransactionUUID).
		Set("payment_method", orderInfo.PaymentMethod).
		Set("order_status", orderInfo.Status).
		Set("updated_at", orderInfo.UpdatedAt).
		Where(sq.Eq{"order_uuid": orderUUID}).
		ToSql()
	if err != nil {
		log.Printf("failed to build query: %v\n", err)
		return err
	}

	_, err = or.pool.Exec(ctx, query, args...)
	if err != nil {
		log.Printf("failed to exec query: %v\n", err)
		return err
	}

	return nil
}
