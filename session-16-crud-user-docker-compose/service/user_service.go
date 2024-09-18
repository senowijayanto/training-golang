package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"session-16-crud-user-docker-compose/entity"
	postgresgormraw "session-16-crud-user-docker-compose/repository/postgres_gorm_raw"
	"time"

	"github.com/redis/go-redis/v9"
)

type IUserService interface {
	CreateUser(ctx context.Context, user *entity.User) (entity.User, error)
	GetUserByID(ctx context.Context, id int) (entity.User, error)
	UpdateUser(ctx context.Context, id int, user entity.User) (entity.User, error)
	DeleteUser(ctx context.Context, id int) error
	GetAllUsers(ctx context.Context) ([]entity.User, error)
}

type userService struct {
	userRepo postgresgormraw.IUserRepository
	rdb      *redis.Client
}

const redisUserByIDKey = "user:%d"

func NewUserService(userRepo postgresgormraw.IUserRepository, rdb *redis.Client) IUserService {
	return &userService{userRepo: userRepo, rdb: rdb}
}

func (s *userService) CreateUser(ctx context.Context, user *entity.User) (entity.User, error) {
	createdUser, err := s.userRepo.CreateUser(ctx, user)
	if err != nil {
		return entity.User{}, fmt.Errorf("error created user: %v", err)
	}

	// serialize the user to JSON
	createdUserJSON, err := json.Marshal(createdUser)
	if err != nil {
		log.Println("failed to marshal user to JSON:", err)
		return createdUser, err
	}

	// set the redis key with a 1-minute expiration
	redisKey := fmt.Sprintf(redisUserByIDKey, createdUser.ID)
	if err := s.rdb.Set(ctx, redisKey, createdUserJSON, 1*time.Minute).Err(); err != nil {
		log.Println("failed to set redis key with expiration:", err)
		return createdUser, err
	}

	// invalidate the "all_users" cache to ensure it will be refreshed
	if err := s.rdb.Del(ctx, "all_users").Err(); err != nil {
		log.Println("failed to validte 'all_users' cache:", err)
	}

	return createdUser, nil
}

func (s *userService) GetUserByID(ctx context.Context, id int) (entity.User, error) {
	var user entity.User
	redisKey := fmt.Sprintf(redisUserByIDKey, id)

	// try to get teh user from redis cache
	val, err := s.rdb.Get(ctx, redisKey).Result()
	if err == nil {
		// unmarshal the cached data
		if err = json.Unmarshal([]byte(val), &user); err == nil {
			return user, nil
		}
		log.Println("failed to unmarshal user data from redis")
	}

	// if cache miss or unmarshal failed, fetch from database
	user, err = s.userRepo.GetUserByID(ctx, id)
	if err != nil {
		return entity.User{}, fmt.Errorf("user not found: %v", err)
	}

	// cache the user data with an expiration time
	userJSON, _ := json.Marshal(user)
	if err = s.rdb.Set(ctx, redisKey, userJSON, 1*time.Minute).Err(); err != nil {
		log.Println("failed to set user data to redis:", err)
	}
	return user, nil
}

func (s *userService) UpdateUser(ctx context.Context, id int, user entity.User) (entity.User, error) {
	updatedUser, err := s.userRepo.UpdateUser(ctx, id, user)
	if err != nil {
		return entity.User{}, fmt.Errorf("failed updated user: %v", err)
	}

	redisKey := fmt.Sprintf(redisUserByIDKey, updatedUser.ID)
	updatedUserJSON, _ := json.Marshal(updatedUser)
	if err = s.rdb.Set(ctx, redisKey, updatedUserJSON, 1*time.Minute).Err(); err != nil {
		log.Println("failed to set user data to redis:", err)
	}

	// invalidate the "all_users" cache to ensure it will be refreshed
	if err := s.rdb.Del(ctx, "all_users").Err(); err != nil {
		log.Println("failed to validte 'all_users' cache:", err)
	}

	return updatedUser, nil
}

func (s *userService) DeleteUser(ctx context.Context, id int) error {
	redisKey := fmt.Sprintf(redisUserByIDKey, id)

	// delete user cache
	if err := s.rdb.Del(ctx, redisKey).Err(); err != nil {
		log.Println("failed to delete user cache:", err)
	}

	// invalidate the "all_users" cache to ensure it will be refreshed
	if err := s.rdb.Del(ctx, "all_users").Err(); err != nil {
		log.Println("failed to validte 'all_users' cache:", err)
	}

	err := s.userRepo.DeleteUser(ctx, id)
	if err != nil {
		return fmt.Errorf("failed deleted user: %v", err)
	}
	return nil
}

func (s *userService) GetAllUsers(ctx context.Context) ([]entity.User, error) {
	redisKey := "all_users"
	var users []entity.User

	// try to get all users from redis cache
	val, err := s.rdb.Get(ctx, redisKey).Result()
	if err == nil {
		// unmarshal teh cached data
		if err = json.Unmarshal([]byte(val), &users); err == nil {
			return users, nil
		}
		log.Println("failed to unmarshal users data from redis")
	}

	// if cache miss or unmarshal failed, fetch from database
	users, err = s.userRepo.GetAllUsers(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieved users: %v", err)
	}

	// cache the users data with an expiration time
	usersJSON, _ := json.Marshal(users)
	if err = s.rdb.Set(ctx, redisKey, usersJSON, 1*time.Minute).Err(); err != nil {
		log.Println("failed to set users data to redis:", err)
	}
	return users, nil
}
