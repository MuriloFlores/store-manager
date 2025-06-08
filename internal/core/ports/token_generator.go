package ports

import "github.com/muriloFlores/StoreManager/internal/core/domain"

type TokenManager interface {
	Generate(identity *domain.Identity) (string, error)
	Validate(tokenString string) (*domain.Identity, error)
}
