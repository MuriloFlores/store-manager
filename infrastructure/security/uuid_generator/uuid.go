package uuid_generator

import (
	"github.com/google/uuid"
	"github.com/muriloFlores/StoreManager/internal/core/ports"
)

type UUIDGenerator struct{}

func NewUUIDGenerator() ports.IDGenerator {
	return &UUIDGenerator{}
}

func (g *UUIDGenerator) Generate() string {
	return uuid.NewString()
}

func (g *UUIDGenerator) Validate(id string) bool {
	if err := uuid.Validate(id); err != nil {
		return false
	}

	return true
}
