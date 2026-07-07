package repository

import (
	"errors"
	"strings"

	"knmp-backend/internal/entity"

	"github.com/go-sql-driver/mysql"
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
	FindByNIP(nip string) (*entity.Pelatih, error)
	UpdateSelf(p *entity.Pelatih, deleteSertifikatIDs []uint, newSertifikat []entity.SertifikatKeahlian) error
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

// FindByNIP mencari pelatih berdasarkan NIP (kunci swakelola edit mandiri).
func (r *pelatihRepository) FindByNIP(nip string) (*entity.Pelatih, error) {
	var p entity.Pelatih
	if err := r.db.Preload("Sertifikat").Where("nip = ?", nip).First(&p).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrPelatihNotFound
		}
		return nil, err
	}
	return &p, nil
}

// UpdateSelf memperbarui field skalar pelatih, menghapus sertifikat tertentu,
// dan menambah sertifikat baru — semuanya dalam satu transaksi. NIP tidak diubah.
func (r *pelatihRepository) UpdateSelf(p *entity.Pelatih, deleteSertifikatIDs []uint, newSertifikat []entity.SertifikatKeahlian) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Perbarui hanya kolom yang diizinkan (NIP sengaja tidak termasuk).
		if err := tx.Model(&entity.Pelatih{}).Where("id = ?", p.ID).Select(
			"NamaLengkap", "Pendidikan", "Jurusan", "Universitas", "UnitKerja",
			"Jabatan", "Golongan", "Kriteria", "LokasiTOT", "KelasJabatan", "CV", "Status",
		).Updates(p).Error; err != nil {
			return err
		}
		if len(deleteSertifikatIDs) > 0 {
			if err := tx.Where("id_pelatih = ? AND id IN ?", p.ID, deleteSertifikatIDs).
				Delete(&entity.SertifikatKeahlian{}).Error; err != nil {
				return err
			}
		}
		if len(newSertifikat) > 0 {
			for i := range newSertifikat {
				newSertifikat[i].IDPelatih = p.ID
			}
			if err := tx.Create(&newSertifikat).Error; err != nil {
				return err
			}
		}
		return nil
	})
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
	if err == nil {
		return false
	}
	var myErr *mysql.MySQLError
	if errors.As(err, &myErr) {
		return myErr.Number == 1062
	}
	return strings.Contains(err.Error(), "Duplicate entry")
}
