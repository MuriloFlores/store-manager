package admin

import (
	"context"
	"testing"

	"github.com/MuriloFlores/order-manager/internal/identity/domain/entity"
	"github.com/MuriloFlores/order-manager/internal/identity/domain/vo"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestChangeUserStatusUseCase_Execute(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	email, _ := vo.NewEmail("test@example.com")
	password, _ := vo.NewPassword("Password123!", "pepper")
	
	setupUser := func() *entity.User {
		user, _ := entity.RestoreUser(userID, email.String(), "testuser", password.String(), []string{vo.EmployeeRole.String()}, true)
		return user
	}

	tests := []struct {
		name    string
		id      string
		active  bool
		setup   func(m *MockUserRepository)
		wantErr bool
		err     error
	}{
		{
			name:   "Success Deactivate",
			id:     userID.String(),
			active: false,
			setup: func(m *MockUserRepository) {
				user := setupUser()
				m.On("FindByID", ctx, userID).Return(user, nil)
				m.On("Update", ctx, mock.MatchedBy(func(u *entity.User) bool {
					return u.ID() == userID && !u.IsActive()
				})).Return(nil)
			},
			wantErr: false,
		},
		{
			name:   "Success Activate",
			id:     userID.String(),
			active: true,
			setup: func(m *MockUserRepository) {
				user := setupUser()
				user.Deactivate()
				m.On("FindByID", ctx, userID).Return(user, nil)
				m.On("Update", ctx, mock.MatchedBy(func(u *entity.User) bool {
					return u.ID() == userID && u.IsActive()
				})).Return(nil)
			},
			wantErr: false,
		},
		{
			name:   "User Not Found",
			id:     userID.String(),
			active: true,
			setup: func(m *MockUserRepository) {
				m.On("FindByID", ctx, userID).Return(nil, nil)
			},
			wantErr: true,
			err:     entity.ErrUserNotFound,
		},
		{
			name:    "Invalid UUID",
			id:      "invalid-uuid",
			active:  true,
			setup:   func(m *MockUserRepository) {},
			wantErr: true,
		},
		{
			name:   "FindByID Error",
			id:     userID.String(),
			active: true,
			setup: func(m *MockUserRepository) {
				m.On("FindByID", ctx, userID).Return(nil, assert.AnError)
			},
			wantErr: true,
			err:     assert.AnError,
		},
		{
			name:   "Update Error",
			id:     userID.String(),
			active: true,
			setup: func(m *MockUserRepository) {
				user := setupUser()
				m.On("FindByID", ctx, userID).Return(user, nil)
				m.On("Update", ctx, mock.Anything).Return(assert.AnError)
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
			uc := NewChangeUserStatusUseCase(m, l)

			err := uc.Execute(ctx, tt.id, tt.active)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.err != nil {
					assert.ErrorIs(t, err, tt.err)
				}
			} else {
				assert.NoError(t, err)
			}
			m.AssertExpectations(t)
		})
	}
}
