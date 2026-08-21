package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/sergssclaude/health-tracker/nutrition-service/internal/service"
)

type NutritionHandler struct {
	service service.NutritionService
}

func NewNutritionService(s service.NutritionService) *NutritionHandler {
	return &NutritionHandler{service: s}
}

func (h *NutritionHandler) LogFood(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(UserIDKey).(int)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var req LogFoodRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	log, err := h.service.LogFood(r.Context(), userID, req.FoodItemID, req.AmountGrams, req.MealType)
	if err != nil {
		if errors.Is(err, service.ErrFoodItemNotFound) {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	resp := toLogFoodResponse(log)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

func (h *NutritionHandler) GetDailyLogs(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(UserIDKey).(int)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	dateStr := r.URL.Query().Get("date")
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		http.Error(w, "Invalid date format, expect YYYY-MM-DD", http.StatusBadRequest)
		return
	}

	logs, err := h.service.GetDailyLogs(r.Context(), userID, date)
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
	}

	resp := toFoodLogListResponse(logs)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(&resp)
}
