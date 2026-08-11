package model

import "time"

type User struct {
	Id           int
	Email        string
	PasswordHash string
	Name         *string
	CreatedAt    time.Time
	UpdatedAt    *time.Time
}

type UserInformation struct {
	UserId           int
	Weight           *float64
	Height           *int
	Age              *int
	Gender           *int
	DailyCalorieNorm *int
	ProfileComplited bool
	CreateAt         time.Time
	UpdatedAt        *time.Time
}

type UserGoal struct {
	UserId       int
	TargetWeight *float64
	CalorieGoal  *int
	CreateAt     time.Time
	UpdatedAt    *time.Time
}
