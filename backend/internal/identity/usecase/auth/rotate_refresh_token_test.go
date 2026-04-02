package auth

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/MuriloFlores/order-manager/internal/identity/domain/entity"
	"github.com/MuriloFlores/order-manager/internal/identity/domain/vo"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRotateRefreshTokenUseCase_Execute(t *testing.T) {
	token := "old-refresh-token"
	userID := uuid.New()
	email, _ := vo.NewEmail("t@t.com")
	pass, _ := vo.NewPassword("Pass123!", "pepper")
	user, _ := entity.NewUser(email, "user", pass, nil)
	
	// Create a deactivated user for testing
	deactivatedUser, _ := entity.RestoreUser(userID, "t@t.com", "user", string(pass), nil, false)

	tests := []struct {
		name      string
		token     string
		setup     func(*MockUserRepository, *MockRefreshTokenRepository, *MockTokenManager)
		wantErr   error
		wantToken bool
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
			wantErr:   nil,
			wantToken: true,
		},
		{
			name:  "Invalid Token",
			token: "invalid",
			setup: func(ur *MockUserRepository, rr *MockRefreshTokenRepository, tm *MockTokenManager) {
				rr.On("GetUserIDByRefreshToken", mock.Anything, "invalid").Return(uuid.Nil, errors.New("not found"))
			},
			wantErr: errors.New("not found"),
		},
		{
			name:  "User Not Found",
			token: token,
			setup: func(ur *MockUserRepository, rr *MockRefreshTokenRepository, tm *MockTokenManager) {
				rr.On("GetUserIDByRefreshToken", mock.Anything, token).Return(userID, nil)
				ur.On("FindByID", mock.Anything, userID).Return(nil, nil)
			},
			wantErr: ErrUserNotFound,
		},
		{
			name:  "User Deactivated",
			token: token,
			setup: func(ur *MockUserRepository, rr *MockRefreshTokenRepository, tm *MockTokenManager) {
				rr.On("GetUserIDByRefreshToken", mock.Anything, token).Return(userID, nil)
				ur.On("FindByID", mock.Anything, userID).Return(deactivatedUser, nil)
			},
			wantErr: ErrUserIsDeactivated,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ur := new(MockUserRepository)
			rr := new(MockRefreshTokenRepository)
			tm := new(MockTokenManager)

			tt.setup(ur, rr, tm)

			uc := NewRotateRefreshTokenUseCase(ur, rr, tm, time.Hour)
			res, err := uc.Execute(context.Background(), tt.token)

			if tt.wantErr != nil {
				assert.ErrorContains(t, err, tt.wantErr.Error())
			} else {
				assert.NoError(t, err)
				if tt.wantToken {
					assert.NotEmpty(t, res.AccessToken)
					assert.NotEmpty(t, res.RefreshToken)
				}
			}
			
			ur.AssertExpectations(t)
			rr.AssertExpectations(t)
			tm.AssertExpectations(t)
		})
	}
}
