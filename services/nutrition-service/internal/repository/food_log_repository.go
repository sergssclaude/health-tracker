package repository

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sergssclaude/health-tracker/nutrition-service/internal/model"
)

var ErrFoodLogNotFound = errors.New("food log not found")

type FoodLogRepository interface {
	Create(ctx context.Context, log *model.FoodLog) error
	GetByUserAndDate(ctx context.Context, userId int, date time.Time) ([]model.FoodLog, error)
	Delete(ctx context.Context, id, userId int) error
}

type postgresFoodLogRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresFoodLogRepository(pool *pgxpool.Pool) FoodLogRepository {
	return &postgresFoodLogRepository{pool: pool}
}

func (r *postgresFoodLogRepository) Create(ctx context.Context, log *model.FoodLog) error {
	query := `
	INSERT INTO food_logs(user_id, food_item_id, amount_grams, meal_type)
	VALUES ($1, $2, $3, $4)
	RETURNING id, logged_at`
	err := r.pool.QueryRow(ctx, query, log.UserID, log.FoodItemID,
		log.AmountGrams, log.MealType).Scan(&log.ID, &log.LoggedAt)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23503" {
			return ErrFoodItemNotFound
		}
		return err
	}
	return nil
}

func (r *postgresFoodLogRepository) GetByUserAndDate(ctx context.Context, userId int, date time.Time) ([]model.FoodLog, error) {
	query := `
	SELECT id, user_id, food_item_id, amount_grams, meal_type, logged_at
	FROM food_logs
	WHERE user_id = $1 AND logged_at::date = $2::date
	ORDER BY logged_at`
	rows, err := r.pool.Query(ctx, query, userId, date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var foodLogs []model.FoodLog
	for rows.Next() {
		var foodLog model.FoodLog
		if err := rows.Scan(&foodLog.ID, &foodLog.UserID, &foodLog.AmountGrams, &foodLog.MealType, &foodLog.LoggedAt); err != nil {
			return nil, err
		}
		foodLogs = append(foodLogs, foodLog)
	}
	return foodLogs, rows.Err()
}

func (r *postgresFoodLogRepository) Delete(ctx context.Context, id int, userId int) error {
	query := `
	DELETE FROM food_logs WHERE id = $1 AND user_id = $2`
	tag, err := r.pool.Exec(ctx, query, id, userId)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrFoodLogNotFound
	}
	return nil
}
