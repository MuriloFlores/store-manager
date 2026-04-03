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
	pepper := "pepper"
	oldPass := "OldPass123!"
	newPass := "NewPass123!"

	tests := []struct {
		name        string
		oldPassword string
		newPassword string
		setup       func(*MockUserRepository, *entity.User)
		wantErr     error
	}{
		{
			name:        "Success",
			oldPassword: oldPass,
			newPassword: newPass,
			setup: func(mr *MockUserRepository, user *entity.User) {
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
			setup: func(mr *MockUserRepository, user *entity.User) {
				mr.On("FindByID", mock.Anything, user.ID()).Return(user, nil)
			},
			wantErr: entity.ErrInvalidOldPassword,
		},
		{
			name:        "New Password Too Short",
			oldPassword: oldPass,
			newPassword: "short",
			setup: func(mr *MockUserRepository, user *entity.User) {
				mr.On("FindByID", mock.Anything, user.ID()).Return(user, nil)
			},
			wantErr: vo.ErrPasswordTooShort,
		},
		{
			name:        "FindByID Error",
			oldPassword: oldPass,
			newPassword: newPass,
			setup: func(mr *MockUserRepository, user *entity.User) {
				mr.On("FindByID", mock.Anything, user.ID()).Return(nil, assert.AnError)
			},
			wantErr: assert.AnError,
		},
		{
			name:        "User Not Found",
			oldPassword: oldPass,
			newPassword: newPass,
			setup: func(mr *MockUserRepository, user *entity.User) {
				mr.On("FindByID", mock.Anything, user.ID()).Return(nil, nil)
			},
			wantErr: entity.ErrUserNotFound,
		},
		{
			name:        "Update Error",
			oldPassword: oldPass,
			newPassword: newPass,
			setup: func(mr *MockUserRepository, user *entity.User) {
				mr.On("FindByID", mock.Anything, user.ID()).Return(user, nil)
				mr.On("Update", mock.Anything, mock.Anything).Return(assert.AnError)
			},
			wantErr: assert.AnError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockUserRepository)
			mockLogger := new(MockLogger)

			email, _ := vo.NewEmail("t@t.com")
			hashedOld, _ := vo.NewPassword(oldPass, pepper)
			user, _ := entity.NewUser(email, "user", hashedOld, nil)

			tt.setup(mockRepo, user)

			uc := NewChangePassword(mockRepo, mockLogger, pepper)
			err := uc.Execute(context.Background(), user.ID(), tt.oldPassword, tt.newPassword)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
				mockRepo.AssertExpectations(t)
			}
		})
	}
}
