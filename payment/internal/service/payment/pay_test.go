package payment

import (
	"github.com/google/uuid"
	"github.com/ploskirev/go-rocket/payment/internal/model"
)

func (s *ServiceSuite) TestPaySuccess() {
	transactionUUID, err := s.service.PayOrder(s.ctx, model.PaymentInfo{})
	s.NotNil(transactionUUID)
	s.Nil(err)
}

func (s *ServiceSuite) TestPay() {
	tests := []struct {
		name    string
		checkFn func(transactionUUID *uuid.UUID, err error)
	}{
		{
			name: "Pay success",
			checkFn: func(transactionUUID *uuid.UUID, err error) {
				s.NotNil(transactionUUID)
				s.Nil(err)
			},
		},
	}

	for _, t := range tests {
		s.Run(t.name, func() {
			transactionUUID, err := s.service.PayOrder(s.ctx, model.PaymentInfo{})

			t.checkFn(transactionUUID, err)
		})
	}
}
