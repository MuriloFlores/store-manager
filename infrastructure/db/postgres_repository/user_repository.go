package postgres_repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/muriloFlores/StoreManager/internal/core/domain"
	"github.com/muriloFlores/StoreManager/internal/core/ports"
)

type PostgresUserRepository struct {
	db *pgxpool.Pool
}

func NewPostgresUserRepository(db *pgxpool.Pool) ports.UserRepository {
	return &PostgresUserRepository{db: db}
}

func (p *PostgresUserRepository) Save(ctx context.Context, user *domain.User) error {
	query := `INSERT INTO users (id, name, email, password_hash, role) VALUES ($1, $2, $3, $4, $5)`

	_, err := p.db.Exec(ctx, query,
		user.ID(),
		user.Name(),
		user.Email(),
		user.Password(),
		user.Role(),
	)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return &domain.ErrConflict{
					Resource: "user",
					Details:  fmt.Sprintf("database error: %s", pgErr.ConstraintName),
				}
			}
		}

		return fmt.Errorf("error saving user: %w", err)
	}

	return nil
}

func (p *PostgresUserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	//TODO implement me
	panic("implement me")
}

func (p *PostgresUserRepository) Update(ctx context.Context, user *domain.User) error {
	//TODO implement me
	panic("implement me")
}

func (p *PostgresUserRepository) Delete(ctx context.Context, userID string) error {
	//TODO implement me
	panic("implement me")
}
