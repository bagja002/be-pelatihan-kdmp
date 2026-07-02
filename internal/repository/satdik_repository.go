package repository

import (
	"errors"

	"knmp-backend/internal/entity"

	"gorm.io/gorm"
)

var ErrSatdikNotFound = errors.New("satdik not found")

type SatdikRepository interface {
	Create(e *entity.Satdik) error
	FindAll() ([]entity.Satdik, error)
	FindByID(id uint) (*entity.Satdik, error)
	FindByKode(kode string) (*entity.Satdik, error)
	FindByNama(nama string) (*entity.Satdik, error)
	Update(e *entity.Satdik) error
	Delete(id uint) error
}

type satdikRepository struct{ db *gorm.DB }

func NewSatdikRepository(db *gorm.DB) SatdikRepository { return &satdikRepository{db: db} }

func (r *satdikRepository) Create(e *entity.Satdik) error { return r.db.Create(e).Error }

func (r *satdikRepository) FindAll() ([]entity.Satdik, error) {
	var items []entity.Satdik
	err := r.db.Order("nama asc").Find(&items).Error
	return items, err
}

func (r *satdikRepository) FindByID(id uint) (*entity.Satdik, error) {
	var e entity.Satdik
	if err := r.db.First(&e, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrSatdikNotFound
		}
		return nil, err
	}
	return &e, nil
}

func (r *satdikRepository) FindByKode(kode string) (*entity.Satdik, error) {
	var e entity.Satdik
	if err := r.db.Where("kode = ?", kode).First(&e).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrSatdikNotFound
		}
		return nil, err
	}
	return &e, nil
}

func (r *satdikRepository) FindByNama(nama string) (*entity.Satdik, error) {
	var e entity.Satdik
	if err := r.db.Where("nama = ?", nama).First(&e).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrSatdikNotFound
		}
		return nil, err
	}
	return &e, nil
}

func (r *satdikRepository) Update(e *entity.Satdik) error { return r.db.Save(e).Error }

func (r *satdikRepository) Delete(id uint) error {
	return r.db.Delete(&entity.Satdik{}, id).Error
}
