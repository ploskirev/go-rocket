package part

import (
	"github.com/ploskirev/go-rocket/inventory/internal/model"
	"github.com/ploskirev/go-rocket/inventory/internal/repository/mocks"
)

func (s *ServiceSuite) Test_ListPartsSuccess() {
	partRepositoryMock := mocks.NewPartRepository(s.T())

	filters := &model.Filters{}

	partRepositoryMock.On("ListParts", s.ctx, filters).Return([]*model.Part{{UUID: "123"}, {UUID: "789"}}, nil)

	service := NewPartService(partRepositoryMock)
	parts, _ := service.ListParts(s.ctx, filters)

	s.NotNil(parts)
	s.Equal(2, len(parts))
}
