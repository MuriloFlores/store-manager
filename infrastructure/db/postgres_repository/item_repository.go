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
	"strings"
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
		INSERT INTO items (id, name, description, sku, type, can_be_sold, price_sale_in_cents, price_cost_in_cents, stock_quantity, minimum_stock_level, unit_of_measure, is_active)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`

	_, err := r.db.Exec(ctx, query,
		item.ID(),
		item.Name(),
		item.Description(),
		item.SKU(),
		item.ItemType(),
		item.CanBeSold(),
		item.PriceInCents(),
		item.PriceCostInCents(),
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
	SELECT id, name, description, sku, type, is_active, can_be_sold, price_sale_in_cents, price_cost_in_cents, stock_quantity,
		unit_of_measure, minimum_stock_level, deleted_at
	FROM items WHERE id = $1 AND deleted_at IS NULL
`
	var deletedAt *time.Time
	var params item.HydrateItemParams

	err := r.db.QueryRow(ctx, query, id).Scan(
		&params.Id, &params.Name, &params.Description, &params.Sku, &params.ItemType, &params.Active, &params.CanBeSold,
		&params.PriceInCents, &params.PriceCostInCents, &params.StockQuantity, &params.UnitOfMeasure,
		&params.MinimumStockLevel, &deletedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, &domain.ErrNotFound{ResourceName: "item", ResourceID: id}
		}

		r.logger.ErrorLevel("Error finding item", err, map[string]interface{}{"item": id})
		return nil, err
	}

	itemDomain := item.HydrateItem(params, deletedAt)

	return itemDomain, err
}

func (r *PostgresItemRepository) FindByIDIncludingDeleted(ctx context.Context, id string) (*item.Item, error) {
	query := `
	SELECT id, name, description, sku, type, is_active, can_be_sold, price_sale_in_cents, price_cost_in_cents, stock_quantity,
		unit_of_measure, minimum_stock_level, deleted_at
	FROM items WHERE id = $1`

	var deletedAt *time.Time
	var params item.HydrateItemParams

	err := r.db.QueryRow(ctx, query, id).Scan(
		&params.Id, &params.Name, &params.Description, &params.Sku, &params.ItemType, &params.Active, &params.CanBeSold,
		&params.PriceInCents, &params.PriceCostInCents, &params.StockQuantity, &params.UnitOfMeasure,
		&params.MinimumStockLevel, &deletedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, &domain.ErrNotFound{ResourceName: "item", ResourceID: id}
		}

		r.logger.ErrorLevel("Error finding item", err, map[string]interface{}{"item": id})
		return nil, err
	}

	itemDomain := item.HydrateItem(params, deletedAt)

	return itemDomain, err
}

func (r *PostgresItemRepository) FindBySKUIncludingDeleted(ctx context.Context, sku string) (*item.Item, error) {
	query := `
	SELECT id, name, description, sku, type, is_active, can_be_sold, price_sale_in_cents, price_cost_in_cents, stock_quantity,
		unit_of_measure, minimum_stock_level, deleted_at
	FROM items WHERE sku = $1`

	var deletedAt *time.Time
	var params item.HydrateItemParams

	err := r.db.QueryRow(ctx, query, sku).Scan(
		&params.Id, &params.Name, &params.Description, &params.Sku, &params.ItemType, &params.Active, &params.CanBeSold,
		&params.PriceInCents, &params.PriceCostInCents, &params.StockQuantity, &params.UnitOfMeasure,
		&params.MinimumStockLevel, &deletedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, &domain.ErrNotFound{ResourceName: "item", ResourceID: sku}
		}

		r.logger.ErrorLevel("Error finding item", err, map[string]interface{}{"item": sku})
		return nil, err
	}

	itemDomain := item.HydrateItem(params, deletedAt)

	return itemDomain, err
}

func (r *PostgresItemRepository) FindBySKU(ctx context.Context, sku string) (*item.Item, error) {
	query := `
	SELECT id, name, description, sku, type, is_active, can_be_sold, price_sale_in_cents, price_cost_in_cents, stock_quantity,
		unit_of_measure, minimum_stock_level, deleted_at
	FROM items WHERE sku = $1 AND deleted_at IS NULL
`
	var deletedAt *time.Time
	var params item.HydrateItemParams

	err := r.db.QueryRow(ctx, query, sku).Scan(
		&params.Id, &params.Name, &params.Description, &params.Sku, &params.ItemType, &params.Active, &params.CanBeSold,
		&params.PriceInCents, &params.PriceCostInCents, &params.StockQuantity, &params.UnitOfMeasure,
		&params.MinimumStockLevel, &deletedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, &domain.ErrNotFound{ResourceName: "item", ResourceID: sku}
		}

		r.logger.ErrorLevel("Error finding item", err, map[string]interface{}{"item": sku})
		return nil, err
	}

	itemDomain := item.HydrateItem(params, deletedAt)

	return itemDomain, err
}

func (r *PostgresItemRepository) Update(ctx context.Context, item *item.Item) error {
	r.logger.InfoLevel("Updating item in repository", map[string]interface{}{"item": item.ID()})

	query := `
		UPDATE items SET 
			name = $1,
			description = $2,
			is_active = $3,
			can_be_sold = $4,
			price_sale_in_cents = $5,
			price_cost_in_cents = $6,
			minimum_stock_level = $7, 
			deleted_at = $8 
		WHERE id = $9`

	commandTag, err := r.db.Exec(ctx, query,
		item.Name(),
		item.Description(),
		item.IsActive(),
		item.CanBeSold(),
		item.PriceInCents(),
		item.PriceCostInCents(),
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
		r.logger.ErrorLevel("Failed to update item", err, map[string]interface{}{"item_id": item.ID()})
		return fmt.Errorf("erro ao atualizar item: %w", err)
	}

	if commandTag.RowsAffected() == 0 {
		return &domain.ErrNotFound{ResourceName: "item", ResourceID: item.ID()}
	}

	return nil
}

func (r *PostgresItemRepository) List(ctx context.Context, paginationParams *pagination.PaginationParams) (*pagination.PaginatedResult[*item.Item], error) {
	r.logger.InfoLevel("Init list item repository")

	var totalItems int64
	countQuery := `SELECT COUNT(*) FROM items WHERE deleted_at IS NULL`

	if err := r.db.QueryRow(ctx, countQuery).Scan(&totalItems); err != nil {
		r.logger.ErrorLevel("error while scan total items", err, map[string]interface{}{})
		return nil, fmt.Errorf("error counting items: %w", err)
	}

	paginationInfo := pagination.PaginationInfo{
		CurrentPage: paginationParams.Page,
		PageSize:    paginationParams.PageSize,
		TotalItems:  totalItems,
	}

	paginationInfo.CalculateTotalPages()
	r.logger.InfoLevel("pagination info complete")

	offset := (paginationParams.Page - 1) * paginationParams.PageSize

	query := `
		SELECT id, name, description, sku, type, is_active, can_be_sold, price_sale_in_cents, price_cost_in_cents,
		       stock_quantity, unit_of_measure, minimum_stock_level, created_at, updated_at, deleted_at
		FROM items
		WHERE deleted_at IS NULL
		ORDER BY name
		LIMIT $1 OFFSET $2`

	rows, err := r.db.Query(ctx, query, paginationParams.PageSize, offset)
	if err != nil {
		r.logger.ErrorLevel("error while execute query", err, map[string]interface{}{})
		return nil, fmt.Errorf("error listing items: %w", err)
	}
	defer rows.Close()

	items := make([]*item.Item, 0, paginationParams.PageSize)
	for rows.Next() {
		var params item.HydrateItemParams
		var createdAt, updatedAt, deletedAt *time.Time

		if err = rows.Scan(
			&params.Id, &params.Name, &params.Description, &params.Sku, &params.ItemType,
			&params.Active, &params.CanBeSold, &params.PriceInCents, &params.PriceCostInCents, &params.StockQuantity,
			&params.UnitOfMeasure, &params.MinimumStockLevel,
			&createdAt, &updatedAt, &deletedAt,
		); err != nil {
			r.logger.ErrorLevel("error while scan item", err, map[string]interface{}{})
			return nil, fmt.Errorf("error scanning item row: %w", err)
		}

		itemDomain := item.HydrateItem(params, deletedAt)
		items = append(items, itemDomain)
	}

	return &pagination.PaginatedResult[*item.Item]{
		Data:       items,
		Pagination: paginationInfo,
	}, nil
}

func (r *PostgresItemRepository) Delete(ctx context.Context, itemID string) error {
	query := `UPDATE items SET deleted_at = NOW() WHERE id = $1`

	commandTag, err := r.db.Exec(ctx, query, itemID)
	if err != nil {
		return fmt.Errorf("error deleting item: %w", err)
	}

	if commandTag.RowsAffected() == 0 {
		return &domain.ErrNotFound{ResourceName: "item", ResourceID: itemID}
	}

	return nil
}

func (r *PostgresItemRepository) ListForUsers(ctx context.Context, paginationParams *pagination.PaginationParams) (*pagination.PaginatedResult[*item.Item], error) {
	var totalItems int64
	countQuery := `SELECT COUNT(*) FROM items WHERE deleted_at IS NULL AND is_active = true AND can_be_sold = true`

	if err := r.db.QueryRow(ctx, countQuery).Scan(&totalItems); err != nil {
		return nil, fmt.Errorf("error counting items: %w", err)
	}

	paginationInfo := pagination.PaginationInfo{
		CurrentPage: paginationParams.Page,
		PageSize:    paginationParams.PageSize,
		TotalItems:  totalItems,
	}

	paginationInfo.CalculateTotalPages()
	offset := (paginationParams.Page - 1) * paginationParams.PageSize

	query := `
		SELECT id, name, description, sku, type, is_active, can_be_sold, price_sale_in_cents, price_cost_in_cents,
		       stock_quantity, unit_of_measure, minimum_stock_level, created_at, updated_at, deleted_at
		FROM items
		WHERE deleted_at IS NULL
		ORDER BY name
		LIMIT $1 OFFSET $2`

	rows, err := r.db.Query(ctx, query, paginationParams.PageSize, offset)
	if err != nil {
		return nil, fmt.Errorf("error listing items: %w", err)
	}
	defer rows.Close()

	items := make([]*item.Item, 0, paginationParams.PageSize)
	for rows.Next() {
		var params item.HydrateItemParams
		var createdAt, updatedAt, deletedAt *time.Time

		if err := rows.Scan(
			&params.Id, &params.Name, &params.Description, &params.Sku, &params.ItemType,
			&params.Active, &params.CanBeSold, &params.PriceInCents, &params.PriceCostInCents, &params.StockQuantity,
			&params.UnitOfMeasure, &params.MinimumStockLevel,
			&createdAt, &updatedAt, &deletedAt,
		); err != nil {
			return nil, fmt.Errorf("error scanning item row: %w", err)
		}

		itemDomain := item.HydrateItem(params, deletedAt)
		items = append(items, itemDomain)
	}

	return &pagination.PaginatedResult[*item.Item]{
		Data:       items,
		Pagination: paginationInfo,
	}, nil
}

func (r *PostgresItemRepository) Search(ctx context.Context, searchTerm string, isPublicSearch bool, params *pagination.PaginationParams) (*pagination.PaginatedResult[*item.Item], error) {
	var queryBuilder strings.Builder

	args := make([]interface{}, 0)

	queryBuilder.WriteString("WHERE (name LIKE $1 OR sku LIKE $2) ")
	args = append(args, fmt.Sprintf("%%%s%%", searchTerm))

	if isPublicSearch {
		queryBuilder.WriteString("AND is_active = true AND can_be_sold = true ")
	}

	queryBuilder.WriteString("AND deleted_at IS NULL ")

	queryClause := queryBuilder.String()

	var totalItems int64
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM items %s", queryClause)

	if err := r.db.QueryRow(ctx, countQuery, args...).Scan(&totalItems); err != nil {
		r.logger.ErrorLevel("Error counting items for search", err, map[string]interface{}{"search_term": searchTerm})
		return nil, fmt.Errorf("error counting items for search: %w", err)
	}

	paginationInfo := pagination.PaginationInfo{
		CurrentPage: params.Page,
		PageSize:    params.PageSize,
		TotalItems:  totalItems,
	}

	paginationInfo.CalculateTotalPages()
	offset := (params.Page - 1) * params.PageSize

	args = append(args, params.PageSize, offset)

	query := fmt.Sprintf(`
		SELECT id, name, description, sku, type, is_active, can_be_sold, price_sale_in_cents,
		       stock_quantity, unit_of_measure, minimum_stock_level, created_at, updated_at, deleted_at
		FROM items
		%s
		ORDER BY name ASC
		LIMIT $%d OFFSET $%d`,
		queryClause, len(args)-1, len(args))

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		r.logger.ErrorLevel("error executing search query for items", err, map[string]interface{}{"search_term": searchTerm})
		return nil, fmt.Errorf("error executing search query for items: %w", err)
	}

	defer rows.Close()

	items := make([]*item.Item, 0, params.PageSize)
	for rows.Next() {
		var itemParams item.HydrateItemParams
		var createdAt, updatedAt, deletedAt *time.Time

		if err = rows.Scan(
			&itemParams.Id, &itemParams.Name, &itemParams.Description, &itemParams.Sku, &itemParams.ItemType,
			&itemParams.Active, &itemParams.CanBeSold, &itemParams.PriceInCents,
			&itemParams.StockQuantity, &itemParams.UnitOfMeasure, &itemParams.MinimumStockLevel,
			&createdAt, &updatedAt, &deletedAt,
		); err != nil {
			r.logger.ErrorLevel("error scanning item row during search", err, map[string]interface{}{"search_term": searchTerm})
			return nil, fmt.Errorf("error scanning item row during search: %w", err)
		}

		itemDomain := item.HydrateItem(itemParams, deletedAt)
		items = append(items, itemDomain)
	}

	if err = rows.Err(); err != nil {
		r.logger.ErrorLevel("error after iterating rows", err, map[string]interface{}{"search_term": searchTerm})
		return nil, fmt.Errorf("error iterating items: %w", err)
	}

	return &pagination.PaginatedResult[*item.Item]{
		Data:       items,
		Pagination: paginationInfo,
	}, nil
}
