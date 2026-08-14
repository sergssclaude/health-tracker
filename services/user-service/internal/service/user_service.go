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
	GetProfile(ctx context.Context, userID int) (*model.UserProfile, error)
	UpdateInformation(ctx context.Context, userID int, weight *float64, height *int, age *int, gender *string, dailyCalorieNorm *int) (*model.UserInformation, error)
	UpdateGoal(ctx context.Context, userID int, targetWeight float64, calorieGoal int) (*model.UserGoal, error)
}

type userService struct {
	userRepo        repository.UserRepository
	userInfoRepo    repository.UserInformationRepository
	userGoalRepo    repository.UserGoalRepository
	userProfileRepo repository.UserProfileRepository
	jwtSecret       string
}

func NewUserService(userRepo repository.UserRepository, userInfoRepo repository.UserInformationRepository,
	userGoalRepo repository.UserGoalRepository, userProfileRepo repository.UserProfileRepository, jwtString string) UserService {
	return &userService{userRepo: userRepo,
		userInfoRepo:    userInfoRepo,
		userGoalRepo:    userGoalRepo,
		userProfileRepo: userProfileRepo,
		jwtSecret:       jwtString}
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

func (s *userService) GetProfile(ctx context.Context, userID int) (*model.UserProfile, error) {
	userProfile, err := s.userProfileRepo.GetFullProfile(ctx, userID)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return userProfile, nil
}

func (s *userService) UpdateInformation(ctx context.Context, userID int, weight *float64, height *int, age *int, gender *string, dailyCalorieNorm *int) (*model.UserInformation, error) {
	//TODO
	return nil, nil
}

func (s *userService) UpdateGoal(ctx context.Context, userID int, targetWeight float64, calorieGoal int) (*model.UserGoal, error) {
	//TODO
	return nil, nil
}
