package converter

import (
	"github.com/ploskirev/go-rocket/order/internal/model"
	order_v1 "github.com/ploskirev/go-rocket/shared/pkg/openapi/order/v1"
)

func OrderToApi(o *model.Order) *order_v1.GetOrderResponse {
	var transactionUUID string
	if o.TransactionUUID != nil {
		transactionUUID = *o.TransactionUUID
	}
	return &order_v1.GetOrderResponse{
		TransactionUUID: transactionUUID,
		OrderUUID:       o.OrderUUID,
		UserUUID:        o.UserUUID,
		PartUuids:       o.PartUUIDs,
		TotalPrice:      float32(o.TotalPrice),
		Status:          order_v1.OrderStatus(o.Status),
		PaymentMethod:   order_v1.PaymentMethod(*o.PaymentMethod),
	}
}

var paymentMethodMap = map[order_v1.PaymentMethod]model.PaymentMethod{
	order_v1.PaymentMethodPAYMENTMETHODUNKNOWN:       model.PaymentMethod_UNKNOWN,
	order_v1.PaymentMethodPAYMENTMETHODCARD:          model.PaymentMethod_CARD,
	order_v1.PaymentMethodPAYMENTMETHODSBP:           model.PaymentMethod_SBP,
	order_v1.PaymentMethodPAYMENTMETHODCREDITCARD:    model.PaymentMethod_CREDIT_CARD,
	order_v1.PaymentMethodPAYMENTMETHODINVESTORMONEY: model.PaymentMethod_INVESTOR_MONEY,
}

func PaymentMthodToModel(pm order_v1.PaymentMethod) model.PaymentMethod {
	return paymentMethodMap[pm]
}
