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
	// userId, ok := r.Context().Value(UserIDKey).(int)
	// if !ok {
	// 	http.Error(w, "user unauthorized", http.StatusUnauthorized)
	// 	return
	// }
	// user, err := h.service.GetProfile(r.Context(), userId)
	// if err != nil {
	// 	http.Error(w, "inernal error", http.StatusInternalServerError)
	// 	return
	// }
	// resp := toUserGetResponse(user)
	// w.Header().Set("Content-Type", "application/json")
	// w.WriteHeader(http.StatusOK)
	// json.NewEncoder(w).Encode(resp)

}

func (h *UserHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {

}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	// idStr := chi.URLParam(r, "id")
	// id, err := strconv.Atoi(idStr)
	// if err != nil {
	// 	http.Error(w, "bad path parametr", http.StatusBadRequest)
	// 	return
	// }

	// user, err := h.service.GetProfile(r.Context(), id)
	// if err != nil {
	// 	http.Error(w, "inernal error", http.StatusInternalServerError)
	// 	return
	// }
	// resp := toUserGetResponse(user)
	// w.Header().Set("Content-Type", "application/json")
	// w.WriteHeader(http.StatusOK)
	// json.NewEncoder(w).Encode(resp)
}
