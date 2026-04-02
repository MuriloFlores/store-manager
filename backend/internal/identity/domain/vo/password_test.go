package vo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const testPepper = "secret-pepper"

func TestNewPassword(t *testing.T) {
	tests := []struct {
		name    string
		plain   string
		wantErr error
	}{
		{"Valid password", "Pass123!", nil},
		{"Too short", "P1!", ErrPasswordTooShort},
		{"Low complexity - no upper", "pass123!", ErrLowPasswordComplexity},
		{"Low complexity - no number", "Password!", ErrLowPasswordComplexity},
		{"Low complexity - no special", "Password123", ErrLowPasswordComplexity},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewPassword(tt.plain, testPepper)
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestPassword_Matches(t *testing.T) {
	plain := "StrongPass123!"
	pwd, err := NewPassword(plain, testPepper)
	assert.NoError(t, err)

	assert.True(t, pwd.Matches(plain, testPepper))
	assert.False(t, pwd.Matches("WrongPass123!", testPepper))
	assert.False(t, pwd.Matches(plain, "wrong-pepper"))
}

func TestRestorePassword(t *testing.T) {
	hash := "$argon2id$v=19$m=65536,t=3,p=4$salt$hash"
	pwd, err := RestorePassword(hash)
	assert.NoError(t, err)
	assert.Equal(t, hash, pwd.String())

	_, err = RestorePassword("")
	assert.ErrorIs(t, err, ErrEmptyPassword)
}
