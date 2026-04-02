package user

import (
	"context"
	"testing"

	"github.com/MuriloFlores/order-manager/internal/_common"
	"github.com/MuriloFlores/order-manager/internal/identity/domain/dto"
	"github.com/MuriloFlores/order-manager/internal/identity/domain/entity"
	"github.com/MuriloFlores/order-manager/internal/identity/domain/vo"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUserRepository implements ports.UserRepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Save(ctx context.Context, user *entity.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) != nil {
		return args.Get(0).(*entity.User), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockUserRepository) FindByEmail(ctx context.Context, email vo.Email) (*entity.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) != nil {
		return args.Get(0).(*entity.User), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockUserRepository) FindByRole(ctx context.Context, roles []vo.Role, pagination _common.Pagination) (*_common.PaginatedResult[*entity.User], error) {
	args := m.Called(ctx, roles, pagination)
	if args.Get(0) != nil {
		return args.Get(0).(*_common.PaginatedResult[*entity.User]), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockUserRepository) Update(ctx context.Context, user *entity.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func TestCreateUserUseCase_Execute(t *testing.T) {
	pepper := "test-pepper"

	tests := []struct {
		name    string
		input   dto.CreateUserInput
		setup   func(mockRepo *MockUserRepository)
		wantErr bool
	}{
		{
			name: "Success",
			input: dto.CreateUserInput{
				Username: "murilo",
				Email:    "murilo@test.com",
				Password: "StrongPass123!",
				Roles:    []string{"ADMIN"},
			},
			setup: func(mockRepo *MockUserRepository) {
				mockRepo.On("Save", mock.Anything, mock.MatchedBy(func(u *entity.User) bool {
					return u.Username() == "murilo" && u.Email().String() == "murilo@test.com"
				})).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "Invalid Email",
			input: dto.CreateUserInput{
				Username: "murilo",
				Email:    "invalid-email",
				Password: "StrongPass123!",
				Roles:    []string{"ADMIN"},
			},
			setup: func(mockRepo *MockUserRepository) {
				// Save should not be called
			},
			wantErr: true,
		},
		{
			name: "Invalid Role",
			input: dto.CreateUserInput{
				Username: "murilo",
				Email:    "murilo@test.com",
				Password: "StrongPass123!",
				Roles:    []string{"INVALID_ROLE"},
			},
			setup: func(mockRepo *MockUserRepository) {
				// Save should not be called
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockUserRepository)
			tt.setup(mockRepo)

			uc := NewCreateUserService(mockRepo, pepper)
			err := uc.Execute(context.Background(), tt.input)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				mockRepo.AssertExpectations(t)
			}
		})
	}
}
