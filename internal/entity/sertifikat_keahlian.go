package entity

import (
	"time"

	"knmp-backend/internal/database"

	"gorm.io/gorm"
)

// SertifikatKeahlian adalah sertifikat keahlian milik seorang Pelatih (relasi 1:N).
type SertifikatKeahlian struct {
	ID             uint   `gorm:"primaryKey" json:"id"`
	IDPelatih      uint   `gorm:"index;not null" json:"idPelatih"`
	NamaSertifikat string `gorm:"type:varchar(255);not null" json:"namaSertifikat"`
	Berkas         string `gorm:"type:varchar(512)" json:"berkas"` // path relatif berkas di disk

	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func init() { database.RegisterModel(&SertifikatKeahlian{}) }
