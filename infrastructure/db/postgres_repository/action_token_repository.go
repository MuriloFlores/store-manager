package postgres_repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/muriloFlores/StoreManager/internal/core/domain"
	"github.com/muriloFlores/StoreManager/internal/core/ports"
	"time"
)

type ActionTokenRepository struct {
	db *pgxpool.Pool
}

func NewActionTokenRepository(db *pgxpool.Pool) ports.ActionTokenRepository {
	return &ActionTokenRepository{db: db}
}

func (r *ActionTokenRepository) Create(ctx context.Context, token *domain.ActionToken) error {
	query := `INSERT INTO action_tokens (token, user_id, type, payload,  expires_at) VALUES ($1, $2, $3, $4, $5)`

	_, err := r.db.Exec(ctx, query,
		token.Token,
		token.UserID,
		token.Type,
		token.Payload,
		token.ExpiresAt,
	)

	if err != nil {
		return fmt.Errorf("error while reset token: %w", err)
	}

	return nil
}

func (r *ActionTokenRepository) FindAndConsume(ctx context.Context, tokenString string, tokenType domain.ActionType) (*domain.ActionToken, error) {
	query := `DELETE FROM action_tokens WHERE token = $1 AND type = $2 RETURNING user_id, type, payload, expires_at`

	var userID, payload string
	var foundType domain.ActionType
	var expiresAt time.Time

	err := r.db.QueryRow(ctx, query, tokenString).Scan(&userID, &foundType, &payload, &expiresAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, &domain.ErrInvalidToken{Reason: "token not found"}
		}

		return nil, fmt.Errorf("error consuming token: %w", err)
	}

	if time.Now().After(expiresAt) {
		return nil, &domain.ErrInvalidToken{Reason: "token expired"}
	}

	return &domain.ActionToken{
		Token:     tokenString,
		UserID:    userID,
		Type:      foundType,
		Payload:   payload,
		ExpiresAt: expiresAt,
	}, nil
}
