package repository

import (
	"knmp-backend/internal/entity"

	"gorm.io/gorm"
)

// ProductRepository defines data-access operations for Product.
type ProductRepository interface {
	Create(e *entity.Product) error
	FindAll() ([]entity.Product, error)
	FindByID(id uint) (*entity.Product, error)
	Update(e *entity.Product) error
	Delete(id uint) error
}

type productRepository struct {
	db *gorm.DB
}

// NewProductRepository returns a GORM-backed ProductRepository.
func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) Create(e *entity.Product) error {
	return r.db.Create(e).Error
}

func (r *productRepository) FindAll() ([]entity.Product, error) {
	var items []entity.Product
	err := r.db.Find(&items).Error
	return items, err
}

func (r *productRepository) FindByID(id uint) (*entity.Product, error) {
	var item entity.Product
	if err := r.db.First(&item, id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *productRepository) Update(e *entity.Product) error {
	return r.db.Save(e).Error
}

func (r *productRepository) Delete(id uint) error {
	return r.db.Delete(&entity.Product{}, id).Error
}
