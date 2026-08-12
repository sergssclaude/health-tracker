package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sergssclaude/health-tracker/user-service/internal/model"
)

type UserProfileRepository interface {
	GetFullProfile(ctx context.Context, userId int) (*model.UserProfile, error)
}

type postgresUserProfileRepository struct {
	pool *pgxpool.Pool
}

func NewUserProfileRepository(pool *pgxpool.Pool) UserProfileRepository {
	return postgresUserProfileRepository{pool: pool}
}

func (r *postgresUserProfileRepository) GetFullProfile(ctx context.Context, userId int) (*model.UserProfile, error) {
	query := `
	SELECT u.id, u.email, u.name, u.created_at, u.updated_at ,
	ui.weight, ui.height, ui.age, ui.gender, ui.daily_calorie_norm, ui.updated_at,
	ug.calorie_goal , ug.target_weight, ug.updated_at
	from users u
	LEFT JOIN user_information ui on ui.user_id = u.id 
	LEFT JOIN user_goal ug on ug.user_id = u.id
	WHERE u.id = $1
	`
	var user model.User
	var info model.UserInformation
	var goal model.UserGoal
	r.pool.QueryRow(ctx, query, userId).Scan(
		&user.Id, &user.Email, &user.Name, &user.CreatedAt, &user.UpdatedAt,
		&info.Weight, &info.Height, &info.Age, &info.Gender, &info.DailyCalorieNorm, &info.UpdatedAt,
		&goal.CalorieGoal, &goal.TargetWeight, &goal.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
	}
}
