package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func main() {
	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()

	// setup service
	var mockUserDBInSlice []User
	userRepo := NewUserRepository(mockUserDBInSlice)
	userService := NewUserService(userRepo)
	userHandler := NewUserHandler(userService)

	// routes
	SetupRouter(r, userHandler)

	log.Println("Running server on port 8080")
	r.Run("localhost:8080")
}

// IUserRepository mendefinisikan interface untuk repository pengguna
type IUserRepository interface {
	GetAllUsers() []User
}

// userRepository adalah implementasi dari IUserRepository yang menggunakan slice untuk menyimpan data pengguna
type userRepository struct {
	db     []User // slice untuk menyimpan data pengguna
	nextID int    // ID berikutnya yang akan digunakan untuk pengguna baru
}

// NewUserRepository membuat instance baru dari userRepository
func NewUserRepository(db []User) IUserRepository {
	return &userRepository{
		db:     db,
		nextID: 1,
	}
}

// GetAllUsers mengembalikan semua pengguna
func (r *userRepository) GetAllUsers() []User {
	return r.db
}

type IUserService interface {
	GetAllUsers() []User
}

type userService struct {
	userRepo IUserRepository
}

func NewUserService(userRepo IUserRepository) IUserService {
	return &userService{userRepo: userRepo}
}

func (s *userService) GetAllUsers() []User {
	return s.userRepo.GetAllUsers()
}

type IUserHandler interface {
	GetAllUsers(c *gin.Context)
}

type UserHandler struct {
	userService IUserService
}

func NewUserHandler(userService IUserService) IUserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) GetAllUsers(c *gin.Context) {
	users := h.userService.GetAllUsers()
	c.JSON(http.StatusOK, users)
}

func SetupRouter(r *gin.Engine, userHandler IUserHandler) {
	usersPublicEndpoint := r.Group("/users")

	usersPublicEndpoint.GET("/", userHandler.GetAllUsers)
}
