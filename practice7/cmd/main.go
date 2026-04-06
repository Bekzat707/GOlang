package main

import (
	"log"
	"os"

	"practice-7/internal/entity"
	v1 "practice-7/internal/controller/http/v1"
	"practice-7/internal/usecase"
	"practice-7/internal/usecase/repo"
	"practice-7/pkg/logger"
	"practice-7/pkg/postgres"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	dsn := os.Getenv("PG_DSN")
	if dsn == "" {
		dsn = "host=localhost user=postgres password=postgres dbname=practice7 port=5432 sslmode=disable"
	}

	pg, err := postgres.New(dsn)
	if err != nil {
		log.Fatalf("Failed to connect to postgres: %v", err)
	}

	pg.Conn.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";")
	pg.Conn.AutoMigrate(&entity.User{})

	userRepo := repo.NewUserRepo(pg)
	userUseCase := usecase.NewUserUseCase(userRepo)
	l := logger.New()

	r := gin.Default()
	v1.NewUserRoutes(r.Group("/api/v1"), userUseCase, l)

	r.Run(":8080")
}
