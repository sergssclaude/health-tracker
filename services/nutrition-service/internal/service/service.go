package service

import (
	"context"
	"errors"
	"time"

	"github.com/sergssclaude/health-tracker/nutrition-service/internal/model"
	"github.com/sergssclaude/health-tracker/nutrition-service/internal/repository"
)

var ErrFoodItemNotFound = errors.New("Food item not found")

var ErrFoodLogNotFound = errors.New("food log not found")

type NutritionService interface {
	LogFood(ctx context.Context, UserID, FoodItemID, amountGrams int, mealType string) (*model.FoodLog, error)
	GetDailyLogs(ctx context.Context, UserID int, date time.Time) ([]model.FoodLog, error)
	SearchFoodItems(ctx context.Context, nameItem string) ([]model.FoodItem, error)
	DeleteLog(ctx context.Context, logID, UserID int) error
}

type nutritionService struct {
	FoodItemRepo repository.FoodItemRepository
	FoodLogRepo  repository.FoodLogRepository
}

func NewNutritionService(foodItemRepo repository.FoodItemRepository, foodLogRepo repository.FoodLogRepository) NutritionService {
	return &nutritionService{FoodItemRepo: foodItemRepo, FoodLogRepo: foodLogRepo}
}

func (service *nutritionService) LogFood(ctx context.Context, UserID, FoodItemID, amountGrams int, mealType string) (*model.FoodLog, error) {
	log := model.FoodLog{
		UserID:      UserID,
		FoodItemID:  FoodItemID,
		AmountGrams: amountGrams,
		MealType:    mealType,
	}

	if err := service.FoodLogRepo.Create(ctx, &log); err != nil {
		if err == repository.ErrFoodItemNotFound {
			return nil, ErrFoodItemNotFound
		}
		return nil, err
	}
	return &log, nil
}

func (service *nutritionService) GetDailyLogs(ctx context.Context, UserID int, date time.Time) ([]model.FoodLog, error) {
	logs, err := service.FoodLogRepo.GetByUserAndDate(ctx, UserID, date)
	if err != nil {
		return nil, err
	}
	return logs, nil
}

func (service *nutritionService) SearchFoodItems(ctx context.Context, nameItem string) ([]model.FoodItem, error) {
	foodItems, err := service.FoodItemRepo.Search(ctx, nameItem)
	if err != nil {
		return nil, err
	}
	return foodItems, err
}

func (service *nutritionService) DeleteLog(ctx context.Context, logID, UserID int) error {
	if err := service.FoodLogRepo.Delete(ctx, logID, UserID); err != nil {
		if err == repository.ErrFoodLogNotFound {
			return ErrFoodLogNotFound
		}
		return err
	}
	return nil
}
