package repository

import (
	"context"

	"github.com/gayathriad/go_api/models"

	"gorm.io/gorm"
)

// interface — contract for DB operations
type UserRepository interface {
	GetAll(ctx context.Context) ([]models.User, error)
	Create(ctx context.Context, user models.User) (models.User, error)
	Delete(ctx context.Context, id string) error
	GetByID(ctx context.Context, id string) (models.User, error)
	Update(ctx context.Context, user models.User) (models.User, error)
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
func (r *userRepo) GetAll(ctx context.Context) ([]models.User, error) {
	var users []models.User
	result := r.db.WithContext(ctx).Find(&users)
	return users, result.Error
}

// create user
func (r *userRepo) Create(ctx context.Context, user models.User) (models.User, error) {
	result := r.db.WithContext(ctx).Create(&user)
	return user, result.Error
}

// delete user
func (r *userRepo) Delete(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Unscoped().Delete(&models.User{}, "id = ?", id)
	return result.Error
}

func (r *userRepo) GetByID(ctx context.Context, id string) (models.User, error) {
	var user models.User
	result := r.db.WithContext(ctx).First(&user, "id = ?", id)
	return user, result.Error
}

// update user
func (r *userRepo) Update(ctx context.Context, user models.User) (models.User, error) {
	result := r.db.WithContext(ctx).Save(&user)
	return user, result.Error
}
