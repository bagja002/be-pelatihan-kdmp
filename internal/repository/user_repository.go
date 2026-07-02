package repository

import (
	"errors"

	"knmp-backend/internal/entity"

	"gorm.io/gorm"
)

// ErrUserNotFound is returned when a user lookup finds no row.
var ErrUserNotFound = errors.New("user not found")

// UserRepository defines data-access operations for User.
type UserRepository interface {
	Create(u *entity.User) error
	FindByEmail(email string) (*entity.User, error)
	FindByID(id uint) (*entity.User, error)
}

type userRepository struct {
	db *gorm.DB
}

// NewUserRepository returns a GORM-backed UserRepository.
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(u *entity.User) error {
	return r.db.Create(u).Error
}

func (r *userRepository) FindByEmail(email string) (*entity.User, error) {
	var u entity.User
	if err := r.db.Where("email = ?", email).First(&u).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &u, nil
}

func (r *userRepository) FindByID(id uint) (*entity.User, error) {
	var u entity.User
	if err := r.db.First(&u, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &u, nil
}
