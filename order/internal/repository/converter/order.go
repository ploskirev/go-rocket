package converter

import (
	"strings"

	"github.com/ploskirev/go-rocket/order/internal/model"
	repoModel "github.com/ploskirev/go-rocket/order/internal/repository/model"
)

// func OrderToRepoModel(o *model.Order) *repoModel.Order {
// 	var paymentMethod repoModel.PaymentMethod
// 	if o.PaymentMethod != nil {
// 		paymentMethod = repoModel.PaymentMethod(*o.PaymentMethod)
// 	}

// 	return &repoModel.Order{
// 		OrderUUID:       sql.NullString{String: o.OrderUUID},
// 		UserUUID:        sql.NullString{String: o.UserUUID},
// 		PartUUIDs:       o.PartUUIDs,
// 		TotalPrice:      sql.NullFloat64{Float64: o.TotalPrice},
// 		Status:          sql.NullString{String: string(o.Status)},
// 		TransactionUUID: sql.NullString{String: *o.TransactionUUID},
// 		PaymentMethod:   sql.NullInt32{Int32: int32(paymentMethod)},
// 	}
// }

func OrderToModel(o *repoModel.Order) *model.Order {
	partUUIDs := []string{}
	if len(o.PartUUIDs.String) > 0 {
		partUUIDs = strings.Split(o.PartUUIDs.String, ",")
	}

	return &model.Order{
		OrderUUID:       o.OrderUUID.String,
		UserUUID:        o.UserUUID.String,
		PartUUIDs:       partUUIDs,
		TotalPrice:      o.TotalPrice.Float64,
		Status:          model.OrderStatus(o.Status.String),
		TransactionUUID: o.TransactionUUID.String,
		PaymentMethod:   model.PaymentMethod(o.PaymentMethod.Int32),
		CreatedAt:       o.CreatedAt.Time,
		UpdatedAt:       o.UpdatedAt.Time,
	}
}
