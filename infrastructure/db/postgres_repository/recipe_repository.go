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
)

type PostgresRecipeRepository struct {
	db        *pgxpool.Pool
	generator ports.IDGenerator
}

func NewPostgresRecipeRepository(db *pgxpool.Pool, generator ports.IDGenerator) repositories.RecipeRepository {
	return &PostgresRecipeRepository{
		db:        db,
		generator: generator,
	}
}

func (r *PostgresRecipeRepository) Save(ctx context.Context, recipe *domain.Recipe) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("cannot start transaction: %w", err)
	}

	defer tx.Rollback(ctx)

	recipeQuery := `INSERT INTO recipes (id, product_id) VALUES ($1, $2)`
	if _, err = tx.Exec(ctx, recipeQuery, recipe.ID(), recipe.ProductID()); err != nil {
		return fmt.Errorf("cannot insert into recipes: %w", err)
	}

	for _, ingredient := range recipe.Ingredients() {
		lineQuery := `INSERT INTO recipe_lines (id, recipe_id, item_id, quantity) VALUES ($1, $2, $3, $4)`

		lineID := r.generator.Generate()

		if _, err = tx.Exec(ctx, lineQuery, lineID, recipe.ID(), ingredient.MaterialID(), ingredient.Quantity()); err != nil {
			return fmt.Errorf("cannot insert into recipe lines: %w", err)
		}

	}

	return tx.Commit(ctx)
}

func (r *PostgresRecipeRepository) FindByProductID(ctx context.Context, productID string) (*domain.Recipe, error) {
	var recipeID string
	recipeQuery := `SELECT id FROM recipes WHERE product_id = $1`

	err := r.db.QueryRow(ctx, recipeQuery, productID).Scan(&recipeID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, &domain.ErrNotFound{ResourceName: "recipe for product", ResourceID: productID}
		}

		return nil, fmt.Errorf("cannot find recipe for product %s: %w", productID, err)
	}

	linesQuery := `SELECT id, item_id, quantity FROM recipe_lines WHERE recipe_id = $1`

	rows, err := r.db.Query(ctx, linesQuery, productID)
	if err != nil {
		return nil, fmt.Errorf("cannot find recipe for product %s: %w", productID, err)
	}

	defer rows.Close()

	var ingredients []domain.RecipeLine
	for rows.Next() {
		var lineID, materialID string
		var quantity float64

		if err = rows.Scan(&lineID, &materialID, &quantity); err != nil {
			return nil, fmt.Errorf("cannot scan line for recipe for product %s: %w", productID, err)
		}

		line, err := domain.NewRecipeLine(lineID, materialID, quantity)
		if err != nil {
			return nil, fmt.Errorf("cannot create recipe line for product %s: %w", productID, err)
		}

		ingredients = append(ingredients, *line)
	}

	recipe := domain.HydrateRecipe(recipeID, productID, ingredients)

	return recipe, nil
}

func (r *PostgresRecipeRepository) Update(ctx context.Context, recipe *domain.Recipe) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("cannot start transaction: %w", err)
	}

	defer tx.Rollback(ctx)

	deleteRecipeLinesQuery := `DELETE FROM recipe_lines WHERE recipe_id = $1`
	if _, err = tx.Exec(ctx, deleteRecipeLinesQuery, recipe.ID()); err != nil {
		return fmt.Errorf("cannot delete recipe lines: %w", err)
	}

	for _, ingredient := range recipe.Ingredients() {
		lineQuery := `INSERT INTO recipe_lines (id, recipe_id, item_id, quantity) VALUES ($1, $2, $3, $4)`

		lineID := r.generator.Generate()

		if _, err = tx.Exec(ctx, lineQuery, lineID, recipe.ID(), ingredient.MaterialID(), ingredient.Quantity()); err != nil {
			return fmt.Errorf("cannot insert into recipe lines: %w", err)
		}

	}

	return tx.Commit(ctx)
}
