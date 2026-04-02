package auth

import (
	"context"
	"testing"

	"github.com/MuriloFlores/order-manager/internal/identity/domain/entity"
	"github.com/MuriloFlores/order-manager/internal/identity/domain/vo"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRotateRefreshTokenUseCase_Execute(t *testing.T) {
	userID := uuid.New()
	token := "valid-refresh-token"
	email, _ := vo.NewEmail("test@example.com")
	password, _ := vo.NewPassword("Password123!", "pepper")
	user, _ := entity.RestoreUser(userID, email.String(), "testuser", password.String(), []string{vo.EmployeeRole.String()}, true)
	deactivatedUser, _ := entity.RestoreUser(userID, email.String(), "testuser", password.String(), []string{vo.EmployeeRole.String()}, false)

	tests := []struct {
		name      string
		token     string
		setup     func(*MockUserRepository, *MockRefreshTokenRepository, *MockTokenManager)
		wantErr   error
		expectErr bool
	}{
		{
			name:  "Success",
			token: token,
			setup: func(ur *MockUserRepository, rr *MockRefreshTokenRepository, tm *MockTokenManager) {
				rr.On("GetUserIDByRefreshToken", mock.Anything, token).Return(userID, nil)
				ur.On("FindByID", mock.Anything, userID).Return(user, nil)
				rr.On("DeleteRefreshToken", mock.Anything, token).Return(nil)
				tm.On("GenerateTokens", mock.Anything, user).Return("new-access", "new-refresh", nil)
				rr.On("SaveRefreshToken", mock.Anything, mock.Anything, "new-refresh", mock.Anything).Return(nil)
			},
			expectErr: false,
		},
		{
			name:  "Invalid Token",
			token: "invalid",
			setup: func(ur *MockUserRepository, rr *MockRefreshTokenRepository, tm *MockTokenManager) {
				rr.On("GetUserIDByRefreshToken", mock.Anything, "invalid").Return(uuid.Nil, entity.ErrSessionNotFound)
			},
			wantErr:   entity.ErrSessionNotFound,
			expectErr: true,
		},
		{
			name:  "User Not Found",
			token: token,
			setup: func(ur *MockUserRepository, rr *MockRefreshTokenRepository, tm *MockTokenManager) {
				rr.On("GetUserIDByRefreshToken", mock.Anything, token).Return(userID, nil)
				ur.On("FindByID", mock.Anything, userID).Return(nil, nil)
			},
			wantErr:   entity.ErrUserNotFound,
			expectErr: true,
		},
		{
			name:  "User Deactivated",
			token: token,
			setup: func(ur *MockUserRepository, rr *MockRefreshTokenRepository, tm *MockTokenManager) {
				rr.On("GetUserIDByRefreshToken", mock.Anything, token).Return(userID, nil)
				ur.On("FindByID", mock.Anything, userID).Return(deactivatedUser, nil)
			},
			wantErr:   entity.ErrUserIsDeactivated,
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ur := new(MockUserRepository)
			rr := new(MockRefreshTokenRepository)
			tm := new(MockTokenManager)

			tt.setup(ur, rr, tm)

			uc := NewRotateRefreshTokenUseCase(ur, rr, tm, 0)
			res, err := uc.Execute(context.Background(), tt.token)

			if tt.expectErr {
				assert.Error(t, err)
				if tt.wantErr != nil {
					assert.ErrorIs(t, err, tt.wantErr)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, res)
			}
		})
	}
}
