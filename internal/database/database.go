package database

import (
	"log"

	"knmp-backend/internal/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Connect opens a GORM connection to MySQL using the given config.
func Connect(cfg *config.Config) *gorm.DB {
	db, err := gorm.Open(mysql.Open(cfg.DSN()), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}
	log.Println("database connected")
	return db
}
