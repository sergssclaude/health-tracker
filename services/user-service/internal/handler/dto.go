package handler

import (
	"github.com/sergssclaude/health-tracker/user-service/internal/model"
)

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserRegisterResponse struct {
	Id    int    `json:"id"`
	Email string `json:"email"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

func toUserRegisterResponse(u *model.User) UserRegisterResponse {
	return UserRegisterResponse{
		Id:    u.Id,
		Email: u.Email}
}

// type UserGetResponse struct {
// 	Id                 int        `json:"id"`
// 	Email              string     `json:"email"`
// 	Password           string     `json:"password"`
// 	Name               *string    `json:"name"`
// 	Weight             *int       `json:"weight"`
// 	Height             *int       `json:"height"`
// 	Age                *int       `json:"age"`
// 	Gender             *int       `json:"gender"`
// 	Goal               *int       `json:"goal"`
// 	DailyCalorieTarget *int       `json:"daily_calorie_target"`
// 	CreatedAt          time.Time  `json:"created_at"`
// 	UpdatedAt          *time.Time `json:"updated_at"`
// }

// func toUserGetResponse(u *model.User) UserGetResponse {
// 	return UserGetResponse{
// 		Id:                 u.Id,
// 		Email:              u.Email,
// 		Password:           u.PasswordHash,
// 		Name:               u.Name,
// 		Weight:             u.Weight,
// 		Height:             u.Height,
// 		Age:                u.Age,
// 		Gender:             u.Gender,
// 		Goal:               u.Goal,
// 		DailyCalorieTarget: u.DailyCalorieTarget,
// 		CreatedAt:          u.CreatedAt,
// 		UpdatedAt:          u.UpdatedAt}
// }
