package model

import "time"

type FoodItem struct {
	ID              int
	Name            string
	CaloriesPer100g int
	ProteinPer100g  *int
	FatsPer100g     *int
	CarbsPer100g    *int
}

type FoodLog struct {
	ID          int
	UserID      int
	FoodItemID  int
	AmountGrams int
	MealType    string
	LoggedAt    time.Time
}
