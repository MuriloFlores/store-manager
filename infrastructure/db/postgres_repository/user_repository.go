package postgres_repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/muriloFlores/StoreManager/internal/core/domain"
	"github.com/muriloFlores/StoreManager/internal/core/ports/repositories"
	"github.com/muriloFlores/StoreManager/internal/core/value_objects"
	"time"
)

type PostgresUserRepository struct {
	db *pgxpool.Pool
}

func NewPostgresUserRepository(db *pgxpool.Pool) repositories.UserRepository {
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
	var id, name, emailScan, passwordHash string
	var role value_objects.Role
	var verifiedAt *time.Time

	query := `SELECT id, name, email, password_hash, role, verified_at FROM users WHERE email = $1 AND deleted_at IS NULL`

	err := p.db.QueryRow(ctx, query, email).Scan(&id, &name, &emailScan, &passwordHash, &role, &verifiedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, &domain.ErrNotFound{ResourceName: "user", ResourceID: email}
		}

		return nil, fmt.Errorf("error finding user by email: %w", err)
	}

	user := domain.HydrateUser(id, name, emailScan, passwordHash, role, verifiedAt, nil)

	return user, nil
}

func (p *PostgresUserRepository) FindByID(ctx context.Context, id string) (*domain.User, error) {
	var idScan, name, email, passwordHash string
	var role value_objects.Role
	var verifiedAt *time.Time

	query := `SELECT id, name, email, password_hash, role, verified_at FROM users WHERE id = $1 AND deleted_at IS NULL`

	err := p.db.QueryRow(ctx, query, id).Scan(&idScan, &name, &email, &passwordHash, &role, &verifiedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, &domain.ErrNotFound{ResourceName: "user", ResourceID: id}
		}

		return nil, fmt.Errorf("error finding user by id: %w", err)
	}

	user := domain.HydrateUser(idScan, name, email, passwordHash, role, verifiedAt, nil)

	return user, nil
}

func (p *PostgresUserRepository) Update(ctx context.Context, user *domain.User) error {
	query := `
        UPDATE users SET 
            name = $1, 
            email = $2, 
            password_hash = $3,  
            role = $4, 
            verified_at = $5,
            deleted_at = $6
        WHERE id = $7`

	commandTag, err := p.db.Exec(ctx, query,
		user.Name(),
		user.Email(),
		user.Password(),
		user.Role(),
		user.VerifiedAt(),
		user.DeletedAt(),
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
	query := `UPDATE users SET deleted_at = NOW(), verified_at = null WHERE id = $1 AND deleted_at IS NULL`

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

func (p *PostgresUserRepository) FindByEmailIncludingDeleted(ctx context.Context, email string) (*domain.User, error) {
	var id, name, passwordHash string
	var role value_objects.Role
	var verifiedAt, deletedAt *time.Time

	query := `SELECT id, name, email, password_hash, role, verified_at, deleted_at FROM users WHERE email = $1`

	err := p.db.QueryRow(ctx, query, email).Scan(
		&id, &name, &email, &passwordHash, &role, &verifiedAt, &deletedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, &domain.ErrNotFound{ResourceName: "user", ResourceID: email}
		}
		return nil, fmt.Errorf("error finding user by email including deleted: %w", err)
	}

	user := domain.HydrateUser(id, name, email, passwordHash, role, verifiedAt, deletedAt)

	return user, nil
}

func (p *PostgresUserRepository) List(ctx context.Context, params *domain.PaginationParams) (*domain.PaginatedUsers, error) {
	var totalItems int

	countQuery := `SELECT COUNT(*) FROM users WHERE deleted_at IS NULL AND verified_at IS NOT NULL`

	err := p.db.QueryRow(ctx, countQuery).Scan(&totalItems)
	if err != nil {
		return nil, fmt.Errorf("error counting users: %w", err)
	}

	paginationInfo := &domain.PaginationInfo{
		CurrentPage: params.Page,
		PageSize:    params.PageSize,
		TotalItems:  int64(totalItems),
	}

	paginationInfo.CalculateTotalPages()

	offset := (params.Page - 1) * params.PageSize

	query := `
        SELECT id, name, email, password_hash, role, verified_at, deleted_at 
        FROM users 
        WHERE deleted_at IS NULL 
        ORDER BY created_at DESC 
        LIMIT $1 OFFSET $2
    `

	rows, err := p.db.Query(ctx, query, params.Page*params.PageSize, offset)
	if err != nil {
		return nil, fmt.Errorf("error listing users: %w", err)
	}

	defer rows.Close()

	users := make([]*domain.User, 0)
	for rows.Next() {
		var id, name, email, passwordHash string
		var role value_objects.Role
		var verifiedAt, deletedAt *time.Time

		if err := rows.Scan(&id, &name, &email, &passwordHash, &role, &verifiedAt, &deletedAt); err != nil {
			return nil, fmt.Errorf("error scanning user row: %w", err)
		}

		user := domain.HydrateUser(id, name, email, passwordHash, role, verifiedAt, deletedAt)
		users = append(users, user)
	}

	return &domain.PaginatedUsers{
		Data:       users,
		Pagination: *paginationInfo,
	}, nil
}
