package user

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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := new(MockUserRepository)
			tt.setup(m)
			uc := NewMyInfoUseCase(m)

			result, err := uc.Execute(ctx, tt.userID)

			if tt.wantErr {
				assert.Error(t)
				assert.Nil(t, result)
				if tt.err != nil {
					assert.Equal(t, tt.err, err)
				}
			} else {
				assert.NoError(t)
				assert.NotNil(t, result)
				assert.Equal(t, username, result.Username)
				assert.Equal(t, emailStr, result.Email)
				assert.Equal(t, []string{"EMPLOYEE"}, result.Role)
			}
			m.AssertExpectations(t)
		})
	}
}
