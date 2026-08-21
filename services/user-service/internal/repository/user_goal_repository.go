package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sergssclaude/health-tracker/user-service/internal/model"
)

type UserGoalRepository interface {
	Upsert(ctx context.Context, model *model.UserGoal) error
	GetByUserId(ctx context.Context, userId int) (*model.UserGoal, error)
}

type postgresUserGoalRepository struct {
	pool *pgxpool.Pool
}

func NewUserGoalRepository(pool *pgxpool.Pool) UserGoalRepository {
	return &postgresUserGoalRepository{pool: pool}
}

func (r *postgresUserGoalRepository) Upsert(ctx context.Context, goal *model.UserGoal) error {
	query := `
        INSERT INTO user_goal (user_id, target_weight, calorie_goal, updated_at)
        VALUES ($1, $2, $3, now())
        ON CONFLICT (user_id) DO UPDATE SET
            target_weight = COALESCE(EXCLUDED.target_weight, user_goal.target_weight),
            calorie_goal = COALESCE(EXCLUDED.calorie_goal, user_goal.calorie_goal),
            updated_at = now()
    `
	_, err := r.pool.Exec(ctx, query, goal.UserId, goal.TargetWeight, goal.CalorieGoal)
	return err
}
func (r *postgresUserGoalRepository) GetByUserId(ctx context.Context, userId int) (*model.UserGoal, error) {
	query := `SELECT user_id, target_weight, calorie_goal, created_at, updated_at
	FROM user_goal WHERE user_id = $1`
	var userGoal model.UserGoal
	err := r.pool.QueryRow(ctx, query, userId).Scan(
		&userGoal.UserId,
		&userGoal.CalorieGoal,
		&userGoal.TargetWeight,
		&userGoal.CreateAt,
		&userGoal.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
	}
	return &userGoal, nil
}
