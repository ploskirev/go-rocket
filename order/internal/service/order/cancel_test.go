package orderservice

import (
	"errors"

	"github.com/ploskirev/go-rocket/order/internal/model"
	repoMocks "github.com/ploskirev/go-rocket/order/internal/repository/mocks"
	"github.com/stretchr/testify/mock"
)

func (s *ServiceSuite) Test_CancelOrder() {
	defaultOrderUUID := "q1w2e3"
	defaultUserUUID := "789"

	getDefaultOrder := func() *model.Order {
		return &model.Order{
			UserUUID:  defaultUserUUID,
			OrderUUID: defaultOrderUUID,
			Status:    model.PENDING_PAYMENT,
		}
	}

	defaultMocksSetup := func(or *repoMocks.OrderRepository) {
		or.EXPECT().GetOrder(s.ctx, defaultOrderUUID).Return(getDefaultOrder(), nil)
		or.EXPECT().UpdateOrder(s.ctx, defaultOrderUUID, mock.Anything).Return(nil)
	}

	tests := []struct {
		name       string
		orderUUID  string
		setupMocks func(*repoMocks.OrderRepository)
		checkFn    func(error)
	}{
		{
			name:      "Cancel order success",
			orderUUID: defaultOrderUUID,
			checkFn: func(err error) {
				s.Nil(err)
			},
		},
		{
			name:      "Pay order fail if getting order return error",
			orderUUID: defaultOrderUUID,
			setupMocks: func(or *repoMocks.OrderRepository) {
				or.EXPECT().GetOrder(s.ctx, defaultOrderUUID).Unset()
				or.EXPECT().UpdateOrder(s.ctx, defaultOrderUUID, mock.Anything).Unset()
				or.EXPECT().GetOrder(s.ctx, defaultOrderUUID).Return(nil, model.ErrOrderNotFound)
			},
			checkFn: func(err error) {
				s.NotNil(err)
				s.True(errors.Is(err, model.ErrOrderNotFound))
			},
		},
		{
			name:      "Pay order fail if order has status PAID",
			orderUUID: defaultOrderUUID,
			setupMocks: func(or *repoMocks.OrderRepository) {
				or.EXPECT().GetOrder(s.ctx, defaultOrderUUID).Unset()
				or.EXPECT().UpdateOrder(s.ctx, defaultOrderUUID, mock.Anything).Unset()
				or.EXPECT().GetOrder(s.ctx, defaultOrderUUID).Return(&model.Order{Status: model.PAID}, nil)
			},
			checkFn: func(err error) {
				s.NotNil(err)
				s.True(errors.Is(err, model.ErrConflict))
			},
		},
		{
			name:      "Pay order fail if order has not status PENDING PAYMENT",
			orderUUID: defaultOrderUUID,
			setupMocks: func(or *repoMocks.OrderRepository) {
				or.EXPECT().GetOrder(s.ctx, defaultOrderUUID).Unset()
				or.EXPECT().UpdateOrder(s.ctx, defaultOrderUUID, mock.Anything).Unset()
				or.EXPECT().GetOrder(s.ctx, defaultOrderUUID).Return(&model.Order{Status: model.CANCELLED}, nil)
			},
			checkFn: func(err error) {
				s.NotNil(err)
				s.True(errors.Is(err, model.ErrBadRequest))
			},
		},
		{
			name:      "Pay order fail if updating order return error",
			orderUUID: defaultOrderUUID,
			setupMocks: func(or *repoMocks.OrderRepository) {
				or.EXPECT().UpdateOrder(s.ctx, defaultOrderUUID, mock.Anything).Unset()
				or.EXPECT().UpdateOrder(s.ctx, defaultOrderUUID, mock.Anything).Return(model.ErrOrderNotFound)
			},
			checkFn: func(err error) {
				s.NotNil(err)
				s.ErrorIs(err, model.ErrOrderNotFound)
			},
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			defer mock.AssertExpectationsForObjects(s.T())
			orderRepositoryMock := repoMocks.NewOrderRepository(s.T())

			defaultMocksSetup(orderRepositoryMock)
			if tt.setupMocks != nil {
				tt.setupMocks(orderRepositoryMock)
			}

			service := NewOrderService(nil, nil, orderRepositoryMock)
			err := service.CancelOrder(s.ctx, tt.orderUUID)

			tt.checkFn(err)
		})
	}
}
