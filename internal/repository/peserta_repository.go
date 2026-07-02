package repository

import (
	"errors"

	"knmp-backend/internal/entity"

	"gorm.io/gorm"
)

var ErrPesertaNotFound = errors.New("peserta not found")

type PesertaRepository interface {
	Create(e *entity.DataPeserta) error
	FindAllScoped(satdikID *uint) ([]entity.DataPeserta, error)
	FindByID(id uint) (*entity.DataPeserta, error)
	FindByNIKAndSatdik(nik string, satdikID uint) (*entity.DataPeserta, error)
	Update(e *entity.DataPeserta) error
	Delete(id uint) error
	AllNIK() (map[string]bool, error)
	CreateBatch(items []entity.DataPeserta) error
}

type pesertaRepository struct{ db *gorm.DB }

func NewPesertaRepository(db *gorm.DB) PesertaRepository { return &pesertaRepository{db: db} }

func (r *pesertaRepository) Create(e *entity.DataPeserta) error { return r.db.Create(e).Error }

// FindAllScoped: nil satdikID → semua (super_admin); non-nil → filter satdik.
func (r *pesertaRepository) FindAllScoped(satdikID *uint) ([]entity.DataPeserta, error) {
	var items []entity.DataPeserta
	q := r.db.Order("nama asc")
	if satdikID != nil {
		q = q.Where("id_satdik = ?", *satdikID)
	}
	err := q.Find(&items).Error
	return items, err
}

func (r *pesertaRepository) FindByID(id uint) (*entity.DataPeserta, error) {
	var e entity.DataPeserta
	if err := r.db.First(&e, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrPesertaNotFound
		}
		return nil, err
	}
	return &e, nil
}

func (r *pesertaRepository) FindByNIKAndSatdik(nik string, satdikID uint) (*entity.DataPeserta, error) {
	var e entity.DataPeserta
	if err := r.db.Where("nik = ? AND id_satdik = ?", nik, satdikID).First(&e).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrPesertaNotFound
		}
		return nil, err
	}
	return &e, nil
}

func (r *pesertaRepository) Update(e *entity.DataPeserta) error { return r.db.Save(e).Error }

func (r *pesertaRepository) Delete(id uint) error {
	return r.db.Delete(&entity.DataPeserta{}, id).Error
}

func (r *pesertaRepository) AllNIK() (map[string]bool, error) {
	var niks []string
	if err := r.db.Model(&entity.DataPeserta{}).Where("nik <> ''").Pluck("nik", &niks).Error; err != nil {
		return nil, err
	}
	set := make(map[string]bool, len(niks))
	for _, n := range niks {
		set[n] = true
	}
	return set, nil
}

func (r *pesertaRepository) CreateBatch(items []entity.DataPeserta) error {
	if len(items) == 0 {
		return nil
	}
	return r.db.CreateInBatches(items, 500).Error
}
