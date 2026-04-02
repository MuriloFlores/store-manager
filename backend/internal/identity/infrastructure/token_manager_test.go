package infrastructure

import (
	"context"
	"testing"
	"time"

	"github.com/MuriloFlores/order-manager/internal/identity/domain/entity"
	"github.com/MuriloFlores/order-manager/internal/identity/domain/vo"
	"github.com/stretchr/testify/assert"
)

func TestJWTTokenManager(t *testing.T) {
	secret := "super-secret-key-for-testing-123"
	ttl := time.Second * 2
	tm := NewJWTTokenManager(secret, ttl)

	email, _ := vo.NewEmail("test@test.com")
	pass, _ := vo.RestorePassword("hash")
	user, _ := entity.NewUser(email, "user", pass, []vo.Role{vo.AdminRole})

	t.Run("Generate and Validate Token", func(t *testing.T) {
		access, _, err := tm.GenerateTokens(context.Background(), user)
		assert.NoError(t, err)

		claims, err := tm.ValidateAccessToken(access)
		assert.NoError(t, err)
		assert.Equal(t, user.ID(), claims.UserID)
		assert.Len(t, claims.Roles, 1)
		assert.Equal(t, vo.AdminRole, claims.Roles[0])
	})

	t.Run("Expired Token", func(t *testing.T) {
		tmShort := NewJWTTokenManager(secret, time.Millisecond)
		access, _, _ := tmShort.GenerateTokens(context.Background(), user)

		time.Sleep(time.Millisecond * 10)

		_, err := tmShort.ValidateAccessToken(access)
		assert.ErrorIs(t, err, ErrExpiredToken)
	})

	t.Run("Invalid Signature", func(t *testing.T) {
		access, _, _ := tm.GenerateTokens(context.Background(), user)
		tmOther := NewJWTTokenManager("other-secret", ttl)

		_, err := tmOther.ValidateAccessToken(access)
		assert.ErrorIs(t, err, ErrInvalidToken)
	})
}
