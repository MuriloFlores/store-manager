package postgres_repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/muriloFlores/StoreManager/internal/core/domain"
	"github.com/muriloFlores/StoreManager/internal/core/ports"
	"github.com/muriloFlores/StoreManager/internal/core/value_objects"
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
					Details:  fmt.Sprintf("data conflict: %s", pgErr.ConstraintName),
				}
			}
		}

		return fmt.Errorf("error saving user: %w", err)
	}

	return nil
}

func (p *PostgresUserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	var id, name, emailScan, password string
	var role value_objects.Role

	query := `SELECT id, name, email, password_hash, role FROM users WHERE email = $1 AND deleted_at IS NULL`

	err := p.db.QueryRow(ctx, query, email).Scan(&id, &name, &emailScan, &password, &role)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, &domain.ErrNotFound{ResourceName: "user", ResourceID: email}
		}

		return nil, fmt.Errorf("error finding user by email: %w", err)
	}

	return domain.NewUser(id, name, emailScan, password, role)
}

func (p *PostgresUserRepository) FindByID(ctx context.Context, id string) (*domain.User, error) {
	var idScan, name, email, passwordHash string
	var role value_objects.Role

	query := `SELECT id, name, email, password_hash, role FROM users WHERE id = $1 AND deleted_at IS NULL`

	err := p.db.QueryRow(ctx, query, id).Scan(&idScan, &name, &email, &passwordHash, &role)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, &domain.ErrNotFound{ResourceName: "user", ResourceID: id}
		}

		return nil, fmt.Errorf("error finding user by id: %w", err)
	}

	return domain.NewUser(idScan, name, email, passwordHash, role)
}

func (p *PostgresUserRepository) Update(ctx context.Context, user *domain.User) error {
	query := `UPDATE users SET name = $1, email = $2, password_hash = $3,  role = $4 WHERE id = $5`

	commandTag, err := p.db.Exec(ctx, query,
		user.Name(),
		user.Email(),
		user.Password(),
		user.Role(),
		user.ID(),
	)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return &domain.ErrConflict{
				Resource: "user",
				Details:  fmt.Sprintf("data conflict: %s", pgErr.ConstraintName),
			}
		}

		return fmt.Errorf("error updating user: %w", err)
	}

	if commandTag.RowsAffected() == 0 {
		return &domain.ErrNotFound{ResourceName: "user", ResourceID: user.ID()}
	}

	return nil
}

func (p *PostgresUserRepository) Delete(ctx context.Context, userID string) error {
	query := `UPDATE users SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL`

	commandTag, err := p.db.Exec(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("error deleting user: %w", err)
	}

	if commandTag.RowsAffected() == 0 {
		return &domain.ErrNotFound{ResourceName: "user", ResourceID: userID}
	}

	return nil
}

func (p *PostgresUserRepository) CountAdmins(ctx context.Context) (int, error) {
	var count int

	query := `SELECT COUNT(*) FROM users WHERE role = 'admin' AND deleted_at IS NULL`

	err := p.db.QueryRow(ctx, query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("error counting admins: %w", err)
	}

	return count, nil
}
