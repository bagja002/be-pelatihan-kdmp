package database

import "gorm.io/gorm"

// models is the auto-migration registry. Each entity registers itself
// here in its init() function, so adding a new entity requires no changes
// to the migration code — only registering its route.
var models []any

// RegisterModel adds a model to the auto-migration registry.
// Call this from an entity's init() function.
func RegisterModel(m any) {
	models = append(models, m)
}

// Migrate runs GORM auto-migration for every registered model.
func Migrate(db *gorm.DB) error {
	if len(models) == 0 {
		return nil
	}
	return db.AutoMigrate(models...)
}
