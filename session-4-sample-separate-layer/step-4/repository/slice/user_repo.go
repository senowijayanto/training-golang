package slice

import "training-golang/session-4-sample-separate-layer/step-4/entity"

// IUserRepository mendefinisikan interface untuk repository pengguna
type IUserRepository interface {
	GetAllUsers() []entity.User
}

// userRepository adalah implementasi dari IUserRepository yang menggunakan slice untuk menyimpan data pengguna
type userRepository struct {
	db     []entity.User // slice untuk menyimpan data pengguna
	nextID int           // ID berikutnya yang akan digunakan untuk pengguna baru
}

// NewUserRepository membuat instance baru dari userRepository
func NewUserRepository(db []entity.User) IUserRepository {
	return &userRepository{
		db:     db,
		nextID: 1,
	}
}

// GetAllUsers mengembalikan semua pengguna
func (r *userRepository) GetAllUsers() []entity.User {
	return r.db
}
