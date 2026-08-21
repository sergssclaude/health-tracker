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

type FoodItemResponse struct {
	ID              int    `json:"id"`
	Name            string `json:"name"`
	CaloriesPer100g int    `json:"calories_per_100g"`
	ProteinPer100g  *int   `json:"protein_per_100g"`
	FatPer100g      *int   `json:"fat_per_100g"`
	CarbsPer100g    *int   `json:"carbs_per_100g"`
}

func toFoodItemListResponse(items []model.FoodItem) []FoodItemResponse {
	resp := make([]FoodItemResponse, 0, len(items))
	for _, item := range items {
		resp = append(resp, FoodItemResponse{
			ID:              item.ID,
			Name:            item.Name,
			CaloriesPer100g: item.CaloriesPer100g,
			ProteinPer100g:  item.ProteinPer100g,
			FatPer100g:      item.FatsPer100g,
			CarbsPer100g:    item.CarbsPer100g,
		})
	}
	return resp
}
