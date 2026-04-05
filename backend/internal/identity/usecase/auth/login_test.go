package auth

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/MuriloFlores/order-manager/internal/identity/domain/dto"
	"github.com/MuriloFlores/order-manager/internal/identity/domain/entity"
	"github.com/MuriloFlores/order-manager/internal/identity/domain/vo"
	"github.com/google/uuid"
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

	// Usuário já bloqueado para teste de lockout
	lockedTime := time.Now().Add(time.Hour)
	lockedUser, _ := entity.RestoreUser(uuid.New(), emailStr, "locked", passwordVO.String(), []string{"ADMIN"}, true, 5, &lockedTime, true)

	threshold := 5
	baseDuration := 15 * time.Minute

	tests := []struct {
		name      string
		input     *dto.LoginRequest
		setup     func(*MockUserRepository, *MockTokenManager, *MockRefreshTokenRepository)
		wantErr   bool
		expectErr error
	}{
		{
			name:  "Success - Resets failed attempts",
			input: &dto.LoginRequest{Email: emailStr, Password: passwordStr},
			setup: func(mr *MockUserRepository, tm *MockTokenManager, rr *MockRefreshTokenRepository) {
				mr.On("FindByEmail", mock.Anything, emailVO).Return(user, nil)
				tm.On("GenerateTokens", mock.Anything, user).Return("access", "refresh", nil)
				rr.On("SaveRefreshToken", mock.Anything, user.ID(), "refresh", mock.Anything).Return(nil)
				mr.On("Update", mock.Anything, mock.MatchedBy(func(u *entity.User) bool {
					return u.FailedAttempts() == 0 && u.LockedUntil() == nil
				})).Return(nil)
			},
			wantErr: false,
		},
		{
			name:  "Login Blocked - User already locked",
			input: &dto.LoginRequest{Email: emailStr, Password: passwordStr},
			setup: func(mr *MockUserRepository, tm *MockTokenManager, rr *MockRefreshTokenRepository) {
				mr.On("FindByEmail", mock.Anything, emailVO).Return(lockedUser, nil)
			},
			wantErr:   true,
			expectErr: entity.ErrUserBlocked,
		},
		{
			name:  "Invalid Password - Increments failed attempts and updates",
			input: &dto.LoginRequest{Email: emailStr, Password: "wrong-password"},
			setup: func(mr *MockUserRepository, tm *MockTokenManager, rr *MockRefreshTokenRepository) {
				mr.On("FindByEmail", mock.Anything, emailVO).Return(user, nil)
				mr.On("Update", mock.Anything, mock.MatchedBy(func(u *entity.User) bool {
					return u.FailedAttempts() > 0
				})).Return(nil)
			},
			wantErr:   true,
			expectErr: entity.ErrInvalidCredentials,
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
			name:  "Token Generation Error",
			input: &dto.LoginRequest{Email: emailStr, Password: passwordStr},
			setup: func(mr *MockUserRepository, tm *MockTokenManager, rr *MockRefreshTokenRepository) {
				mr.On("FindByEmail", mock.Anything, emailVO).Return(user, nil)
				tm.On("GenerateTokens", mock.Anything, user).Return("", "", errors.New("tm error"))
			},
			wantErr: true,
		},
		{
			name:    "Invalid Email Format",
			input:   &dto.LoginRequest{Email: "invalid", Password: passwordStr},
			setup:   func(mr *MockUserRepository, tm *MockTokenManager, rr *MockRefreshTokenRepository) {},
			wantErr: true,
		},
		{
			name:  "Repository FindByEmail Error",
			input: &dto.LoginRequest{Email: emailStr, Password: passwordStr},
			setup: func(mr *MockUserRepository, tm *MockTokenManager, rr *MockRefreshTokenRepository) {
				mr.On("FindByEmail", mock.Anything, emailVO).Return(nil, errors.New("db error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockUserRepository)
			mockTM := new(MockTokenManager)
			mockRR := new(MockRefreshTokenRepository)
			mockLogger := new(MockLogger)

			tt.setup(mockRepo, mockTM, mockRR)

			uc := NewLogin(mockRepo, mockTM, mockRR, mockLogger, pepper, baseDuration, threshold, time.Hour)
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
