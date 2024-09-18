package main

import (
	"log"
	"net"
	grpcHandler "training-golang/session-9-crud-user-grpc/handler/grpc"
	pb "training-golang/session-9-crud-user-grpc/proto/user_service/v1"
	postgresgormraw "training-golang/session-9-crud-user-grpc/repository/postgres_gorm_raw"
	"training-golang/session-9-crud-user-grpc/service"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	gin.SetMode(gin.ReleaseMode)

	dsn := "postgresql://postgres:postgres@localhost:5432/go_db"
	gormDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{SkipDefaultTransaction: true})
	if err != nil {
		log.Fatalln(err)
	}

	// setup service
	userRepo := postgresgormraw.NewUserRepository(gormDB)
	userService := service.NewUserService(userRepo)
	userHandler := grpcHandler.NewUserHandler(userService)

	// run the grpc server
	grpcServer := grpc.NewServer()
	pb.RegisterUserServiceServer(grpcServer, userHandler)

	lis, err := net.Listen("tcp", "localhost:50051")
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("Running server on port :50051")
	grpcServer.Serve(lis)
}
