package postgres_repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/muriloFlores/StoreManager/internal/core/domain"
	"github.com/muriloFlores/StoreManager/internal/core/domain/item"
	"github.com/muriloFlores/StoreManager/internal/core/domain/pagination"
	"github.com/muriloFlores/StoreManager/internal/core/ports"
	"github.com/muriloFlores/StoreManager/internal/core/ports/repositories"
	"time"
)

type PostgresItemRepository struct {
	db     *pgxpool.Pool
	logger ports.Logger
}

func NewPostgresItemRepository(db *pgxpool.Pool, logger ports.Logger) repositories.ItemRepository {
	return &PostgresItemRepository{
		db:     db,
		logger: logger,
	}
}

func (r *PostgresItemRepository) Save(ctx context.Context, item *item.Item) error {
	query := `
		INSERT INTO items (id, name, description, sku, type, can_be_sold, price_sale_in_cents, stock_quantity, minimum_stock_level, unit_of_measure, is_active)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`

	_, err := r.db.Exec(ctx, query,
		item.ID(),
		item.Name(),
		item.Description(),
		item.SKU(),
		item.ItemType(),
		item.CanBeSold(),
		item.PriceInCents(),
		item.StockQuantity(),
		item.MinimumStockLevel(),
		item.UnitOfMeasure(),
		item.IsActive(),
	)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return &domain.ErrConflict{Resource: "Item", Details: fmt.Sprintf("constraint violation: %s", pgErr.ConstraintName)}
		}
		r.logger.ErrorLevel("Error saving item", err, map[string]interface{}{"item_id": item.ID()})
		return fmt.Errorf("error saving item: %w", err)
	}
	return nil
}

func (r *PostgresItemRepository) FindByID(ctx context.Context, id string) (*item.Item, error) {
	query := `
	SELECT id, name, description, sku, type, is_active, can_be_sold, price_sale_in_cents, stock_quantity,
		unit_of_measure, minimum_stock_level, deleted_at
	FROM items WHERE id = $1 AND deleted_at IS NULL
`
	var deletedAt *time.Time
	var params item.HydrateItemParams

	err := r.db.QueryRow(ctx, query, id).Scan(
		&params.Id, &params.Name, &params.Description, &params.Sku, &params.ItemType, &params.Active, &params.CanBeSold,
		&params.PriceInCents, &params.StockQuantity, &params.UnitOfMeasure,
		&params.MinimumStockLevel, &deletedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, &domain.ErrNotFound{ResourceName: "item", ResourceID: id}
		}

		r.logger.ErrorLevel("Error finding item", err, map[string]interface{}{"item": id})
		return nil, err
	}

	item := item.HydrateItem(params, deletedAt)

	return item, err
}

func (r *PostgresItemRepository) FindBySKU(ctx context.Context, sku string) (*item.Item, error) {
	query := `
	SELECT id, name, description, sku, type, is_active, can_be_sold, price_sale_in_cents, stock_quantity,
		unit_of_measure, minimum_stock_level, deleted_at
	FROM items WHERE sku = $1 AND deleted_at IS NULL
`
	var deletedAt *time.Time
	var params item.HydrateItemParams

	err := r.db.QueryRow(ctx, query, sku).Scan(
		&params.Id, &params.Name, &params.Description, &params.Sku, &params.ItemType, &params.Active, &params.CanBeSold,
		&params.PriceInCents, &params.StockQuantity, &params.UnitOfMeasure,
		&params.MinimumStockLevel, &deletedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, &domain.ErrNotFound{ResourceName: "item", ResourceID: sku}
		}

		r.logger.ErrorLevel("Error finding item", err, map[string]interface{}{"item": sku})
		return nil, err
	}

	item := item.HydrateItem(params, deletedAt)

	return item, err
}

func (r *PostgresItemRepository) Update(ctx context.Context, item *item.Item) error {
	query := `
		UPDATE items SET
			name = $1, description = $2, sku = $3, is_active = $4, can_be_sold = $5,
			price_sale_in_cents = $6, stock_quantity = $7, unit_of_measure = $8,
			minimum_stock_level = $9, deleted_at = $10
		WHERE id = $11`

	commandTag, err := r.db.Exec(ctx, query,
		item.Name(),
		item.Description(),
		item.SKU(),
		item.IsActive(),
		item.CanBeSold(),
		item.PriceInCents(),
		item.StockQuantity(),
		item.UnitOfMeasure(),
		item.MinimumStockLevel(),
		item.DeletedAt(),
		item.ID(),
	)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return &domain.ErrConflict{
				Resource: "Item",
				Details:  fmt.Sprintf("constraint violation: %s", pgErr.ConstraintName),
			}
		}
		return fmt.Errorf("error updating user: %w", err)
	}

	if commandTag.RowsAffected() == 0 {
		return &domain.ErrNotFound{ResourceName: "item", ResourceID: item.ID()}

	}

	return nil
}

func (r *PostgresItemRepository) List(ctx context.Context, PaginationParams *pagination.PaginationParams) (*pagination.PaginatedResult[*item.Item], error) {
	var totalItems int64
	countQuery := `SELECT COUNT(*) FROM items WHERE deleted_at IS NULL`

	if err := r.db.QueryRow(ctx, countQuery).Scan(&totalItems); err != nil {
		return nil, fmt.Errorf("error counting items: %w", err)
	}

	paginationInfo := pagination.PaginationInfo{ /* ... */ }
	paginationInfo.CalculateTotalPages()
	offset := (PaginationParams.Page - 1) * PaginationParams.PageSize

	query := `
		SELECT id, name, description, sku, type, is_active, can_be_sold, price_sale_in_cents,
		       stock_quantity, unit_of_measure, minimum_stock_level, created_at, updated_at, deleted_at
		FROM items
		WHERE deleted_at IS NULL
		ORDER BY name
		LIMIT $1 OFFSET $2`

	rows, err := r.db.Query(ctx, query, PaginationParams.PageSize, offset)
	if err != nil {
		return nil, fmt.Errorf("error listing items: %w", err)
	}
	defer rows.Close()

	items := make([]*item.Item, 0, PaginationParams.PageSize)
	for rows.Next() {
		var params item.HydrateItemParams
		var createdAt, updatedAt, deletedAt *time.Time

		if err := rows.Scan(
			&params.Id, &params.Name, &params.Description, &params.Sku, &params.ItemType,
			&params.Active, &params.CanBeSold, &params.PriceInCents, &params.StockQuantity,
			&params.UnitOfMeasure, &params.MinimumStockLevel,
			&createdAt, &updatedAt, &deletedAt,
		); err != nil {
			return nil, fmt.Errorf("error scanning item row: %w", err)
		}

		item := item.HydrateItem(params, deletedAt)
		items = append(items, item)
	}

	return &pagination.PaginatedResult[*item.Item]{
		Data:       items,
		Pagination: paginationInfo,
	}, nil
}
