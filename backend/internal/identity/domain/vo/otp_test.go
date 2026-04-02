package vo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewOTP(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		wantErr error
	}{
		{"Valid OTP", "123456", nil},
		{"Too short", "12345", ErrInvalidOTPFormat},
		{"Too long", "1234567", ErrInvalidOTPFormat},
		{"Non-numeric", "123a56", ErrInvalidOTPFormat},
		{"Empty", "", ErrInvalidOTPFormat},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewOTP(tt.value)
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.value, got.String())
			}
		})
	}
}

func TestGenerateOTP(t *testing.T) {
	otp, err := GenerateOTP()
	assert.NoError(t, err)
	assert.Len(t, otp.String(), 6)

	_, err = NewOTP(otp.String())
	assert.NoError(t, err, "Generated OTP is invalid")
}
