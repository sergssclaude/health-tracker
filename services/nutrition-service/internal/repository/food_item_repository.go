package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sergssclaude/health-tracker/nutrition-service/internal/model"
)

type FoodItemRepository interface {
	Create(ctx context.Context, foodItem *model.FoodItem) error
	GetById(ctx context.Context, id int) (*model.FoodItem, error)
	Search(ctx context.Context, nameItem string) ([]model.FoodItem, error)
}

type postgresFoodItemRepository struct {
	pool *pgxpool.Pool
}

func newPostgresFoodItemRepository(pool *pgxpool.Pool) FoodItemRepository {
	return &postgresFoodItemRepository{pool: pool}
}

var ErrFoodItemNotFound = errors.New("Food item not found")

func (r *postgresFoodItemRepository) Create(ctx context.Context, foodItem *model.FoodItem) error {
	query := `
	INSERT INTO food_items(name, calories_per_100g, protein_per_100g, fats_per_100g, carbs_per_100g)
	VALUES ($1, $2, $3, $4, $5)
	`
	_, err := r.pool.Exec(ctx, query, foodItem.Name, foodItem.CaloriesPer100g, foodItem.ProteinPer100g, foodItem.FatsPer100g, foodItem.CarbsPer100g)
	if err != nil {
		return err
	}
	return nil
}

func (r *postgresFoodItemRepository) GetById(ctx context.Context, id int) (*model.FoodItem, error) {
	query := `
	SELECT (fi.id, fi.name, fi.calories_per_100g,  fi.protein_per_100g,  fi.fats_per_100g, fi.carbs_per_100g)
	FROM food_items fi
	WHERE fi.id = $1
	`
	var foodItem model.FoodItem
	err := r.pool.QueryRow(ctx, query, id).Scan(&foodItem)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrFoodItemNotFound
		}
		return nil, err
	}
	return &foodItem, nil
}

func (r *postgresFoodItemRepository) Search(ctx context.Context, nameItem string) ([]model.FoodItem, error) {
	query := `
	SELECT (id, name, calories_per_100g, protein_per_100g, fats_per_100g, carbs_per_100g)
	FROM food_items
	WHERE name ILIKE '%' || $1 || '%'
	`
	rows, err := r.pool.Query(ctx, query, nameItem)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []model.FoodItem

	for rows.Next() {
		var item model.FoodItem
		if err := rows.Scan(&item.ID, &item.Name, &item.CaloriesPer100g, &item.ProteinPer100g, &item.FatsPer100g, &item.CarbsPer100g); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return nil, rows.Err()
}
