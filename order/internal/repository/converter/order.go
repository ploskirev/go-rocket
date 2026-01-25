package converter

import (
	"github.com/ploskirev/go-rocket/order/internal/model"
	repoModel "github.com/ploskirev/go-rocket/order/internal/repository/model"
)

func OrderToRepoModel(o *model.Order) *repoModel.Order {
	var paymentMethod repoModel.PaymentMethod
	if o.PaymentMethod != nil {
		paymentMethod = repoModel.PaymentMethod(*o.PaymentMethod)
	}

	return &repoModel.Order{
		OrderUUID:       o.OrderUUID,
		UserUUID:        o.UserUUID,
		PartUUIDs:       o.PartUUIDs,
		TotalPrice:      o.TotalPrice,
		Status:          repoModel.OrderStatus(o.Status),
		TransactionUUID: o.TransactionUUID,
		PaymentMethod:   &paymentMethod,
	}
}

func OrderToModel(o *repoModel.Order) *model.Order {
	var paymentMethod model.PaymentMethod
	if o.PaymentMethod != nil {
		paymentMethod = model.PaymentMethod(*o.PaymentMethod)
	}

	return &model.Order{
		OrderUUID:       o.OrderUUID,
		UserUUID:        o.UserUUID,
		PartUUIDs:       o.PartUUIDs,
		TotalPrice:      o.TotalPrice,
		Status:          model.OrderStatus(o.Status),
		TransactionUUID: o.TransactionUUID,
		PaymentMethod:   &paymentMethod,
	}
}
