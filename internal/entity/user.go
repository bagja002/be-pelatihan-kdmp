package entity

import (
	"time"

	"knmp-backend/internal/database"

	"gorm.io/gorm"
)

// User adalah principal internal (super_admin / admin operator satdik).
type User struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	Nama     string `gorm:"type:varchar(255);not null" json:"nama"`
	Username string `gorm:"type:varchar(128);uniqueIndex;not null" json:"username"`
	Password string `gorm:"type:varchar(255);not null" json:"-"`
	Type     string `gorm:"type:varchar(32);not null" json:"type"` // super_admin | admin
	IDSatdik *uint  `gorm:"index" json:"idSatdik"`                 // null untuk super_admin

	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func init() { database.RegisterModel(&User{}) }
