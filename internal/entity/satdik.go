package entity

import (
	"time"

	"knmp-backend/internal/database"

	"gorm.io/gorm"
)

// Satdik adalah satuan pendidikan yang menaungi data peserta.
type Satdik struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Kode      string         `gorm:"type:varchar(64);uniqueIndex;not null" json:"kode"`
	Nama      string         `gorm:"type:varchar(255);not null" json:"nama"`
	Lokasi    string         `gorm:"type:varchar(255)" json:"lokasi"`
	Provinsi  string         `gorm:"type:varchar(128)" json:"provinsi"`
	PicSatdik string         `gorm:"type:varchar(255)" json:"picSatdik"`
	NoPic     string         `gorm:"type:varchar(64)" json:"noPic"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func init() { database.RegisterModel(&Satdik{}) }
