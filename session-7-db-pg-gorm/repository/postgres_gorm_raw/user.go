package postgresgormraw

import (
	"context"
	"errors"
	"log"
	"training-golang/session-7-db-pg-gorm/entity"

	"gorm.io/gorm"
)

type GormDBIface interface {
	WithContext(ctx context.Context) *gorm.DB
	Create(value interface{}) *gorm.DB
	First(dest interface{}, conds ...interface{}) *gorm.DB
	Save(value interface{}) *gorm.DB
	Delete(value interface{}, conds ...interface{}) *gorm.DB
	Find(dest interface{}, conds ...interface{}) *gorm.DB
}

// IUserRepository mendefinisikan interface untuk repository pengguna
type IUserRepository interface {
	CreateUser(ctx context.Context, user *entity.User) (entity.User, error)
	GetUserByID(ctx context.Context, id int) (entity.User, error)
	UpdateUser(ctx context.Context, id int, user entity.User) (entity.User, error)
	DeleteUser(ctx context.Context, id int) error
	GetAllUsers(ctx context.Context) ([]entity.User, error)
}

// userRepository adalah implementasi dari IUserRepository yang menggunakan db untuk menyimpan data pengguna
type userRepository struct {
	db GormDBIface
}

// NewUserRepository membuat instance baru dari userRepository
func NewUserRepository(db GormDBIface) IUserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) CreateUser(ctx context.Context, user *entity.User) (entity.User, error) {
	query := "INSERT INTO users(name, email, password, created_at, updated_at) VALUES ($1, $2, $3, NOW(), NOW()) RETURNING id"
	var createdID int
	if err := r.db.WithContext(ctx).Raw(query, user.Name, user.Email, user.Password).Scan(&createdID).Error; err != nil {
		log.Printf("Error creating user: %v\n", err)
		return entity.User{}, err
	}

	user.ID = createdID
	return *user, nil
}

func (r *userRepository) GetUserByID(ctx context.Context, id int) (entity.User, error) {
	var user entity.User
	query := "SELECT id, name, email, password, created_at, updated_at FROM users WHERE id = $1"
	if err := r.db.WithContext(ctx).Raw(query, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.User{}, nil
		}
		log.Printf("Error getting user by ID: %v\n", err)
		return entity.User{}, err
	}
	return user, nil
}

func (r *userRepository) UpdateUser(ctx context.Context, id int, user entity.User) (entity.User, error) {
	query := "UPDATE users SET name = $1, email = $2, password = $3, updated_at = NOW() WHERE id = $4"
	if err := r.db.WithContext(ctx).Exec(query, user.Name, user.Email, user.Password, id).Error; err != nil {
		log.Printf("Error updating users: %v\n", err)
		return entity.User{}, err
	}

	return user, nil
}

func (r *userRepository) DeleteUser(ctx context.Context, id int) error {
	query := "DELETE FROM users WHERE id = $1"
	if err := r.db.WithContext(ctx).Exec(query, id).Error; err != nil {
		log.Printf("Error deleting user: %v\n", err)
		return err
	}
	return nil
}

// GetAllUsers mengembalikan semua pengguna
func (r *userRepository) GetAllUsers(ctx context.Context) ([]entity.User, error) {
	var users []entity.User
	query := "SELECT id, name, email, password, created_at, updated_at FROM users"
	if err := r.db.WithContext(ctx).Raw(query).Scan(&users).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return users, nil
		}
		log.Printf("Error getting all users: %v\n", err)
		return nil, err
	}
	return users, nil
}
