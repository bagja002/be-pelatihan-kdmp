package entity

import (
	"time"

	"knmp-backend/internal/database"

	"gorm.io/gorm"
)

// Product represents the product domain model.
// TODO: adjust the fields below to match your needs.
type Product struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"type:varchar(255);not null" json:"name"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// init registers this model for auto-migration, so you only ever need to
// register the route — the database migration happens automatically.
func init() {
	database.RegisterModel(&Product{})
}
