package repository

import (
	"context"

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

func (r *postgresUserGoalRepository) Upsert(ctx context.Context, model *model.UserGoal) error {
	//TODO
	return nil
}

func (r *postgresUserGoalRepository) GetByUserId(ctx context.Context, userId int) (*model.UserGoal, error) {
	//TODO
	return nil, nil
}
