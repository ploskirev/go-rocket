package orderservice

import (
	"errors"

	"github.com/stretchr/testify/mock"

	clientMocks "github.com/ploskirev/go-rocket/order/internal/client/grpc/mocks"
	"github.com/ploskirev/go-rocket/order/internal/model"
	repoMocks "github.com/ploskirev/go-rocket/order/internal/repository/mocks"
)

func (s *ServiceSuite) Test_PayOrder() {
	defaultOrderUUID := "q1w2e3"
	defaultUserUUID := "789"

	defaultPaymentInfo := &model.PaymentInfo{
		OrderUUID:     defaultOrderUUID,
		UserUUID:      defaultUserUUID,
		PaymentMethod: model.PaymentMethod_CARD,
	}

	defaultTransactionUUID := "y6u7i8"

	defaultOrder := &model.Order{
		UserUUID:  defaultUserUUID,
		OrderUUID: defaultOrderUUID,
	}

	defaultMocksSetup := func(pc *clientMocks.PaymentClient, or *repoMocks.OrderRepository) {
		or.EXPECT().GetOrder(s.ctx, defaultOrderUUID).Return(defaultOrder, nil)
		pc.EXPECT().PayOrder(s.ctx, defaultPaymentInfo).Return(defaultTransactionUUID, nil)
		or.EXPECT().UpdateOrder(s.ctx, defaultOrderUUID, mock.Anything).Return(nil)
	}

	tests := []struct {
		name          string
		paymenentInfo *model.PaymentInfo
		setupMocks    func(*clientMocks.PaymentClient, *repoMocks.OrderRepository)
		checkFn       func(string, error)
	}{
		{
			name:          "Pay order success",
			paymenentInfo: defaultPaymentInfo,
			checkFn: func(transactionUUID string, err error) {
				s.Nil(err)
				s.NotNil(transactionUUID)
				s.Equal(defaultTransactionUUID, transactionUUID)
			},
		},
		{
			name:          "Pay order fail if getting order return error",
			paymenentInfo: defaultPaymentInfo,
			setupMocks: func(pc *clientMocks.PaymentClient, or *repoMocks.OrderRepository) {
				or.EXPECT().GetOrder(s.ctx, defaultOrderUUID).Unset()
				pc.EXPECT().PayOrder(s.ctx, defaultPaymentInfo).Unset()
				or.EXPECT().UpdateOrder(s.ctx, defaultOrderUUID, mock.Anything).Unset()
				or.EXPECT().GetOrder(s.ctx, defaultOrderUUID).Return(nil, model.ErrOrderNotFound)
			},
			checkFn: func(transactionUUID string, err error) {
				s.Equal("", transactionUUID)
				s.NotNil(err)
				s.True(errors.Is(err, model.ErrOrderNotFound))
			},
		},
		{
			name:          "Pay order fail if paying order return error",
			paymenentInfo: defaultPaymentInfo,
			setupMocks: func(pc *clientMocks.PaymentClient, or *repoMocks.OrderRepository) {
				pc.EXPECT().PayOrder(s.ctx, defaultPaymentInfo).Unset()
				pc.EXPECT().PayOrder(s.ctx, defaultPaymentInfo).Return("", model.ErrPayOrder)
				or.EXPECT().UpdateOrder(s.ctx, defaultOrderUUID, mock.Anything).Unset()
			},
			checkFn: func(transactionUUID string, err error) {
				s.Equal("", transactionUUID)
				s.NotNil(err)
				s.True(errors.Is(err, model.ErrPayOrder))
			},
		},
		{
			name:          "Pay order fail if paying order return error",
			paymenentInfo: defaultPaymentInfo,
			setupMocks: func(pc *clientMocks.PaymentClient, or *repoMocks.OrderRepository) {
				or.EXPECT().UpdateOrder(s.ctx, defaultOrderUUID, mock.Anything).Unset()
				or.EXPECT().UpdateOrder(s.ctx, defaultOrderUUID, mock.Anything).Return(model.ErrOrderNotFound)
			},
			checkFn: func(transactionUUID string, err error) {
				s.Equal("", transactionUUID)
				s.NotNil(err)
				s.True(errors.Is(err, model.ErrOrderNotFound))
			},
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			orderRepositoryMock := repoMocks.NewOrderRepository(s.T())
			clientPaymentMocks := clientMocks.NewPaymentClient(s.T())

			defaultMocksSetup(clientPaymentMocks, orderRepositoryMock)
			if tt.setupMocks != nil {
				tt.setupMocks(clientPaymentMocks, orderRepositoryMock)
			}

			service := NewOrderService(nil, clientPaymentMocks, orderRepositoryMock)
			transactionUUID, err := service.PayOrder(s.ctx, tt.paymenentInfo)

			tt.checkFn(transactionUUID, err)
		})
	}
}
