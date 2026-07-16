package entity

import (
	"time"

	"knmp-backend/internal/database"

	"gorm.io/gorm"
)

// BahanAjarKategori mengelompokkan bahan ajar (mis. "Kompetensi Umum").
// Nama unik ditegakkan di service (bukan unique index) agar nama bekas
// kategori yang di-soft-delete bisa dipakai ulang.
type BahanAjarKategori struct {
	ID     uint   `gorm:"primaryKey" json:"id"`
	Nama   string `gorm:"type:varchar(255);not null;index" json:"nama"`
	Urutan int    `json:"urutan"`

	Items []BahanAjar `gorm:"foreignKey:KategoriID" json:"items"`

	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// BahanAjar adalah satu unit kompetensi/modul dengan slot berkas PDF & PPT.
// Path berkas relatif terhadap upload root dan tidak diekspos ke JSON.
type BahanAjar struct {
	ID         uint   `gorm:"primaryKey" json:"id"`
	KategoriID uint   `gorm:"index;not null" json:"kategoriId"`
	Judul      string `gorm:"type:varchar(255);not null" json:"judul"`
	Urutan     int    `json:"urutan"`
	FilePdf    string `gorm:"type:varchar(512)" json:"-"`
	FilePpt    string `gorm:"type:varchar(512)" json:"-"`

	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func init() {
	database.RegisterModel(&BahanAjarKategori{})
	database.RegisterModel(&BahanAjar{})
}
