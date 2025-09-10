package part

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/xgmsx/rsf/inventory/internal/repository/mocks"
)

type ServiceSuite struct {
	suite.Suite

	ctx      context.Context //nolint:containedctx
	partRepo *mocks.PartRepository
	service  *partService
}

func (s *ServiceSuite) SetupTest() {
	s.ctx = context.Background()
	s.partRepo = mocks.NewPartRepository(s.T())
	s.service = NewPartService(s.partRepo)
}

func (s *ServiceSuite) TearDownTest() {}

func TestPartService(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}
