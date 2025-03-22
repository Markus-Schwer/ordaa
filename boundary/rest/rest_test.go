package rest

import (
	"context"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/suite"
	"gitlab.com/sfz.aalen/hackwerk/dotinder/boundary/auth"
	"gitlab.com/sfz.aalen/hackwerk/dotinder/boundary/utils"
	"gitlab.com/sfz.aalen/hackwerk/dotinder/entity"
	"gorm.io/gorm"
)

type Suite struct {
	suite.Suite
	ctx          context.Context
	repo         *entity.RepositoryMock
	user         *entity.User
	e            *echo.Echo
	authService  *auth.AuthService
	restBoundary *RestBoundary
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}

func (s *Suite) SetupSuite() {
	userUuid := uuid.Must(uuid.NewV4())
	s.user = &entity.User{
		Uuid: &userUuid,
		Name: "test",
	}

	s.ctx = context.Background()
	s.e = echo.New()
	s.e.Validator = utils.NewValidator()
	s.repo = &entity.RepositoryMock{
		TransactionFunc: func(fn func(tx *gorm.DB) error) error {
			return fn(nil)
		},
		GetUserFunc: func(tx *gorm.DB, uuidMoqParam *uuid.UUID) (*entity.User, error) {
			if *uuidMoqParam == userUuid {
				return s.user, nil
			}
			return nil, entity.ErrUserNotFound
		},
	}

	s.authService = auth.NewAuthService(s.ctx, s.repo)
	s.restBoundary = NewRestBoundary(s.ctx, s.repo, s.authService)
}

func (s *Suite) AfterTest(suiteName string, testName string) {
	userUuid := uuid.Must(uuid.NewV4())
	s.user = &entity.User{
		Uuid: &userUuid,
		Name: "test",
	}

	s.repo = &entity.RepositoryMock{
		TransactionFunc: func(fn func(tx *gorm.DB) error) error {
			return fn(nil)
		},
		GetUserFunc: func(tx *gorm.DB, uuidMoqParam *uuid.UUID) (*entity.User, error) {
			if *uuidMoqParam == userUuid {
				return s.user, nil
			}
			return nil, entity.ErrUserNotFound
		},
	}
	s.authService = auth.NewAuthService(s.ctx, s.repo)
	s.restBoundary = NewRestBoundary(s.ctx, s.repo, s.authService)
}
