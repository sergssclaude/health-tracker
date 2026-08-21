package handler

import (
	"time"

	"github.com/sergssclaude/health-tracker/nutrition-service/internal/model"
)

type LogFoodRequest struct {
	FoodItemID  int    `json:"food_item_id"`
	AmountGrams int    `json:"amount_grams"`
	MealType    string `json:"meal_type"`
}

type LogFoodResponse struct {
	ID          int       `json:"id"`
	FoodItemID  int       `json:"food_item_id"`
	AmountGrams int       `json:"amount_grams"`
	MealType    string    `json:"meal_type"`
	LoggedAt    time.Time `json:"logged_at"`
}

func toLogFoodResponse(log *model.FoodLog) LogFoodResponse {
	return LogFoodResponse{
		ID:          log.ID,
		FoodItemID:  log.FoodItemID,
		AmountGrams: log.AmountGrams,
		MealType:    log.MealType,
		LoggedAt:    log.LoggedAt,
	}
}

func toFoodLogListResponse(logs []model.FoodLog) []LogFoodResponse {
	resp := make([]LogFoodResponse, 0, len(logs))
	for _, log := range logs {
		resp = append(resp, toLogFoodResponse(&log))
	}
	return resp
}
