package repository

import (
	"errors"
	"strings"

	"knmp-backend/internal/entity"

	"gorm.io/gorm"
)

var (
	ErrPelatihNotFound    = errors.New("pelatih not found")
	ErrSertifikatNotFound = errors.New("sertifikat not found")
)

type PelatihRepository interface {
	ExistsNIP(nip string) (bool, error)
	Create(p *entity.Pelatih) error
	FindAll() ([]entity.Pelatih, error)
	FindByID(id uint) (*entity.Pelatih, error)
	Delete(id uint) error
	FindSertifikat(id uint) (*entity.SertifikatKeahlian, error)
}

type pelatihRepository struct{ db *gorm.DB }

func NewPelatihRepository(db *gorm.DB) PelatihRepository { return &pelatihRepository{db: db} }

func (r *pelatihRepository) ExistsNIP(nip string) (bool, error) {
	var count int64
	if err := r.db.Model(&entity.Pelatih{}).Where("nip = ?", nip).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// Create menyimpan pelatih beserta sertifikatnya (asosiasi otomatis GORM).
func (r *pelatihRepository) Create(p *entity.Pelatih) error { return r.db.Create(p).Error }

func (r *pelatihRepository) FindAll() ([]entity.Pelatih, error) {
	var items []entity.Pelatih
	err := r.db.Preload("Sertifikat").Order("created_at desc").Find(&items).Error
	return items, err
}

func (r *pelatihRepository) FindByID(id uint) (*entity.Pelatih, error) {
	var p entity.Pelatih
	if err := r.db.Preload("Sertifikat").First(&p, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrPelatihNotFound
		}
		return nil, err
	}
	return &p, nil
}

// Delete menghapus (soft) pelatih beserta sertifikatnya.
func (r *pelatihRepository) Delete(id uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id_pelatih = ?", id).Delete(&entity.SertifikatKeahlian{}).Error; err != nil {
			return err
		}
		return tx.Delete(&entity.Pelatih{}, id).Error
	})
}

func (r *pelatihRepository) FindSertifikat(id uint) (*entity.SertifikatKeahlian, error) {
	var s entity.SertifikatKeahlian
	if err := r.db.First(&s, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrSertifikatNotFound
		}
		return nil, err
	}
	return &s, nil
}

// IsDuplicateNIP mengenali error unique-constraint MySQL (kode 1062).
func IsDuplicateNIP(err error) bool {
	return err != nil && strings.Contains(err.Error(), "Duplicate entry")
}
