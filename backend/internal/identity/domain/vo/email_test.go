package vo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewEmail(t *testing.T) {
	tests := []struct {
		name    string
		email   string
		wantErr error
	}{
		{"Valid email", "test@example.com", nil},
		{"Empty email", "", ErrEmptyEmail},
		{"Invalid format", "invalid-email", ErrInvalidEmail},
		{"Missing domain", "test@", ErrInvalidEmail},
		{"Missing user", "@example.com", ErrInvalidEmail},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewEmail(tt.email)
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.email, got.String())
			}
		})
	}
}

func TestEmail_Equals(t *testing.T) {
	e1, _ := NewEmail("test@example.com")
	e2, _ := NewEmail("test@example.com")
	e3, _ := NewEmail("other@example.com")

	assert.True(t, e1.Equals(e2))
	assert.False(t, e1.Equals(e3))
}
