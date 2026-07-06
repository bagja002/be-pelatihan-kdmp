package entity

import (
	"time"

	"knmp-backend/internal/database"

	"gorm.io/gorm"
)

// Pelatih adalah data pelatih SDM KNMP yang diisi mandiri lewat link publik.
type Pelatih struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	NamaLengkap string `gorm:"type:varchar(255);not null" json:"namaLengkap"`
	NIP         string `gorm:"column:nip;type:varchar(64);uniqueIndex;not null" json:"nip"` // kunci unik: tanpa data ganda
	Pendidikan  string `gorm:"type:varchar(128)" json:"pendidikan"` // pendidikan terakhir
	Jurusan     string `gorm:"type:varchar(255)" json:"jurusan"`
	Universitas string `gorm:"type:varchar(255)" json:"universitas"`
	UnitKerja   string `gorm:"type:varchar(255)" json:"unitKerja"`
	Jabatan     string `gorm:"type:varchar(128)" json:"jabatan"`
	Golongan    string `gorm:"type:varchar(16)" json:"golongan"`      // golongan PNS, mis. III/a
	Kriteria    string `gorm:"type:varchar(16)" json:"kriteria"`       // Expert | KKP | Non KKP
	LokasiTOT   string `gorm:"column:lokasi_tot;type:varchar(255)" json:"lokasiTot"`
	CV          string `gorm:"column:cv;type:varchar(512)" json:"cv"` // path relatif berkas CV
	Status      string `gorm:"type:varchar(32)" json:"status,omitempty"`

	Sertifikat []SertifikatKeahlian `gorm:"foreignKey:IDPelatih" json:"sertifikat"`

	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func init() { database.RegisterModel(&Pelatih{}) }
