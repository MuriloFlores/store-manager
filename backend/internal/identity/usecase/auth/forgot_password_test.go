package auth

import (
	"context"
	"testing"
	"time"

	"github.com/MuriloFlores/order-manager/internal/identity/domain/entity"
	"github.com/MuriloFlores/order-manager/internal/identity/domain/vo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestForgotPasswordUseCase_Execute(t *testing.T) {
	emailStr := "test@example.com"
	emailVO, _ := vo.NewEmail(emailStr)
	user, _ := entity.NewUser(emailVO, "testuser", vo.Password(""), []vo.Role{vo.EmployeeRole})

	tests := []struct {
		name      string
		email     string
		setup     func(*MockUserRepository, *MockOTPRepository, *MockNotificationService)
		wantErr   error
		expectErr bool
	}{
		{
			name:  "Success",
			email: emailStr,
			setup: func(ur *MockUserRepository, or *MockOTPRepository, ns *MockNotificationService) {
				ur.On("FindByEmail", mock.Anything, emailVO).Return(user, nil)
				or.On("SaveOTP", mock.Anything, emailVO, mock.Anything, mock.Anything).Return(nil)
				ns.On("SendForgotPasswordEmail", mock.Anything, emailVO, mock.Anything).Return(nil)
			},
			expectErr: false,
		},
		{
			name:  "User Not Found",
			email: emailStr,
			setup: func(ur *MockUserRepository, or *MockOTPRepository, ns *MockNotificationService) {
				ur.On("FindByEmail", mock.Anything, emailVO).Return(nil, nil)
			},
			wantErr:   nil, // O use case atual retorna nil mesmo se nao encontrar (para evitar enumeraçao)
			expectErr: false,
		},
		{
			name:    "Invalid Email",
			email:   "invalid",
			setup:   func(ur *MockUserRepository, or *MockOTPRepository, ns *MockNotificationService) {},
			expectErr: true,
		},
		{
			name:  "Save OTP Error",
			email: emailStr,
			setup: func(ur *MockUserRepository, or *MockOTPRepository, ns *MockNotificationService) {
				ur.On("FindByEmail", mock.Anything, emailVO).Return(user, nil)
				or.On("SaveOTP", mock.Anything, emailVO, mock.Anything, mock.Anything).Return(assert.AnError)
			},
			expectErr: true,
		},
		{
			name:  "Send Email Error",
			email: emailStr,
			setup: func(ur *MockUserRepository, or *MockOTPRepository, ns *MockNotificationService) {
				ur.On("FindByEmail", mock.Anything, emailVO).Return(user, nil)
				or.On("SaveOTP", mock.Anything, emailVO, mock.Anything, mock.Anything).Return(nil)
				ns.On("SendForgotPasswordEmail", mock.Anything, emailVO, mock.Anything).Return(assert.AnError)
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ur := new(MockUserRepository)
			or := new(MockOTPRepository)
			ns := new(MockNotificationService)

			tt.setup(ur, or, ns)

			uc := NewForgotPassword(or, ur, ns, time.Hour)
			err := uc.Execute(context.Background(), tt.email)

			if tt.expectErr {
				assert.Error(t, err)
				if tt.wantErr != nil {
					assert.ErrorIs(t, err, tt.wantErr)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
