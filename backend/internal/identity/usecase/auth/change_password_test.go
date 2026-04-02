package auth

import (
	"context"
	"testing"

	"github.com/MuriloFlores/order-manager/internal/identity/domain/entity"
	"github.com/MuriloFlores/order-manager/internal/identity/domain/vo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestChangePasswordUseCase_Execute(t *testing.T) {
	pepper := "test-pepper"
	emailStr := "test@example.com"
	emailVO, _ := vo.NewEmail(emailStr)
	passwordStr := "OldPassword123!"
	passwordVO, _ := vo.NewPassword(passwordStr, pepper)
	user, _ := entity.NewUser(emailVO, "testuser", passwordVO, []vo.Role{vo.EmployeeRole})

	tests := []struct {
		name        string
		oldPassword string
		newPassword string
		setup       func(*MockUserRepository)
		wantErr     error
	}{
		{
			name:        "Success",
			oldPassword: passwordStr,
			newPassword: "NewPassword123!",
			setup: func(mr *MockUserRepository) {
				mr.On("FindByID", mock.Anything, user.ID()).Return(user, nil)
				mr.On("Update", mock.Anything, mock.MatchedBy(func(u *entity.User) bool {
					return u.Password().Matches("NewPassword123!", pepper)
				})).Return(nil)
			},
			wantErr: nil,
		},
		{
			name:        "Invalid Old Password",
			oldPassword: "wrong-password",
			newPassword: "NewPassword123!",
			setup: func(mr *MockUserRepository) {
				mr.On("FindByID", mock.Anything, user.ID()).Return(user, nil)
			},
			wantErr: entity.ErrInvalidOldPassword,
		},
		{
			name:        "User Not Found",
			oldPassword: passwordStr,
			newPassword: "NewPassword123!",
			setup: func(mr *MockUserRepository) {
				mr.On("FindByID", mock.Anything, user.ID()).Return(nil, nil)
			},
			wantErr: entity.ErrUserNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mr := new(MockUserRepository)
			tt.setup(mr)

			uc := NewChangePassword(mr, pepper)
			err := uc.Execute(context.Background(), user.ID(), tt.oldPassword, tt.newPassword)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
			}
			mr.AssertExpectations(t)
		})
	}
}
