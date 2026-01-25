package converter

import (
	"github.com/ploskirev/go-rocket/order/internal/model"
	inventory_v1 "github.com/ploskirev/go-rocket/shared/pkg/proto/inventory/v1"
	payment_v1 "github.com/ploskirev/go-rocket/shared/pkg/proto/payment/v1"
)

func PartToModel(p *inventory_v1.Part) *model.Part {
	return &model.Part{
		UUID:  p.Uuid,
		Name:  p.Name,
		Price: p.Price,
	}
}

func PaymentInfoToProto(info *model.PaymentInfo) *payment_v1.PayOrderRequest {
	return &payment_v1.PayOrderRequest{
		OrderUuid:     info.OrderUUID,
		UserUuid:      info.UserUUID,
		PaymentMethod: payment_v1.PaymentMethod(info.PaymentMethod),
	}
}
