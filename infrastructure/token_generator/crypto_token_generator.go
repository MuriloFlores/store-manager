package token_generator

import (
	"crypto/rand"
	"github.com/muriloFlores/StoreManager/internal/core/ports"
	"io"
)

type CryptoTokenGenerator struct{}

func NewCryptoTokenGenerator() ports.SecureTokenGenerator {
	return &CryptoTokenGenerator{}
}

func (g *CryptoTokenGenerator) Generate() (string, error) {
	table := [...]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}

	otp := make([]byte, 6)

	_, err := io.ReadFull(rand.Reader, otp)
	if err != nil {
		return "", err
	}

	for i := 0; i < len(otp); i++ {
		otp[i] = table[int(otp[i])%len(table)]
	}

	return string(otp), nil
}
