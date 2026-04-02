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

type MockOTPRepository struct {
	mock.Mock
}

func (m *MockOTPRepository) SaveOTP(ctx context.Context, email vo.Email, otp vo.OTP, expiresIn time.Duration) error {
	return m.Called(ctx, email, otp, expiresIn).Error(0)
}
func (m *MockOTPRepository) GetOTP(ctx context.Context, email vo.Email) (vo.OTP, error) {
	args := m.Called(ctx, email)
	return args.Get(0).(vo.OTP), args.Error(1)
}
func (m *MockOTPRepository) DeleteOTP(ctx context.Context, email vo.Email) error {
	return m.Called(ctx, email).Error(0)
}

type MockNotificationService struct {
	mock.Mock
}

func (m *MockNotificationService) SendChangePasswordEmail(ctx context.Context, toEmail vo.Email, resetToken vo.OTP) error {
	return m.Called(ctx, toEmail, resetToken).Error(0)
}
func (m *MockNotificationService) SendForgotPasswordEmail(ctx context.Context, toEmail vo.Email, resetToken vo.OTP) error {
	return m.Called(ctx, toEmail, resetToken).Error(0)
}

func TestForgotPasswordUseCase_Execute(t *testing.T) {
	emailStr := "test@example.com"
	emailVO, _ := vo.NewEmail(emailStr)
	pass, _ := vo.NewPassword("Pass123!", "pepper")
	user, _ := entity.NewUser(emailVO, "user", pass, nil)

	tests := []struct {
		name    string
		email   string
		setup   func(*MockUserRepository, *MockOTPRepository, *MockNotificationService)
		wantErr bool
	}{
		{
			name:  "Success",
			email: emailStr,
			setup: func(ur *MockUserRepository, or *MockOTPRepository, ns *MockNotificationService) {
				ur.On("FindByEmail", mock.Anything, emailVO).Return(user, nil)
				or.On("SaveOTP", mock.Anything, emailVO, mock.AnythingOfType("vo.OTP"), mock.Anything).Return(nil)
				ns.On("SendForgotPasswordEmail", mock.Anything, emailVO, mock.AnythingOfType("vo.OTP")).Return(nil)
			},
			wantErr: false,
		},
		{
			name:  "Invalid Email Format",
			email: "invalid",
			setup: func(ur *MockUserRepository, or *MockOTPRepository, ns *MockNotificationService) {},
			wantErr: true,
		},
		{
			name:  "User Not Found (Silent Success)",
			email: emailStr,
			setup: func(ur *MockUserRepository, or *MockOTPRepository, ns *MockNotificationService) {
				ur.On("FindByEmail", mock.Anything, emailVO).Return(nil, nil)
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ur := new(MockUserRepository)
			or := new(MockOTPRepository)
			ns := new(MockNotificationService)

			tt.setup(ur, or, ns)

			uc := NewForgotPassword(or, ur, ns, time.Minute*15)
			err := uc.Execute(context.Background(), tt.email)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				ur.AssertExpectations(t)
				or.AssertExpectations(t)
				ns.AssertExpectations(t)
			}
		})
	}
}
