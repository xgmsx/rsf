package part

import (
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/xgmsx/rsf/inventory/internal/model"
	"github.com/xgmsx/rsf/inventory/internal/model/converter"
	"github.com/xgmsx/rsf/inventory/tests/testutil"
	genInventoryV1 "github.com/xgmsx/rsf/shared/pkg/proto/inventory/v1"
)

func (s *ServiceSuite) TestListPartsHandlerSuccess() {
	testCases := []struct {
		name         string
		gotFilter    *genInventoryV1.PartsFilter
		gotParts     []model.Part
		gotErr       error
		expectedCode codes.Code
		setupMock    func(*model.PartsFilter, []model.Part, error)
	}{
		{
			name:     "Happy path",
			gotParts: testutil.GetNewParts(3),
			setupMock: func(filter *model.PartsFilter, parts []model.Part, err error) {
				s.service.On("ListParts", s.ctx, filter).Return(parts, err).Once()
			},
		},
		{
			name:         "No parts found",
			expectedCode: codes.NotFound,
			setupMock: func(filter *model.PartsFilter, parts []model.Part, err error) {
				s.service.On("ListParts", s.ctx, filter).Return(parts, err).Once()
			},
		},
		{
			name:         "Internal error",
			gotErr:       fmt.Errorf("test error"),
			expectedCode: codes.Internal,
			setupMock: func(filter *model.PartsFilter, parts []model.Part, err error) {
				s.service.On("ListParts", s.ctx, filter).Return(parts, err).Once()
			},
		},
		{
			name:         "Request validation error",
			gotFilter:    &genInventoryV1.PartsFilter{Uuids: []string{"invalid-uuid 1"}},
			expectedCode: codes.InvalidArgument,
			setupMock:    func(filter *model.PartsFilter, parts []model.Part, err error) {},
		},
		{
			name:         "Response validation error",
			expectedCode: codes.Internal,
			gotParts:     []model.Part{{UUID: "invalid-uuid 1"}, {UUID: "invalid-uuid 2"}, {UUID: "invalid-uuid 3"}},
			setupMock: func(filter *model.PartsFilter, parts []model.Part, err error) {
				s.service.On("ListParts", s.ctx, filter).Return(parts, err).Once()
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			// arrange
			tc.setupMock(converter.PartFilterFromProto(tc.gotFilter), tc.gotParts, tc.gotErr)

			// act
			resp, err := s.api.ListParts(s.ctx, &genInventoryV1.ListPartsRequest{Filter: tc.gotFilter})

			// assert
			if tc.expectedCode == codes.OK {
				s.Require().NoError(err)
				s.Require().NotNil(resp)
				s.Require().NoError(resp.ValidateAll())
				s.Require().EqualValues(converter.PartsToProto(tc.gotParts), resp.Parts)
			} else {
				s.Require().Error(err)
				s.Require().Nil(resp)
				s.Require().NoError(resp.ValidateAll())

				st, ok := status.FromError(err)
				s.Require().True(ok)
				s.Require().Equal(tc.expectedCode, st.Code())
			}
		})
	}
}
