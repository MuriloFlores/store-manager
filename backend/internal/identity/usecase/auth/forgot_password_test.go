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

type MockNotificationService struct {
	mock.Mock
}

func (m *MockNotificationService) SendForgotPasswordEmail(ctx context.Context, email vo.Email, otp vo.OTP) error {
	args := m.Called(ctx, email, otp)
	return args.Error(0)
}

func (m *MockNotificationService) SendChangePasswordEmail(ctx context.Context, email vo.Email, otp vo.OTP) error {
	args := m.Called(ctx, email, otp)
	return args.Error(0)
}

type MockOTPRepository struct {
	mock.Mock
}

func (m *MockOTPRepository) SaveOTP(ctx context.Context, email vo.Email, otp vo.OTP, expiresIn time.Duration) error {
	args := m.Called(ctx, email, otp, expiresIn)
	return args.Error(0)
}

func (m *MockOTPRepository) GetOTP(ctx context.Context, email vo.Email) (vo.OTP, error) {
	args := m.Called(ctx, email)
	return args.Get(0).(vo.OTP), args.Error(1)
}

func (m *MockOTPRepository) DeleteOTP(ctx context.Context, email vo.Email) error {
	args := m.Called(ctx, email)
	return args.Error(0)
}

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
