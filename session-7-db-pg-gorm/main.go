package main

import (
	"log"
	"training-golang/session-7-db-pg-gorm/handler"
	postgresgormraw "training-golang/session-7-db-pg-gorm/repository/postgres_gorm_raw"
	"training-golang/session-7-db-pg-gorm/router"
	"training-golang/session-7-db-pg-gorm/service"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()

	dsn := "postgresql://postgres:postgres@localhost:5432/go_db"
	gormDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{SkipDefaultTransaction: true})
	if err != nil {
		log.Fatalln(err)
	}

	// setup service
	userRepo := postgresgormraw.NewUserRepository(gormDB)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	// routes
	router.SetupRouter(r, userHandler)

	log.Println("Running server on port 8080")
	r.Run("localhost:8080")
}
