package repository

import (
	"errors"

	"knmp-backend/internal/entity"

	"gorm.io/gorm"
)

var ErrUserNotFound = errors.New("user not found")

type UserRepository interface {
	Create(u *entity.User) error
	FindAll() ([]entity.User, error)
	FindByUsername(username string) (*entity.User, error)
	FindByID(id uint) (*entity.User, error)
	Update(u *entity.User) error
	Delete(id uint) error
}

type userRepository struct{ db *gorm.DB }

func NewUserRepository(db *gorm.DB) UserRepository { return &userRepository{db: db} }

func (r *userRepository) Create(u *entity.User) error { return r.db.Create(u).Error }

func (r *userRepository) FindAll() ([]entity.User, error) {
	var items []entity.User
	err := r.db.Order("nama asc").Find(&items).Error
	return items, err
}

func (r *userRepository) FindByUsername(username string) (*entity.User, error) {
	var u entity.User
	if err := r.db.Where("username = ?", username).First(&u).Error; err != nil {
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

func (r *userRepository) Update(u *entity.User) error { return r.db.Save(u).Error }

func (r *userRepository) Delete(id uint) error {
	return r.db.Delete(&entity.User{}, id).Error
}
