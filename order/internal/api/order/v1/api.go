package orderapiv1

import (
	"context"
	"net/http"

	"github.com/ploskirev/go-rocket/order/internal/service"
	order_v1 "github.com/ploskirev/go-rocket/shared/pkg/openapi/order/v1"
)

type api struct {
	os service.OrderService
}

func NewApi(os service.OrderService) *api {
	return &api{
		os: os,
	}
}

func (a *api) NewError(ctx context.Context, err error) *order_v1.GenericErrorStatusCode {
	return &order_v1.GenericErrorStatusCode{
		StatusCode: http.StatusInternalServerError,
		Response: order_v1.GenericError{
			Code:    order_v1.NewOptInt(http.StatusInternalServerError),
			Message: order_v1.NewOptString(err.Error()),
		},
	}
}
