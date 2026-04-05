package entity

import (
	"testing"
	"time"

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
	now := time.Now()
	u, err := RestoreUser(id, "test@test.com", "user", "hash", []string{"ADMIN"}, false, 3, &now, true)

	assert.NoError(t, err)
	assert.Equal(t, id, u.ID())
	assert.False(t, u.IsActive())
	assert.Equal(t, 3, u.FailedAttempts())
	assert.NotNil(t, u.LockedUntil())
	assert.True(t, u.EmailVerified())
}

func TestUser_AccountLockout(t *testing.T) {
	email, _ := vo.NewEmail("test@test.com")
	u, _ := NewUser(email, "user", "pass", []vo.Role{vo.EmployeeRole})

	threshold := 5
	baseDuration := 15 * time.Minute
	now := time.Now()

	t.Run("Should increment failed attempts", func(t *testing.T) {
		u.RecordFailedLogin(threshold, baseDuration, now)
		assert.Equal(t, 1, u.FailedAttempts())
		assert.Nil(t, u.LockedUntil())
	})

	t.Run("Should lock account when threshold is reached", func(t *testing.T) {
		// Já temos 1 falha, vamos para a 5ª
		for i := 0; i < 4; i++ {
			u.RecordFailedLogin(threshold, baseDuration, now)
		}

		assert.Equal(t, 5, u.FailedAttempts())
		assert.NotNil(t, u.LockedUntil())
		expectedLock := now.Add(baseDuration)
		assert.True(t, u.LockedUntil().Equal(expectedLock))
	})

	t.Run("Should implement progressive lockout (backoff)", func(t *testing.T) {
		// No 10º erro (fator 2), o tempo deve ser maior
		for i := 0; i < 5; i++ {
			u.RecordFailedLogin(threshold, baseDuration, now)
		}

		assert.Equal(t, 10, u.FailedAttempts())
		// fator = 10/5 = 2. duração = 15m * 2 = 30m
		expectedLock := now.Add(30 * time.Minute)
		assert.True(t, u.LockedUntil().Equal(expectedLock))
	})

	t.Run("Should reset failed attempts", func(t *testing.T) {
		u.ResetFailedAttempts()
		assert.Equal(t, 0, u.FailedAttempts())
		assert.Nil(t, u.LockedUntil())
	})
}
