package jwt_manager

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/muriloFlores/StoreManager/internal/core/domain"
	"github.com/muriloFlores/StoreManager/internal/core/ports"
	"github.com/muriloFlores/StoreManager/internal/core/value_objects"
	"time"
)

type jwtGenerator struct {
	secretKey string
}

func NewJWTGenerator(secretKey string) ports.TokenManager {
	return &jwtGenerator{
		secretKey: secretKey,
	}
}

func (j *jwtGenerator) Generate(identity *domain.Identity) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)

	claims := &AppClaims{
		UserID: identity.UserID,
		Role:   string(identity.Role),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			Issuer:    "order_manager",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(j.secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (j *jwtGenerator) Validate(tokenString string) (*domain.Identity, error) {
	claims := &AppClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(j.secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	identity := &domain.Identity{
		UserID: claims.UserID,
		Role:   value_objects.Role(claims.Role),
	}

	return identity, nil
}
