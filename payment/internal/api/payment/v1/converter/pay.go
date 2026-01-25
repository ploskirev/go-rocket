package converter

import (
	"github.com/ploskirev/go-rocket/payment/internal/model"
	payment_v1 "github.com/ploskirev/go-rocket/shared/pkg/proto/payment/v1"
)

func PaymentInfoToModel(o *payment_v1.PayOrderRequest) *model.PaymentInfo {
	return &model.PaymentInfo{
		OrderUUID:     o.OrderUuid,
		UserUUID:      o.UserUuid,
		PaymentMethod: model.PaymentMethod(o.PaymentMethod),
	}
}
