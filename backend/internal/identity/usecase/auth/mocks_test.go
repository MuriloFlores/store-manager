package auth

import (
	"context"
	"time"

	"github.com/MuriloFlores/order-manager/internal/_common"
	"github.com/MuriloFlores/order-manager/internal/identity/domain/dto"
	"github.com/MuriloFlores/order-manager/internal/identity/domain/entity"
	"github.com/MuriloFlores/order-manager/internal/identity/domain/vo"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

// MockUserRepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Save(ctx context.Context, user *entity.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserRepository) FindByEmail(ctx context.Context, email vo.Email) (*entity.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserRepository) GetUsersInfo(ctx context.Context, roles []vo.Role, pagination _common.Pagination) (*_common.PaginatedResult[*entity.User], error) {
	args := m.Called(ctx, roles, pagination)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*_common.PaginatedResult[*entity.User]), args.Error(1)
}

func (m *MockUserRepository) Update(ctx context.Context, user *entity.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

// MockTokenManager
type MockTokenManager struct {
	mock.Mock
}

func (m *MockTokenManager) GenerateTokens(ctx context.Context, user *entity.User) (string, string, error) {
	args := m.Called(ctx, user)
	return args.String(0), args.String(1), args.Error(2)
}

func (m *MockTokenManager) ValidateAccessToken(tokenString string) (*dto.UserClaims, error) {
	args := m.Called(tokenString)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.UserClaims), args.Error(1)
}

// MockRefreshTokenRepository
type MockRefreshTokenRepository struct {
	mock.Mock
}

func (m *MockRefreshTokenRepository) SaveRefreshToken(ctx context.Context, userID uuid.UUID, refreshToken string, expiresIn time.Duration) error {
	args := m.Called(ctx, userID, refreshToken, expiresIn)
	return args.Error(0)
}

func (m *MockRefreshTokenRepository) GetUserIDByRefreshToken(ctx context.Context, refreshToken string) (uuid.UUID, error) {
	args := m.Called(ctx, refreshToken)
	return args.Get(0).(uuid.UUID), args.Error(1)
}

func (m *MockRefreshTokenRepository) DeleteRefreshToken(ctx context.Context, refreshToken string) error {
	args := m.Called(ctx, refreshToken)
	return args.Error(0)
}
