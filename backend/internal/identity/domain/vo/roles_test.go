package vo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewRole(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		want    Role
		wantErr error
	}{
		{"Valid ADMIN", "ADMIN", AdminRole, nil},
		{"Valid manager lowercase", "manager", ManagerRole, nil},
		{"Valid employee with spaces", "  EMPLOYEE  ", EmployeeRole, nil},
		{"Invalid role", "GOD", "", ErrInvalidRole},
		{"Empty role", "", "", ErrEmptyRole},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewRole(tt.value)
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
