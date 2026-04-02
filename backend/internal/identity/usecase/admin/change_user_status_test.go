package admin

import (
	"context"
	"testing"

	"github.com/MuriloFlores/order-manager/internal/_common"
	"github.com/MuriloFlores/order-manager/internal/identity/domain/entity"
	"github.com/MuriloFlores/order-manager/internal/identity/domain/vo"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUserRepository implements ports.UserRepository for testing
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

func TestChangeUserStatusUseCase_Execute(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	email, _ := vo.NewEmail("test@example.com")
	password, _ := vo.NewPassword("Password123!", "pepper")
	
	setupUser := func() *entity.User {
		user, _ := entity.NewUser(email, "testuser", password, []vo.Role{vo.EmployeeRole})
		// Re-create to match the specific ID for easier testing
		user, _ = entity.RestoreUser(userID, email.String(), "testuser", password.String(), []string{vo.EmployeeRole.String()}, true)
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := new(MockUserRepository)
			tt.setup(m)
			uc := NewChangeUserStatusUseCase(m)

			err := uc.Execute(ctx, tt.id, tt.active)

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
