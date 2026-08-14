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

type UserProfileGetResponse struct {
	User        UserResponse            `json:"user"`
	Information UserInformationResponse `json:"information"`
	Goal        UserGoalResponse        `json:"goal"`
}

type UserResponse struct {
	ID    int     `json:"id"`
	Email string  `json:"email"`
	Name  *string `json:"name"`
}

type UserInformationResponse struct {
	Weight           *float64 `json:"weight"`
	Height           *int     `json:"height"`
	Age              *int     `json:"age"`
	Gender           *string  `json:"gender"`
	DailyCalorieNorm *int     `json:"daily_calorie_target"`
	ProfileCompleted bool     `json:"profile_completed"`
}

type UserGoalResponse struct {
	TargetWeight *float64 `json:"target_weight"`
	CalorieGoal  *int     `json:"calorie_goal"`
}

func toUserGetResponse(u *model.UserProfile) UserProfileGetResponse {
	return UserProfileGetResponse{
		User: UserResponse{
			ID:    u.User.Id,
			Email: u.User.Email,
			Name:  u.User.Name,
		},
		Information: UserInformationResponse{
			Weight:           u.Information.Weight,
			Height:           u.Information.Height,
			Age:              u.Information.Age,
			Gender:           u.Information.Gender,
			DailyCalorieNorm: u.Information.DailyCalorieNorm,
			ProfileCompleted: u.Information.ProfileComplited,
		},
		Goal: UserGoalResponse{
			TargetWeight: u.Goal.TargetWeight,
			CalorieGoal:  u.Goal.CalorieGoal,
		},
	}
}
