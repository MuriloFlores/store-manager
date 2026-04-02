package infrastructure

import (
	"context"
	"errors"
	"time"

	"github.com/MuriloFlores/order-manager/internal/identity/domain/dto"
	"github.com/MuriloFlores/order-manager/internal/identity/domain/entity"
	"github.com/MuriloFlores/order-manager/internal/identity/domain/vo"
	"github.com/MuriloFlores/order-manager/internal/identity/ports/auth"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var (
	ErrInSigningProcess = errors.New("error in signing process")
	ErrUnexpectedMethod = errors.New("error unexpected signing method")
	ErrInvalidToken     = errors.New("invalid token")
	ErrExpiredToken     = errors.New("expired token")
)

type jwtCustomClaims struct {
	UserID string   `json:"user_id"`
	Roles  []string `json:"roles"`
	jwt.RegisteredClaims
}

type jwtTokenManager struct {
	secretKey      []byte
	accessTokenTTL time.Duration
}

func NewJWTTokenManager(secretKey string, accessTokenTTL time.Duration) auth.TokenManager {
	return &jwtTokenManager{
		secretKey:      []byte(secretKey),
		accessTokenTTL: accessTokenTTL,
	}
}

func (m *jwtTokenManager) GenerateTokens(ctx context.Context, user *entity.User) (string, string, error) {
	rolesStr := make([]string, 0, len(user.Roles()))
	for _, role := range user.Roles() {
		rolesStr = append(rolesStr, role.String())
	}

	claims := jwtCustomClaims{
		UserID: user.ID().String(),
		Roles:  rolesStr,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(m.accessTokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "store-manager",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	accessToken, err := token.SignedString(m.secretKey)
	if err != nil {
		return "", "", ErrInSigningProcess
	}

	return accessToken, uuid.New().String(), nil
}

func (m *jwtTokenManager) ValidateAccessToken(tokenString string) (*dto.UserClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwtCustomClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrUnexpectedMethod
		}

		return m.secretKey, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}

		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*jwtCustomClaims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return nil, ErrInvalidToken
	}

	roles := make([]vo.Role, 0, len(claims.Roles))
	for _, role := range claims.Roles {
		restRole, err := vo.NewRole(role)
		if err != nil {
			return nil, err
		}

		roles = append(roles, restRole)
	}

	return &dto.UserClaims{
		UserID: userID,
		Roles:  roles,
	}, nil
}
