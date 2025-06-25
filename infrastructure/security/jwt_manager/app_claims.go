package jwt_manager

import "github.com/golang-jwt/jwt/v5"

type AppClaims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	Name   string `json:"name"`
	jwt.RegisteredClaims
}
