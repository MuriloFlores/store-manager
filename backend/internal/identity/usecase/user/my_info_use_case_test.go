package user

import (
	"context"
	"testing"

	"github.com/MuriloFlores/order-manager/internal/identity/domain/entity"
	"github.com/MuriloFlores/order-manager/internal/identity/domain/vo"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestMyInfoUseCase_Execute(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	emailStr := "test@example.com"
	username := "testuser"
	
	setupUser := func() *entity.User {
		password, _ := vo.NewPassword("Password123!", "pepper")
		user, _ := entity.RestoreUser(userID, emailStr, username, password.String(), []string{"EMPLOYEE"}, true)
		return user
	}

	tests := []struct {
		name    string
		userID  uuid.UUID
		setup   func(m *MockUserRepository)
		wantErr bool
		err     error
	}{
		{
			name:   "Success",
			userID: userID,
			setup: func(m *MockUserRepository) {
				m.On("FindByID", ctx, userID).Return(setupUser(), nil)
			},
			wantErr: false,
		},
		{
			name:   "User Not Found",
			userID: userID,
			setup: func(m *MockUserRepository) {
				m.On("FindByID", ctx, userID).Return(nil, nil)
			},
			wantErr: true,
			err:     entity.ErrUserNotFound,
		},
		{
			name:   "Repository Error",
			userID: userID,
			setup: func(m *MockUserRepository) {
				m.On("FindByID", ctx, userID).Return(nil, assert.AnError)
			},
			wantErr: true,
			err:     assert.AnError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := new(MockUserRepository)
			l := new(MockLogger)
			tt.setup(m)
			uc := NewMyInfoUseCase(m, l)

			result, err := uc.Execute(ctx, tt.userID)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, result)
				if tt.err != nil {
					assert.ErrorIs(t, err, tt.err)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, username, result.Username)
				assert.Equal(t, emailStr, result.Email)
				assert.Equal(t, []string{"EMPLOYEE"}, result.Role)
			}
			m.AssertExpectations(t)
		})
	}
}
