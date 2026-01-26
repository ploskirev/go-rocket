package orderservice

import (
	"errors"

	clientMocks "github.com/ploskirev/go-rocket/order/internal/client/grpc/mocks"
	"github.com/ploskirev/go-rocket/order/internal/model"
	repoMocks "github.com/ploskirev/go-rocket/order/internal/repository/mocks"
)

func (s *ServiceSuite) Test_GetOrder() {
	defaultOrderUUID := "123"
	defaultOrder := &model.Order{}

	tests := []struct {
		name       string
		uuid       string
		setupMocks func(*repoMocks.OrderRepository)
		checkFn    func(*model.Order, error)
	}{
		{
			name: "Get order success",
			uuid: defaultOrderUUID,
			setupMocks: func(mockRepo *repoMocks.OrderRepository) {
				mockRepo.On("GetOrder", s.ctx, defaultOrderUUID).Return(defaultOrder, nil)
			},
			checkFn: func(part *model.Order, err error) {
				s.NotNil(part)
				s.Nil(err)
			},
		},
		{
			name: "Get order failed",
			uuid: defaultOrderUUID,
			setupMocks: func(mockRepo *repoMocks.OrderRepository) {
				mockRepo.On("GetOrder", s.ctx, defaultOrderUUID).Return(nil, model.ErrOrderNotFound)
			},
			checkFn: func(part *model.Order, err error) {
				s.Nil(part)
				s.NotNil(err)
				s.True(errors.Is(err, model.ErrOrderNotFound))
			},
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			orderRepositoryMock := repoMocks.NewOrderRepository(s.T())
			clientInventoryMocks := clientMocks.NewInventoryClient(s.T())

			tt.setupMocks(orderRepositoryMock)

			service := NewOrderService(clientInventoryMocks, nil, orderRepositoryMock)
			order, err := service.GetOrder(s.ctx, tt.uuid)

			tt.checkFn(order, err)
		})
	}
}
