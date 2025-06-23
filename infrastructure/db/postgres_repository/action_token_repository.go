package postgres_repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/muriloFlores/StoreManager/internal/core/domain"
	"github.com/muriloFlores/StoreManager/internal/core/ports"
	"github.com/muriloFlores/StoreManager/internal/core/ports/repositories"
	"time"
)

type ActionTokenRepository struct {
	db     *pgxpool.Pool
	logger ports.Logger
}

func NewActionTokenRepository(db *pgxpool.Pool, logger ports.Logger) repositories.ActionTokenRepository {
	return &ActionTokenRepository{db: db, logger: logger}
}

func (r *ActionTokenRepository) Create(ctx context.Context, token *domain.ActionToken) error {
	r.logger.InfoLevel("Invoking Create for token", map[string]interface{}{"token": token.Token, "type": token.Type})

	query := `INSERT INTO action_tokens (token, user_id, type, payload,  expires_at) VALUES ($1, $2, $3, $4, $5)`

	_, err := r.db.Exec(ctx, query,
		token.Token,
		token.UserID,
		token.Type,
		token.Payload,
		token.ExpiresAt,
	)

	if err != nil {
		r.logger.ErrorLevel("Error creating token", err, map[string]interface{}{"token": token.Token, "type": token.Type})
		return fmt.Errorf("error while reset token: %w", err)
	}

	r.logger.InfoLevel("Token created successfully", map[string]interface{}{"token": token.Token, "type": token.Type})

	return nil
}

func (r *ActionTokenRepository) FindAndConsume(ctx context.Context, tokenString string, tokenType domain.ActionType) (*domain.ActionToken, error) {
	r.logger.InfoLevel("Invoking FindAndConsume for token", map[string]interface{}{"token": tokenString, "type": tokenType})

	query := `DELETE FROM action_tokens WHERE token = $1 AND type = $2 RETURNING user_id, type, payload, expires_at`

	var userID, payload string
	var foundType domain.ActionType
	var expiresAt time.Time

	err := r.db.QueryRow(ctx, query, tokenString, tokenType).Scan(&userID, &foundType, &payload, &expiresAt)
	if err != nil {
		r.logger.ErrorLevel("Error consuming token", err, map[string]interface{}{"token": tokenString, "type": tokenType})

		if errors.Is(err, pgx.ErrNoRows) {
			r.logger.ErrorLevel("Token not found", err, map[string]interface{}{"token": tokenString, "type": tokenType})

			return nil, &domain.ErrInvalidToken{Reason: "token not found"}
		}

		return nil, fmt.Errorf("error consuming token: %w", err)
	}

	if time.Now().After(expiresAt) {
		r.logger.ErrorLevel("Token expired", nil, map[string]interface{}{"token": tokenString, "type": tokenType})
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
