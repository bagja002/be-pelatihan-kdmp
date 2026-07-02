package service

import (
	"knmp-backend/internal/dto"
	"knmp-backend/internal/entity"
	"knmp-backend/internal/repository"
)

// ProductService defines the business logic for Product.
type ProductService interface {
	Create(req *dto.CreateProductRequest) (*entity.Product, error)
	GetAll() ([]entity.Product, error)
	GetByID(id uint) (*entity.Product, error)
	Update(id uint, req *dto.UpdateProductRequest) (*entity.Product, error)
	Delete(id uint) error
}

type productService struct {
	repo repository.ProductRepository
}

// NewProductService wires a ProductService to its repository.
func NewProductService(repo repository.ProductRepository) ProductService {
	return &productService{repo: repo}
}

func (s *productService) Create(req *dto.CreateProductRequest) (*entity.Product, error) {
	e := &entity.Product{
		Name: req.Name,
	}
	if err := s.repo.Create(e); err != nil {
		return nil, err
	}
	return e, nil
}

func (s *productService) GetAll() ([]entity.Product, error) {
	return s.repo.FindAll()
}

func (s *productService) GetByID(id uint) (*entity.Product, error) {
	return s.repo.FindByID(id)
}

func (s *productService) Update(id uint, req *dto.UpdateProductRequest) (*entity.Product, error) {
	e, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	e.Name = req.Name
	if err := s.repo.Update(e); err != nil {
		return nil, err
	}
	return e, nil
}

func (s *productService) Delete(id uint) error {
	return s.repo.Delete(id)
}
