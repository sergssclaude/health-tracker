package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"github.com/pressly/goose/v3"
	"github.com/sergssclaude/health-tracker/nutrition-service/migrations"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("no .env file found")
	}
	dsn := os.Getenv("DATABASE_USERS_URL")
	if dsn == "" {
		log.Fatal("Database url is not set")
	}

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatalf("failed to open db connection: %v", err)
	}

	goose.SetBaseFS(migrations.FS)

	if err := goose.SetDialect("postgres"); err != nil {
		log.Fatalf("failed to set dialect: %v", err)
	}

	if err := goose.Up(db, "."); err != nil {
		log.Fatalf("migration failed: %v", err)
	}

}
