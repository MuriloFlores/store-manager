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

func TestChangeUserRoleUseCase_Execute(t *testing.T) {
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
		roles   []string
		setup   func(m *MockUserRepository)
		wantErr bool
		err     error
	}{
		{
			name:  "Success Change Roles",
			id:    userID.String(),
			roles: []string{"ADMIN", "MANAGER"},
			setup: func(m *MockUserRepository) {
				user := setupUser()
				m.On("FindByID", ctx, userID).Return(user, nil)
				m.On("Update", ctx, mock.MatchedBy(func(u *entity.User) bool {
					roles := u.Roles()
					return len(roles) == 2 && roles[0] == vo.AdminRole && roles[1] == vo.ManagerRole
				})).Return(nil)
			},
			wantErr: false,
		},
		{
			name:  "Invalid Role Name",
			id:    userID.String(),
			roles: []string{"INVALID_ROLE"},
			setup: func(m *MockUserRepository) {
				// No FindByID expected as NewRole should fail first
			},
			wantErr: true,
			err:     vo.ErrInvalidRole,
		},
		{
			name:  "User Not Found",
			id:    userID.String(),
			roles: []string{"ADMIN"},
			setup: func(m *MockUserRepository) {
				m.On("FindByID", ctx, userID).Return(nil, nil)
			},
			wantErr: true,
			err:     entity.ErrUserNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := new(MockUserRepository)
			tt.setup(m)
			uc := NewChangeUserRoleUseCase(m)

			err := uc.Execute(ctx, tt.id, tt.roles)

			if tt.wantErr {
				assert.Error(t)
				if tt.err != nil {
					assert.Equal(t, tt.err, err)
				}
			} else {
				assert.NoError(t)
			}
			m.AssertExpectations(t)
		})
	}
}
