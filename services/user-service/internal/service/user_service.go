package service

import (
	"context"
	"errors"
	"log"

	"github.com/sergssclaude/health-tracker/user-service/internal/model"
	"github.com/sergssclaude/health-tracker/user-service/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("Invalid email or password")
	ErrEmailAlreadyExists = errors.New("Email already exists")
	ErrUserNotFound       = errors.New("User not found")
)

type UserService interface {
	Register(ctx context.Context, email, password string) (*model.User, error)
	Login(ctx context.Context, email, password string) (string, error)
	GetProfile(ctx context.Context, userID int) (*model.User, error)
	UpdateInformation(ctx context.Context, userID int, weight float64, height, age int, gender, dailyCalorieNorm string) (*model.UserInformation, error)
	UpdateGoal(ctx context.Context, userID int, targetWeight float64, calorieGoal int) (*model.UserGoal, error)
}

type userService struct {
	userRepo     repository.UserRepository
	userInfoRepo repository.UserInformationRepository
	userGoalRepo repository.UserGoalRepository
	jwtSecret    string
}

func NewUserService(userRepo repository.UserRepository, userInfoRepo repository.UserInformationRepository,
	userGoalRepo repository.UserGoalRepository, jwtString string) UserService {
	return &userService{userRepo: userRepo, userInfoRepo: userInfoRepo, userGoalRepo: userGoalRepo, jwtSecret: jwtString}
}

func (s *userService) Register(ctx context.Context, email, password string) (*model.User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user := &model.User{
		Email:        email,
		PasswordHash: string(hash),
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		if errors.Is(err, repository.ErrEmailAlreadyExists) {
			return nil, ErrEmailAlreadyExists
		}
		return nil, err
	}
	return user, nil
}

func (s *userService) Login(ctx context.Context, email, password string) (string, error) {
	user, err := s.userRepo.GetByMail(ctx, email)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return "", ErrInvalidCredentials
		}
		return "", nil
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", ErrInvalidCredentials
	}

	token, err := generateJWT(user.Id, s.jwtSecret)
	if err != nil {
		log.Printf("failed generater JWT token for user %v: %v", user.Id, err)
		return "", err
	}
	return token, nil
}

func (s *userService) GetProfile(ctx context.Context, userID int) (*model.User, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return user, nil
}

func (s *userService) UpdateProfile(ctx context.Context, userID, weight, height, age int, gender, goal string) (*model.User, error) {
	//TODO
	return nil, nil
}
