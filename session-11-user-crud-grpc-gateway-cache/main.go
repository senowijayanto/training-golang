package main

import (
	"context"
	"log"
	"net"
	grpcHandler "training-golang/session-11-user-crud-grpc-gateway-cache/handler/grpc"
	"training-golang/session-11-user-crud-grpc-gateway-cache/middleware"
	pb "training-golang/session-11-user-crud-grpc-gateway-cache/proto/user_service/v1"
	postgresgormraw "training-golang/session-11-user-crud-grpc-gateway-cache/repository/postgres_gorm_raw"
	"training-golang/session-11-user-crud-grpc-gateway-cache/service"

	"github.com/gin-gonic/gin"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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

	// setup redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	// setup service
	userRepo := postgresgormraw.NewUserRepository(gormDB)
	userService := service.NewUserService(userRepo, rdb)
	userHandler := grpcHandler.NewUserHandler(userService)

	// run the grpc server
	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(middleware.UnaryAuthInterceptor()))
	pb.RegisterUserServiceServer(grpcServer, userHandler)

	lis, err := net.Listen("tcp", "localhost:50051")
	if err != nil {
		log.Fatalln(err)
	}

	go func() {
		log.Println("Running server on port :50051")
		grpcServer.Serve(lis)
	}()

	// run the grpc gateway
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalln("Failed to dial server:", err)
	}

	gwmux := runtime.NewServeMux()
	if err = pb.RegisterUserServiceHandler(context.Background(), gwmux, conn); err != nil {
		log.Fatalln("Failed to register gateway:", err)
	}

	// run gin server
	ginServer := gin.Default()

	ginServer.Group("/v1/*{grpc_gateway}").Any("", gin.WrapH(gwmux))

	log.Println("Running grpc gateway server in port :8080")

	ginServer.Run("localhost:8080")
}
