package part

import (
	"github.com/brianvoe/gofakeit/v7"

	"github.com/xgmsx/rsf/inventory/internal/model"
	"github.com/xgmsx/rsf/inventory/tests/testutil"
)

func (s *ServiceSuite) TestGetPart() {
	testCases := []struct {
		name        string
		partUuid    string
		gotPart     model.Part
		gotErr      error
		expectedErr error
		setupMock   func(string, model.Part, error)
	}{
		{
			name:     "Happy path",
			partUuid: gofakeit.UUID(),
			gotPart:  testutil.GetNewPart(),
			setupMock: func(uuid string, part model.Part, err error) {
				s.partRepo.On("GetPart", s.ctx, uuid).Return(part, err).Once()
			},
		},
		{
			name:        "Part not found",
			partUuid:    gofakeit.UUID(),
			gotErr:      model.ErrPartDoesNotExist,
			expectedErr: model.ErrPartDoesNotExist,
			setupMock: func(uuid string, part model.Part, err error) {
				s.partRepo.On("GetPart", s.ctx, uuid).Return(part, err).Once()
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			// arrange
			tc.setupMock(tc.partUuid, tc.gotPart, tc.gotErr)

			// act
			part, err := s.service.GetPart(s.ctx, tc.partUuid)

			// assert
			if tc.expectedErr != nil {
				s.Require().Error(err)
				s.Require().ErrorIs(err, tc.expectedErr)
				s.Require().Empty(part)
			} else {
				s.Require().NoError(err)
				s.Require().Equal(part, tc.gotPart)
			}
		})
	}
}

func (s *ServiceSuite) TestListParts() {
	testCases := []struct {
		name        string
		partFilter  *model.PartsFilter
		gotParts    []model.Part
		gotErr      error
		expectedErr error
		setupMock   func(*model.PartsFilter, []model.Part, error)
	}{
		{
			name:       "Happy path",
			partFilter: &model.PartsFilter{},
			gotParts:   testutil.GetNewParts(3),
			setupMock: func(filter *model.PartsFilter, parts []model.Part, err error) {
				s.partRepo.On("ListParts", s.ctx, filter).Return(parts, err).Once()
			},
		},
		{
			name:       "Empty list of parts",
			partFilter: &model.PartsFilter{},
			gotParts:   testutil.GetNewParts(0),
			setupMock: func(filter *model.PartsFilter, parts []model.Part, err error) {
				s.partRepo.On("ListParts", s.ctx, filter).Return(parts, err).Once()
			},
		},
		{
			name:       "Large list of parts",
			partFilter: &model.PartsFilter{},
			gotParts:   testutil.GetNewParts(100),
			setupMock: func(filter *model.PartsFilter, parts []model.Part, err error) {
				s.partRepo.On("ListParts", s.ctx, filter).Return(parts, err).Once()
			},
		},
		{
			name:        "Error case",
			partFilter:  &model.PartsFilter{},
			gotErr:      model.ErrPartDoesNotExist,
			expectedErr: model.ErrPartDoesNotExist,
			setupMock: func(filter *model.PartsFilter, parts []model.Part, err error) {
				s.partRepo.On("ListParts", s.ctx, filter).Return(parts, err).Once()
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			// arrange
			tc.setupMock(tc.partFilter, tc.gotParts, tc.gotErr)

			// act
			parts, err := s.service.ListParts(s.ctx, tc.partFilter)

			// assert
			if tc.expectedErr != nil {
				s.Require().Error(err)
				s.Require().ErrorIs(err, tc.expectedErr)
				s.Require().Nil(parts)
			} else {
				s.Require().NoError(err)
				s.Require().Equal(parts, tc.gotParts)
			}
		})
	}
}
