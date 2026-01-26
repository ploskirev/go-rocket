package orderservice

import (
	"errors"

	"github.com/stretchr/testify/mock"

	clientMocks "github.com/ploskirev/go-rocket/order/internal/client/grpc/mocks"
	"github.com/ploskirev/go-rocket/order/internal/model"
	repoMocks "github.com/ploskirev/go-rocket/order/internal/repository/mocks"
)

func (s *ServiceSuite) Test_CreateOrder() {
	defaultOrderUUID := "q1w2e3"
	defaultUserUUID := "789"
	defaultPartsUUIDs := []string{"123", "456"}

	defaultParts := []*model.Part{
		{UUID: defaultPartsUUIDs[1], Name: "First part", Price: 123.5},
		{UUID: defaultPartsUUIDs[0], Name: "Second part", Price: 15.7},
	}

	defaultCreateOrderData := &model.CreateOrderData{
		UserUUID:  defaultUserUUID,
		PartUUIDs: defaultPartsUUIDs,
	}

	defaultOrder := &model.Order{
		UserUUID:  defaultUserUUID,
		PartUUIDs: defaultPartsUUIDs,
		OrderUUID: defaultOrderUUID,
	}

	defaultMocksSetup := func(ic *clientMocks.InventoryClient, or *repoMocks.OrderRepository) {
		ic.EXPECT().ListParts(s.ctx).Return(defaultParts, nil)
		or.EXPECT().CreateOrder(s.ctx, mock.Anything).Return(defaultOrder)
	}

	tests := []struct {
		name            string
		createOrderData *model.CreateOrderData
		setupMocks      func(*clientMocks.InventoryClient, *repoMocks.OrderRepository)
		checkFn         func(*model.Order, error)
	}{
		{
			name:            "Create order success",
			createOrderData: defaultCreateOrderData,
			checkFn: func(order *model.Order, err error) {
				s.Nil(err)
				s.NotNil(order)
				s.Equal(defaultOrderUUID, order.OrderUUID)
				s.Equal(defaultUserUUID, order.UserUUID)
			},
		},
		{
			name:            "Create order failed if inventory client returned error",
			createOrderData: defaultCreateOrderData,
			setupMocks: func(ic *clientMocks.InventoryClient, or *repoMocks.OrderRepository) {
				ic.EXPECT().ListParts(s.ctx).Unset()
				or.EXPECT().CreateOrder(s.ctx, mock.Anything).Unset()
				ic.EXPECT().ListParts(s.ctx).Return(nil, model.ErrGetListParts)
			},
			checkFn: func(order *model.Order, err error) {
				s.NotNil(err)
				s.True(errors.Is(err, model.ErrGetListParts))
				s.Nil(order)
			},
		},
		{
			name: "Create order failed if inventory client returned error",
			createOrderData: &model.CreateOrderData{
				PartUUIDs: []string{"wrong part"},
			},
			setupMocks: func(ic *clientMocks.InventoryClient, or *repoMocks.OrderRepository) {
				or.EXPECT().CreateOrder(s.ctx, mock.Anything).Unset()
			},
			checkFn: func(order *model.Order, err error) {
				s.NotNil(err)
				s.True(errors.Is(err, model.ErrPartNotFound))
				s.Nil(order)
			},
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			orderRepositoryMock := repoMocks.NewOrderRepository(s.T())
			clientInventoryMocks := clientMocks.NewInventoryClient(s.T())

			defaultMocksSetup(clientInventoryMocks, orderRepositoryMock)
			if tt.setupMocks != nil {
				tt.setupMocks(clientInventoryMocks, orderRepositoryMock)
			}

			service := NewOrderService(clientInventoryMocks, nil, orderRepositoryMock)
			order, err := service.CreateOrder(s.ctx, tt.createOrderData)

			tt.checkFn(order, err)
		})
	}
}
