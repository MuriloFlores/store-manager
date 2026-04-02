package auth

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/MuriloFlores/order-manager/internal/_common"
	"github.com/MuriloFlores/order-manager/internal/identity/domain/dto"
	"github.com/MuriloFlores/order-manager/internal/identity/domain/entity"
	"github.com/MuriloFlores/order-manager/internal/identity/domain/vo"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
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
	if args.Get(0) != nil {
		return args.Get(0).(*entity.User), args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *MockUserRepository) FindByEmail(ctx context.Context, email vo.Email) (*entity.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) != nil {
		return args.Get(0).(*entity.User), args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *MockUserRepository) FindByRole(ctx context.Context, roles []vo.Role, pagination _common.Pagination) (*_common.PaginatedResult[*entity.User], error) {
	args := m.Called(ctx, roles, pagination)
	if args.Get(0) != nil {
		return args.Get(0).(*_common.PaginatedResult[*entity.User]), args.Error(1)
	}
	return nil, args.Error(1)
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
	if args.Get(0) != nil {
		return args.Get(0).(*dto.UserClaims), args.Error(1)
	}
	return nil, args.Error(1)
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

func TestLoginUseCase_Execute(t *testing.T) {
	pepper := "test-pepper"
	emailStr := "test@example.com"
	passStr := "Pass123!"

	rawPass, _ := vo.NewPassword(passStr, pepper)
	emailVO, _ := vo.NewEmail(emailStr)
	user, _ := entity.NewUser(emailVO, "testuser", rawPass, []vo.Role{vo.EmployeeRole})

	tests := []struct {
		name      string
		input     *dto.LoginRequest
		setup     func(*MockUserRepository, *MockTokenManager, *MockRefreshTokenRepository)
		wantErr   error
		wantToken bool
	}{
		{
			name:  "Successful login",
			input: &dto.LoginRequest{Email: emailStr, Password: passStr},
			setup: func(mr *MockUserRepository, tm *MockTokenManager, rr *MockRefreshTokenRepository) {
				mr.On("FindByEmail", mock.Anything, emailVO).Return(user, nil)
				tm.On("GenerateTokens", mock.Anything, user).Return("access", "refresh", nil)
				rr.On("SaveRefreshToken", mock.Anything, user.ID(), "refresh", mock.Anything).Return(nil)
			},
			wantErr:   nil,
			wantToken: true,
		},
		{
			name:  "User not found",
			input: &dto.LoginRequest{Email: emailStr, Password: passStr},
			setup: func(mr *MockUserRepository, tm *MockTokenManager, rr *MockRefreshTokenRepository) {
				mr.On("FindByEmail", mock.Anything, emailVO).Return(nil, nil)
			},
			wantErr: ErrInvalidCredentials,
		},
		{
			name:  "Password mismatch",
			input: &dto.LoginRequest{Email: emailStr, Password: "WrongPassword1!"},
			setup: func(mr *MockUserRepository, tm *MockTokenManager, rr *MockRefreshTokenRepository) {
				mr.On("FindByEmail", mock.Anything, emailVO).Return(user, nil)
			},
			wantErr: ErrInvalidCredentials,
		},
		{
			name:  "Token generation error",
			input: &dto.LoginRequest{Email: emailStr, Password: passStr},
			setup: func(mr *MockUserRepository, tm *MockTokenManager, rr *MockRefreshTokenRepository) {
				mr.On("FindByEmail", mock.Anything, emailVO).Return(user, nil)
				tm.On("GenerateTokens", mock.Anything, user).Return("", "", errors.New("tm error"))
			},
			wantErr: errors.New("tm error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockUserRepository)
			mockTM := new(MockTokenManager)
			mockRR := new(MockRefreshTokenRepository)

			tt.setup(mockRepo, mockTM, mockRR)

			uc := NewLogin(mockRepo, mockTM, mockRR, pepper, time.Hour)
			res, err := uc.Execute(context.Background(), tt.input)

			if tt.wantErr != nil {
				assert.ErrorContains(t, err, tt.wantErr.Error())
			} else {
				assert.NoError(t, err)
				if tt.wantToken {
					assert.NotEmpty(t, res.AccessToken)
					assert.NotEmpty(t, res.RefreshToken)
				}
				mockRepo.AssertExpectations(t)
				mockTM.AssertExpectations(t)
				mockRR.AssertExpectations(t)
			}
		})
	}
}
