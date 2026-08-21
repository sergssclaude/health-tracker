package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"

	"github.com/sergssclaude/health-tracker/nutrition-service/internal/db"
	"github.com/sergssclaude/health-tracker/nutrition-service/internal/handler"
	"github.com/sergssclaude/health-tracker/nutrition-service/internal/repository"
	"github.com/sergssclaude/health-tracker/nutrition-service/internal/service"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found, relying on real env vars")
	}

	dsn := os.Getenv("DATABASE_URL")
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET is not set")
	}

	ctx := context.Background()
	pool, err := db.NewPool(ctx, dsn)
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}
	defer pool.Close()

	foodItemRepo := repository.NewPostgresFoodItemRepository(pool)
	foodLogRepo := repository.NewPostgresFoodLogRepository(pool)

	nutritionService := service.NewNutritionService(foodItemRepo, foodLogRepo)

	nutritionHandler := handler.NewNutritionHandler(nutritionService)

	r := chi.NewRouter()
	r.Group(func(r chi.Router) {
		r.Use(handler.AuthMiddleware(jwtSecret))

		r.Post("/food-logs", nutritionHandler.LogFood)
		r.Get("/food-logs", nutritionHandler.GetDailyLogs)
		r.Delete("/food-logs/{id}", nutritionHandler.DeleteLog)
		r.Get("/food-items/search", nutritionHandler.SearchFoodItems)
	})

	log.Println("starting nutrition-service on :8081")
	if err := http.ListenAndServe(":8081", r); err != nil {
		log.Fatalf("server failed: %v", err)
	}

}
