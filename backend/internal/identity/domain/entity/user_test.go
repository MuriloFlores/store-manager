package entity

import (
	"testing"

	"github.com/MuriloFlores/order-manager/internal/identity/domain/vo"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewUser(t *testing.T) {
	email, _ := vo.NewEmail("test@test.com")
	password, _ := vo.RestorePassword("hashed")
	roles := []vo.Role{vo.AdminRole}

	u, err := NewUser(email, "username", password, roles)
	assert.NoError(t, err)
	assert.Equal(t, "username", u.Username())
	assert.True(t, u.IsActive())

	_, err = NewUser(email, "", password, roles)
	assert.ErrorIs(t, err, ErrEmptyUsername)
}

func TestUser_RoleManagement(t *testing.T) {
	email, _ := vo.NewEmail("test@test.com")
	password, _ := vo.RestorePassword("hashed")
	u, _ := NewUser(email, "user", password, []vo.Role{vo.EmployeeRole})

	u.AddRole(vo.AdminRole)
	assert.Len(t, u.Roles(), 2)

	// Test duplicate role
	u.AddRole(vo.AdminRole)
	assert.Len(t, u.Roles(), 2)

	u.RemoveRole(vo.EmployeeRole)
	assert.Len(t, u.Roles(), 1)
	assert.Equal(t, vo.AdminRole, u.Roles()[0])
}

func TestRestoreUser(t *testing.T) {
	id := uuid.New()
	u, err := RestoreUser(id, "test@test.com", "user", "hash", []string{"ADMIN"}, false)

	assert.NoError(t, err)
	assert.Equal(t, id, u.ID())
	assert.False(t, u.IsActive())
}
