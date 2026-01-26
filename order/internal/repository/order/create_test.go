package orderrepo

import (
	"context"

	"github.com/ploskirev/go-rocket/order/internal/model"
)

func (s *ServiceSuite) Test_CreateOrder() {
	tests := []struct {
		name           string
		orderInfo      *model.Order
		expectedResult *model.Order
	}{
		{
			name: "Should execute successfully",
			orderInfo: &model.Order{
				OrderUUID: "123",
			},
			expectedResult: &model.Order{
				OrderUUID: "123",
			},
		},
	}
	for _, tt := range tests {
		s.Run(tt.name, func() {
			or := NewOrderRepo()
			result := or.CreateOrder(context.Background(), tt.orderInfo)

			s.Equal(tt.expectedResult, result)
		})
	}
}
