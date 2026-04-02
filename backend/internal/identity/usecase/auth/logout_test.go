package auth

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestLogoutUseCase_Execute(t *testing.T) {
	tests := []struct {
		name         string
		refreshToken string
		setup        func(*MockRefreshTokenRepository)
		wantErr      error
	}{
		{
			name:         "Success",
			refreshToken: "valid-token",
			setup: func(rr *MockRefreshTokenRepository) {
				rr.On("DeleteRefreshToken", mock.Anything, "valid-token").Return(nil)
			},
			wantErr: nil,
		},
		{
			name:         "Token Not Found",
			refreshToken: "invalid-token",
			setup: func(rr *MockRefreshTokenRepository) {
				rr.On("DeleteRefreshToken", mock.Anything, "invalid-token").Return(errors.New("not found"))
			},
			wantErr: errors.New("not found"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rr := new(MockRefreshTokenRepository)
			tt.setup(rr)

			uc := NewLogoutUseCase(rr)
			err := uc.Execute(context.Background(), tt.refreshToken)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.ErrorContains(t, err, tt.wantErr.Error())
			} else {
				assert.NoError(t, err)
			}
			rr.AssertExpectations(t)
		})
	}
}
