package grpc

import (
	"context"
	"fmt"
	"log"
	"training-golang/session-11-user-crud-grpc-gateway-cache/entity"
	pb "training-golang/session-11-user-crud-grpc-gateway-cache/proto/user_service/v1"
	"training-golang/session-11-user-crud-grpc-gateway-cache/service"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type UserHandler struct {
	pb.UnimplementedUserServiceServer
	userService service.IUserService
}

func NewUserHandler(userService service.IUserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) GetUsers(ctx context.Context, _ *emptypb.Empty) (*pb.GetUsersResponse, error) {
	users, err := h.userService.GetAllUsers(ctx)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var userProto []*pb.User
	for _, user := range users {
		userProto = append(userProto, &pb.User{
			Id:        int32(user.ID),
			Name:      user.Name,
			Email:     user.Email,
			Password:  user.Password,
			CreatedAt: timestamppb.New(user.CreatedAt),
			UpdatedAt: timestamppb.New(user.UpdatedAt),
		})
	}

	return &pb.GetUsersResponse{
		Users: userProto,
	}, nil
}

func (h *UserHandler) GetUserByID(ctx context.Context, req *pb.GetUserByIDRequest) (*pb.GetUserByIDResponse, error) {
	user, err := h.userService.GetUserByID(ctx, int(req.Id))
	if err != nil {
		log.Println(err)
		return nil, err
	}

	if user.ID == 0 {
		log.Println("User not found for ID:", req.Id)
		return nil, status.Errorf(codes.NotFound, "user with ID %d not found", req.Id)
	}

	log.Println("Fetched user:", user)

	return &pb.GetUserByIDResponse{User: &pb.User{
		Id:        int32(user.ID),
		Name:      user.Name,
		Email:     user.Email,
		Password:  user.Password,
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: timestamppb.New(user.UpdatedAt),
	}}, nil
}

func (h *UserHandler) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.MutationResponse, error) {
	createdUser, err := h.userService.CreateUser(ctx, &entity.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	})

	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &pb.MutationResponse{
		Message: fmt.Sprintf("Success created user with ID %d", createdUser.ID),
	}, nil
}

func (h *UserHandler) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.MutationResponse, error) {
	updatedUser, err := h.userService.UpdateUser(ctx, int(req.Id), entity.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	})

	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &pb.MutationResponse{
		Message: fmt.Sprintf("Success updated user with id %d", updatedUser.ID),
	}, nil
}

func (h *UserHandler) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.MutationResponse, error) {
	if err := h.userService.DeleteUser(ctx, int(req.Id)); err != nil {
		log.Println(err)
		return nil, err
	}

	return &pb.MutationResponse{
		Message: fmt.Sprintf("Success deleted user with id %d", req.Id),
	}, nil
}
