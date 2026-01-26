package part

import (
	"errors"

	"github.com/ploskirev/go-rocket/inventory/internal/model"
	"github.com/ploskirev/go-rocket/inventory/internal/repository/mocks"
)

func (s *ServiceSuite) TestGetPart() {
	defaultPartUUID := "123"
	defaultPart := &model.Part{UUID: "789"}

	tests := []struct {
		name       string
		uuid       string
		setupMocks func(*mocks.PartRepository)
		checkFn    func(*model.Part, error)
	}{
		{
			name: "Get order success",
			uuid: defaultPartUUID,
			setupMocks: func(mockRepo *mocks.PartRepository) {
				mockRepo.On("GetPart", s.ctx, defaultPartUUID).Return(defaultPart, nil)
			},
			checkFn: func(part *model.Part, err error) {
				s.NotNil(part)
				s.Nil(err)
			},
		},
		{
			name: "Get order failed",
			uuid: defaultPartUUID,
			setupMocks: func(mockRepo *mocks.PartRepository) {
				mockRepo.On("GetPart", s.ctx, defaultPartUUID).Return(nil, model.ErrPartNotFound)
			},
			checkFn: func(part *model.Part, err error) {
				s.Nil(part)
				s.NotNil(err)
				s.True(errors.Is(err, model.ErrPartNotFound))
			},
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			partRepositoryMock := mocks.NewPartRepository(s.T())

			tt.setupMocks(partRepositoryMock)

			service := NewPartService(partRepositoryMock)
			part, err := service.GetPart(s.ctx, tt.uuid)

			tt.checkFn(part, err)
		})
	}
}
