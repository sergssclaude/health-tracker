package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"github.com/sergssclaude/health-tracker/user-service/internal/db"
	"github.com/sergssclaude/health-tracker/user-service/internal/handler"
	"github.com/sergssclaude/health-tracker/user-service/internal/repository"
	"github.com/sergssclaude/health-tracker/user-service/internal/service"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("no .env file found")
	}
	dsn := os.Getenv("DATABASE_USERS_URL")
	if dsn == "" {
		log.Fatal("Database url is not set")
	}
	jwt := os.Getenv("JWT_SECRET")
	if jwt == "" {
		log.Fatal("Jwt secret is not set")
	}

	ctx := context.Background()
	pool, err := db.NewPool(ctx, dsn)
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}
	defer pool.Close()

	userRepo := repository.NewUserRepository(pool)
	userInfoRepo := repository.NewUserInformationRepository(pool)
	userGoalRepo := repository.NewUserGoalRepository(pool)
	userProfileRepo := repository.NewUserProfileRepository(pool)

	userService := service.NewUserService(userRepo, userInfoRepo, userGoalRepo, userProfileRepo, jwt)
	userHandler := handler.NewUserService(userService)

	r := chi.NewRouter()
	//Проверить 2 ендпоинта register, login (работают)
	r.Post("/register", userHandler.Register)
	r.Post("/login", userHandler.Login)

	r.Group(func(r chi.Router) {
		r.Use(handler.AuthMiddleware(jwt))
		r.Get("/profile", userHandler.GetProfile)
		// r.Put("/profile/information", userHandler.UpdateInformation)
		// r.Get("/profile/goal", userHandler.UpdateGoal)
	})

	log.Println("starting server on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("failed start server :8080, %v", err)
	}
}
