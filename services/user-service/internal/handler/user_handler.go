package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/sergssclaude/health-tracker/user-service/internal/service"
)

type UserHandler struct {
	service service.UserService
}

func NewUserService(service service.UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	user, err := h.service.Register(r.Context(), req.Email, req.Password)
	if err != nil {
		if errors.Is(err, service.ErrEmailAlreadyExists) {
			http.Error(w, "email already exists", http.StatusConflict)
			return
		}
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	resp := toUserRegisterResponse(user)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Bad request body", http.StatusBadRequest)
		return
	}

	token, err := h.service.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(LoginResponse{Token: token})

}

func (h *UserHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	userId, ok := r.Context().Value(UserIDKey).(int)
	if !ok {
		http.Error(w, "user unauthorized", http.StatusUnauthorized)
		return
	}
	userProfile, err := h.service.GetProfile(r.Context(), userId)
	if err != nil {
		http.Error(w, "inernal error", http.StatusInternalServerError)
		return
	}
	resp := toUserGetResponse(userProfile)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func (h *UserHandler) UpdateInformation(w http.ResponseWriter, r *http.Request) {
	userId, ok := r.Context().Value(UserIDKey).(int)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	var req UpdateInformationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Bad request body", http.StatusBadRequest)
		return
	}

	userInfo, err := h.service.UpdateInformation(r.Context(), userId, req.Weight, req.Height, req.Age, req.Gender, req.DailyCalorieNorm)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	resp := toUpdateInformationResponse(userInfo)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)

}

func (h *UserHandler) UpdateGoal(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(UserIDKey).(int)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var req UpdateGoalRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	goal, err := h.service.UpdateGoal(r.Context(), userID, req.TargetWeight, req.CalorieGoal)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	resp := toUserGoalResponse(goal)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}
