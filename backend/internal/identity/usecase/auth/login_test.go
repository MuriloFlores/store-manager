package auth

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/MuriloFlores/order-manager/internal/identity/domain/dto"
	"github.com/MuriloFlores/order-manager/internal/identity/domain/entity"
	"github.com/MuriloFlores/order-manager/internal/identity/domain/vo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestLoginUseCase_Execute(t *testing.T) {
	pepper := "test-pepper"
	emailStr := "test@example.com"
	emailVO, _ := vo.NewEmail(emailStr)
	passwordStr := "Password123!"
	passwordVO, _ := vo.NewPassword(passwordStr, pepper)
	user, _ := entity.NewUser(emailVO, "testuser", passwordVO, []vo.Role{vo.EmployeeRole})

	tests := []struct {
		name      string
		input     *dto.LoginRequest
		setup     func(*MockUserRepository, *MockTokenManager, *MockRefreshTokenRepository)
		wantErr   bool
		expectErr error
	}{
		{
			name:  "Success",
			input: &dto.LoginRequest{Email: emailStr, Password: passwordStr},
			setup: func(mr *MockUserRepository, tm *MockTokenManager, rr *MockRefreshTokenRepository) {
				mr.On("FindByEmail", mock.Anything, emailVO).Return(user, nil)
				tm.On("GenerateTokens", mock.Anything, user).Return("access", "refresh", nil)
				rr.On("SaveRefreshToken", mock.Anything, user.ID(), "refresh", mock.Anything).Return(nil)
			},
			wantErr: false,
		},
		{
			name:  "User Not Found",
			input: &dto.LoginRequest{Email: emailStr, Password: passwordStr},
			setup: func(mr *MockUserRepository, tm *MockTokenManager, rr *MockRefreshTokenRepository) {
				mr.On("FindByEmail", mock.Anything, emailVO).Return(nil, nil)
			},
			wantErr:   true,
			expectErr: entity.ErrInvalidCredentials,
		},
		{
			name:  "Invalid Password",
			input: &dto.LoginRequest{Email: emailStr, Password: "wrong-password"},
			setup: func(mr *MockUserRepository, tm *MockTokenManager, rr *MockRefreshTokenRepository) {
				mr.On("FindByEmail", mock.Anything, emailVO).Return(user, nil)
			},
			wantErr:   true,
			expectErr: entity.ErrInvalidCredentials,
		},
		{
			name:  "Token Generation Error",
			input: &dto.LoginRequest{Email: emailStr, Password: passwordStr},
			setup: func(mr *MockUserRepository, tm *MockTokenManager, rr *MockRefreshTokenRepository) {
				mr.On("FindByEmail", mock.Anything, emailVO).Return(user, nil)
				tm.On("GenerateTokens", mock.Anything, user).Return("", "", errors.New("tm error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockUserRepository)
			mockTM := new(MockTokenManager)
			mockRR := new(MockRefreshTokenRepository)

			tt.setup(mockRepo, mockTM, mockRR)

			uc := NewLogin(mockRepo, mockTM, mockRR, pepper, time.Hour)
			result, err := uc.Execute(context.Background(), tt.input)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.expectErr != nil {
					assert.ErrorIs(t, err, tt.expectErr)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, "access", result.AccessToken)
				assert.Equal(t, "refresh", result.RefreshToken)
			}

			mockRepo.AssertExpectations(t)
			mockTM.AssertExpectations(t)
			mockRR.AssertExpectations(t)
		})
	}
}
