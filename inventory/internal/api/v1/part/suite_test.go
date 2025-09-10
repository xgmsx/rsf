package part

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/xgmsx/rsf/inventory/internal/service/mocks"
)

type ServiceSuite struct {
	suite.Suite

	ctx     context.Context //nolint:containedctx
	service *mocks.PartService
	api     *partAPI
}

func (s *ServiceSuite) SetupTest() {
	s.ctx = context.Background()
	s.service = mocks.NewPartService(s.T())
	s.api = NewPartAPI(s.service)
}

func (s *ServiceSuite) TearDownTest() {}

func TestPartService(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}
