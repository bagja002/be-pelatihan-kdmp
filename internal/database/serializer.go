package database

import (
	"context"
	"fmt"
	"reflect"

	"knmp-backend/pkg/crypto"

	"gorm.io/gorm/schema"
)

// EncryptedSerializer transparently encrypts a string field with AES-256-GCM
// before it is written to the database and decrypts it on read. Apply it with:
//
//	Field string `gorm:"type:varchar(512);serializer:encrypted"`
type EncryptedSerializer struct{}

func init() {
	schema.RegisterSerializer("encrypted", EncryptedSerializer{})
}

// Scan decrypts the stored value into the Go field.
func (EncryptedSerializer) Scan(ctx context.Context, field *schema.Field, dst reflect.Value, dbValue any) error {
	target := field.ReflectValueOf(ctx, dst)

	if dbValue == nil {
		target.SetString("")
		return nil
	}

	var enc string
	switch v := dbValue.(type) {
	case string:
		enc = v
	case []byte:
		enc = string(v)
	default:
		return fmt.Errorf("encrypted serializer: unsupported source type %T", dbValue)
	}

	if enc == "" {
		target.SetString("")
		return nil
	}

	plain, err := crypto.Decrypt(enc)
	if err != nil {
		return fmt.Errorf("encrypted serializer: decrypt failed: %w", err)
	}
	target.SetString(plain)
	return nil
}

// Value encrypts the Go field before it is persisted.
func (EncryptedSerializer) Value(_ context.Context, _ *schema.Field, _ reflect.Value, fieldValue any) (any, error) {
	plain, ok := fieldValue.(string)
	if !ok {
		plain = fmt.Sprintf("%v", fieldValue)
	}
	if plain == "" {
		return "", nil
	}
	return crypto.Encrypt(plain)
}
