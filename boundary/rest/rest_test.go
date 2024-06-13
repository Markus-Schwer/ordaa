package rest

import (
	"context"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/suite"
	"gitlab.com/sfz.aalen/hackwerk/dotinder/boundary/auth"
	"gitlab.com/sfz.aalen/hackwerk/dotinder/boundary/utils"
	"gitlab.com/sfz.aalen/hackwerk/dotinder/entity"
)

type Suite struct {
	suite.Suite
	ctx context.Context
	e *echo.Echo
	repo entity.Repository
	authService *auth.AuthService
	restBoundary *RestBoundary
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}

func (s *Suite) SetupSuite() {
	s.ctx = context.Background()
	s.e = echo.New()
	s.e.Validator = utils.NewValidator()
	s.repo = entity.NewMockRepository()
	s.authService = auth.NewAuthService(s.ctx, s.repo)
	s.restBoundary = NewRestBoundary(s.ctx, s.repo, s.authService)
}

func (s *Suite) AfterTest(suiteName string, testName string) {
	s.repo = entity.NewMockRepository()
	s.authService = auth.NewAuthService(s.ctx, s.repo)
	s.restBoundary = NewRestBoundary(s.ctx, s.repo, s.authService)
}
