package repository

import (
	"github.com/gayathriad/go_api/models"

	"gorm.io/gorm"
)

// interface — contract for DB operations
type UserRepository interface {
	GetAll() ([]models.User, error)
	Create(user models.User) (models.User, error)
	Delete(id string) error
}

// actual implementation
type userRepo struct {
	db *gorm.DB
}

// constructor
func NewUserRepo(db *gorm.DB) UserRepository {
	return &userRepo{db: db}
}

// get all users
func (r *userRepo) GetAll() ([]models.User, error) {
	var users []models.User
	result := r.db.Find(&users)
	return users, result.Error
}

// create user
func (r *userRepo) Create(user models.User) (models.User, error) {
	result := r.db.Create(&user)
	return user, result.Error
}

// delete user
func (r *userRepo) Delete(id string) error {
	result := r.db.Unscoped().Delete(&models.User{}, "id = ?", id)
	return result.Error
}
