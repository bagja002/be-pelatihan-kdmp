package repository

import (
	"errors"

	"knmp-backend/internal/entity"

	"gorm.io/gorm"
)

var (
	ErrKategoriNotFound  = errors.New("kategori bahan ajar not found")
	ErrBahanAjarNotFound = errors.New("bahan ajar not found")
)

type BahanAjarRepository interface {
	ListKategori() ([]entity.BahanAjarKategori, error)
	FindKategoriByID(id uint) (*entity.BahanAjarKategori, error)
	ExistsKategoriNama(nama string, excludeID uint) (bool, error)
	MaxKategoriUrutan() (int, error)
	CreateKategori(k *entity.BahanAjarKategori) error
	UpdateKategori(k *entity.BahanAjarKategori) error
	CountItems(kategoriID uint) (int64, error)
	DeleteKategori(id uint) error

	FindItemByID(id uint) (*entity.BahanAjar, error)
	MaxItemUrutan(kategoriID uint) (int, error)
	CreateItem(b *entity.BahanAjar) error
	UpdateItem(b *entity.BahanAjar) error
	DeleteItem(id uint) error
}

type bahanAjarRepository struct{ db *gorm.DB }

func NewBahanAjarRepository(db *gorm.DB) BahanAjarRepository {
	return &bahanAjarRepository{db: db}
}

// ListKategori mengembalikan semua kategori + itemnya, terurut urutan lalu id.
func (r *bahanAjarRepository) ListKategori() ([]entity.BahanAjarKategori, error) {
	var items []entity.BahanAjarKategori
	err := r.db.
		Preload("Items", func(tx *gorm.DB) *gorm.DB { return tx.Order("urutan asc, id asc") }).
		Order("urutan asc, id asc").
		Find(&items).Error
	return items, err
}

func (r *bahanAjarRepository) FindKategoriByID(id uint) (*entity.BahanAjarKategori, error) {
	var k entity.BahanAjarKategori
	if err := r.db.First(&k, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrKategoriNotFound
		}
		return nil, err
	}
	return &k, nil
}

// ExistsKategoriNama memeriksa nama kategori kembar (kecuali id sendiri saat edit).
func (r *bahanAjarRepository) ExistsKategoriNama(nama string, excludeID uint) (bool, error) {
	var count int64
	q := r.db.Model(&entity.BahanAjarKategori{}).Where("nama = ?", nama)
	if excludeID != 0 {
		q = q.Where("id <> ?", excludeID)
	}
	if err := q.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *bahanAjarRepository) MaxKategoriUrutan() (int, error) {
	var n int
	err := r.db.Model(&entity.BahanAjarKategori{}).
		Select("COALESCE(MAX(urutan), 0)").Scan(&n).Error
	return n, err
}

func (r *bahanAjarRepository) CreateKategori(k *entity.BahanAjarKategori) error {
	return r.db.Create(k).Error
}

func (r *bahanAjarRepository) UpdateKategori(k *entity.BahanAjarKategori) error {
	return r.db.Model(&entity.BahanAjarKategori{}).Where("id = ?", k.ID).
		Select("Nama", "Urutan").Updates(k).Error
}

func (r *bahanAjarRepository) CountItems(kategoriID uint) (int64, error) {
	var n int64
	err := r.db.Model(&entity.BahanAjar{}).Where("kategori_id = ?", kategoriID).Count(&n).Error
	return n, err
}

func (r *bahanAjarRepository) DeleteKategori(id uint) error {
	return r.db.Delete(&entity.BahanAjarKategori{}, id).Error
}

func (r *bahanAjarRepository) FindItemByID(id uint) (*entity.BahanAjar, error) {
	var b entity.BahanAjar
	if err := r.db.First(&b, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrBahanAjarNotFound
		}
		return nil, err
	}
	return &b, nil
}

func (r *bahanAjarRepository) MaxItemUrutan(kategoriID uint) (int, error) {
	var n int
	err := r.db.Model(&entity.BahanAjar{}).Where("kategori_id = ?", kategoriID).
		Select("COALESCE(MAX(urutan), 0)").Scan(&n).Error
	return n, err
}

func (r *bahanAjarRepository) CreateItem(b *entity.BahanAjar) error {
	return r.db.Create(b).Error
}

func (r *bahanAjarRepository) UpdateItem(b *entity.BahanAjar) error {
	return r.db.Model(&entity.BahanAjar{}).Where("id = ?", b.ID).
		Select("KategoriID", "Judul", "Urutan", "FilePdf", "FilePpt").Updates(b).Error
}

func (r *bahanAjarRepository) DeleteItem(id uint) error {
	return r.db.Delete(&entity.BahanAjar{}, id).Error
}
