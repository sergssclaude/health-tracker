package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sergssclaude/health-tracker/user-service/internal/model"
)

type UserInformationRepository interface {
	Upsert(ctx context.Context, info *model.UserInformation) error
	GetByUserId(ctx context.Context, userId int) (*model.UserInformation, error)
}

type postgresUserInformationRepository struct {
	pool *pgxpool.Pool
}

func NewUserInformationRepository(pool *pgxpool.Pool) UserInformationRepository {
	return &postgresUserInformationRepository{pool: pool}
}

func (r *postgresUserInformationRepository) Upsert(ctx context.Context, info *model.UserInformation) error {
	query := `
	INSERT INTO user_information (user_id, weight, height, age, gender, daily_calorie_norm, profile_complited, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7, now())
	ON CONFLICT (user_id) DO UPDATE SET
	weight = EXCLUDED.weight,
	height = EXCLUDED.height,
	age = EXCLUDED.age,
	gender = EXCLUDED.gender,
	daily_calorie_norm = EXCLUDED.daily_calorie_norm,
	profile_complited = EXCLUDED.profile_complited,
	updated_at = now()
	`
	_, err := r.pool.Exec(ctx, query, info.UserId, info.Weight, info.Height, info.Age,
		info.Gender, info.DailyCalorieNorm, info.ProfileComplited)

	return err
}

func (r *postgresUserInformationRepository) GetByUserId(ctx context.Context, userId int) (*model.UserInformation, error) {
	query := `SELECT user_id, weight, height, age, gender, daily_calorie_norm, profile_complited, created_at, updated_at
	FROM user_information WHERE user_id = $1`
	var userInformation model.UserInformation
	err := r.pool.QueryRow(ctx, query, userId).Scan(
		&userInformation.UserId,
		&userInformation.Weight,
		&userInformation.Height,
		&userInformation.Age,
		&userInformation.Gender,
		&userInformation.DailyCalorieNorm,
		&userInformation.ProfileComplited,
		&userInformation.CreateAt,
		&userInformation.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
	}
	return &userInformation, nil
}
