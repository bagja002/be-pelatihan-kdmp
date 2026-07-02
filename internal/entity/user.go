package entity

import (
	"time"

	"knmp-backend/internal/database"

	"gorm.io/gorm"
)

// User is the authentication principal.
type User struct {
	ID    uint   `gorm:"primaryKey" json:"id"`
	Email string `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`
	// Password stores the Argon2id hash. json:"-" ensures it is never
	// serialized in any API response.
	Password string `gorm:"type:varchar(255);not null" json:"-"`
	Role     string `gorm:"type:varchar(32);not null;default:user" json:"role"`
	// Phone is example PII, encrypted at rest with AES-256-GCM via the
	// "encrypted" serializer. The column stores ciphertext; the Go field
	// holds plaintext after decryption.
	Phone     string         `gorm:"type:varchar(512);serializer:encrypted" json:"phone,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func init() {
	database.RegisterModel(&User{})
}
