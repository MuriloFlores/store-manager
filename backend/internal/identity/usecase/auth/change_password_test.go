package auth

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

type mockCPUserRepo struct {
	mock.Mock
}

func (m *mockCPUserRepo) Save(ctx context.Context, user *entity.User) error {
	return m.Called(ctx, user).Error(0)
}
func (m *mockCPUserRepo) FindByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) != nil {
		return args.Get(0).(*entity.User), args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *mockCPUserRepo) FindByEmail(ctx context.Context, email vo.Email) (*entity.User, error) {
	return nil, nil
}
func (m *mockCPUserRepo) FindByRole(ctx context.Context, roles []vo.Role, pagination _common.Pagination) (*_common.PaginatedResult[*entity.User], error) {
	return nil, nil
}
func (m *mockCPUserRepo) Update(ctx context.Context, user *entity.User) error {
	return m.Called(ctx, user).Error(0)
}

func TestChangePasswordUseCase_Execute(t *testing.T) {
	pepper := "pepper"
	oldPass := "OldPass123!"
	newPass := "NewPass123!"

	email, _ := vo.NewEmail("t@t.com")
	hashedOld, _ := vo.NewPassword(oldPass, pepper)
	user, _ := entity.NewUser(email, "user", hashedOld, nil)

	tests := []struct {
		name        string
		oldPassword string
		newPassword string
		setup       func(*mockCPUserRepo)
		wantErr     error
	}{
		{
			name:        "Success",
			oldPassword: oldPass,
			newPassword: newPass,
			setup: func(mr *mockCPUserRepo) {
				mr.On("FindByID", mock.Anything, user.ID()).Return(user, nil)
				mr.On("Update", mock.Anything, mock.MatchedBy(func(u *entity.User) bool {
					return u.Password().Matches(newPass, pepper)
				})).Return(nil)
			},
			wantErr: nil,
		},
		{
			name:        "Invalid Old Password",
			oldPassword: "WrongOld1!",
			newPassword: newPass,
			setup: func(mr *mockCPUserRepo) {
				mr.On("FindByID", mock.Anything, user.ID()).Return(user, nil)
			},
			wantErr: ErrInvalidOldPassword,
		},
		{
			name:        "New Password Too Short",
			oldPassword: oldPass,
			newPassword: "short",
			setup: func(mr *mockCPUserRepo) {
				mr.On("FindByID", mock.Anything, user.ID()).Return(user, nil)
			},
			wantErr: vo.ErrPasswordTooShort,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(mockCPUserRepo)
			tt.setup(mockRepo)

			uc := NewChangePassword(mockRepo, pepper)
			err := uc.Execute(context.Background(), user.ID(), tt.oldPassword, tt.newPassword)

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
				mockRepo.AssertExpectations(t)
			}
		})
	}
}
