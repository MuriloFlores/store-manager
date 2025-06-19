package reset_token

import (
	"crypto/rand"
	"encoding/hex"
	"github.com/muriloFlores/StoreManager/internal/core/ports"
)

type CryptoTokenGenerator struct{}

func NewCryptoTokenGenerator() ports.SecureTokenGenerator {
	return &CryptoTokenGenerator{}
}

func (g *CryptoTokenGenerator) Generate() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	return hex.EncodeToString(bytes), nil
}
