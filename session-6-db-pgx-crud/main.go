package main

import (
	"context"
	"log"
	"training-golang/session-6-db-pgx-crud/handler"
	postgrespgx "training-golang/session-6-db-pgx-crud/repository/postgres_pgx"
	"training-golang/session-6-db-pgx-crud/router"
	"training-golang/session-6-db-pgx-crud/service"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()

	dsn := "postgresql://postgres:postgres@localhost:5432/go_db"
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		log.Fatalln(err)
	}

	// setup service
	userRepo := postgrespgx.NewUserRepository(pool)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	// routes
	router.SetupRouter(r, userHandler)

	log.Println("Running server on port 8080")
	r.Run("localhost:8080")
}
