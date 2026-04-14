package orderservice

import (
	"context"
	"log"

	"github.com/google/uuid"

	"github.com/ploskirev/go-rocket/order/internal/model"
)

func (os *orederService) CreateOrder(ctx context.Context, createOrderInfo *model.CreateOrderData) (*model.Order, error) {
	existedPartsList, err := os.ic.ListParts(ctx)
	if err != nil {
		log.Printf("Error get parts list: %s", err)
		return nil, err
	}

	totalPrice := float64(0)

	partsMap := map[string]*model.Part{}
	for _, p := range existedPartsList {
		partsMap[p.UUID] = p
	}

	for _, pu := range createOrderInfo.PartUUIDs {
		if partInfo, ok := partsMap[pu]; !ok {
			log.Printf("Error parts not found")
			return nil, model.ErrPartNotFound
		} else {
			totalPrice += partInfo.Price
		}
	}

	orderUUID, err := uuid.NewV7()
	if err != nil {
		log.Printf("Error create uuid: %s", err)
		return nil, model.ErrGenUUID
	}
	orderUUIDString := orderUUID.String()

	orderInfo, err := os.or.CreateOrder(ctx, &model.Order{
		OrderUUID:  orderUUIDString,
		UserUUID:   createOrderInfo.UserUUID,
		PartUUIDs:  createOrderInfo.PartUUIDs,
		TotalPrice: totalPrice,
		Status:     model.PENDING_PAYMENT,
	})
	if err != nil {
		log.Printf("Error create order: %s", err)
		return nil, model.ErrCreateOrder
	}

	return orderInfo, nil
}
