package handler

import (
	"net/http"
	"strconv"
	"time"
	"training-golang/session-2-latihan-crud-user-gin/entity"

	"github.com/gin-gonic/gin"
)

var (
	users  []entity.User
	nextID int = 1
)

// Create User
func CreateUser(c *gin.Context) {
	var user entity.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user.ID = nextID
	nextID++
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	users = append(users, user)
	c.JSON(http.StatusCreated, user)
}

// Get User by ID
func GetUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	for _, user := range users {
		if user.ID == id {
			c.JSON(http.StatusOK, user)
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "User not found!"})
}

// Update User
func UpdateUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var user entity.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for i, u := range users {
		updateUser := entity.User{
			ID:        id,
			Name:      user.Name,
			Email:     user.Email,
			Password:  u.Password,
			CreatedAt: u.CreatedAt,
			UpdatedAt: time.Now(),
		}

		users[i] = updateUser
		c.JSON(http.StatusOK, updateUser)
		return
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "User not found!"})
}

// Delete User
func DeleteUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	// users := []int{0, 1, 2, 3, 4}
	// i := 2
	// users = append(users[:i], users[i+1]...)
	// users[:i] will be [0, 1]
	// users[i+1] will be [3, 4]
	// users = append([]int{0,1}, []int{3,4}...)
	// users slice will be [0, 1, 3, 4]
	for i, user := range users {
		if user.ID == id {
			users = append(users[:i], users[i+1:]...)
			c.JSON(http.StatusOK, gin.H{"message": "User Deleted"})
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "User not found!"})
}

// Get All Users
func GetAllUsers(c *gin.Context) {
	c.JSON(http.StatusOK, users)
}
