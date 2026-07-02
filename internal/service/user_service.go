package service

import (
	"errors"

	"knmp-backend/internal/dto"
	"knmp-backend/internal/entity"
	"knmp-backend/internal/repository"
	"knmp-backend/pkg/hash"
)

var ErrUsernameTaken = errors.New("username already registered")

type UserService interface {
	Create(req *dto.CreateUserRequest) (*entity.User, error)
	GetAll() ([]entity.User, error)
	GetByID(id uint) (*entity.User, error)
	Update(id uint, req *dto.UpdateUserRequest) (*entity.User, error)
	Delete(id uint) error
}

type userService struct{ repo repository.UserRepository }

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) Create(req *dto.CreateUserRequest) (*entity.User, error) {
	if existing, _ := s.repo.FindByUsername(req.Username); existing != nil {
		return nil, ErrUsernameTaken
	}
	hashed, err := hash.Password(req.Password)
	if err != nil {
		return nil, err
	}
	u := &entity.User{
		Nama: req.Nama, Username: req.Username, Password: hashed,
		Type: req.Type, IDSatdik: req.IDSatdik,
	}
	if err := s.repo.Create(u); err != nil {
		return nil, err
	}
	return u, nil
}

func (s *userService) GetAll() ([]entity.User, error)        { return s.repo.FindAll() }
func (s *userService) GetByID(id uint) (*entity.User, error) { return s.repo.FindByID(id) }

func (s *userService) Update(id uint, req *dto.UpdateUserRequest) (*entity.User, error) {
	u, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if req.Nama != "" {
		u.Nama = req.Nama
	}
	if req.Type != "" {
		u.Type = req.Type
	}
	if req.IDSatdik != nil {
		u.IDSatdik = req.IDSatdik
	}
	if req.Password != "" {
		hashed, err := hash.Password(req.Password)
		if err != nil {
			return nil, err
		}
		u.Password = hashed
	}
	if err := s.repo.Update(u); err != nil {
		return nil, err
	}
	return u, nil
}

func (s *userService) Delete(id uint) error { return s.repo.Delete(id) }
