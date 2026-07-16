// Command seed prepares the data-knmp database for development: it cleans up
// legacy template columns/tables, runs migration, and ensures a super admin
// account exists. Real satdik & peserta data is loaded via the Excel importer.
package main

import (
	"log"

	"knmp-backend/internal/config"
	"knmp-backend/internal/database"
	"knmp-backend/internal/entity"
	"knmp-backend/pkg/crypto"
	"knmp-backend/pkg/hash"

	"gorm.io/gorm"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}
	if err := crypto.SetKey(cfg.EncryptionKey); err != nil {
		log.Fatalf("crypto: %v", err)
	}
	db := database.Connect(cfg)

	cleanupLegacy(db)

	if err := database.Migrate(db); err != nil {
		log.Fatalf("migrate: %v", err)
	}

	seedSuperAdmin(db)
	seedBahanAjarKategori(db)
	log.Println("seed selesai")
}

// cleanupLegacy removes columns/tables left over from the starter template so
// the reshaped User schema (username/type/id_satdik) migrates cleanly.
func cleanupLegacy(db *gorm.DB) {
	m := db.Migrator()
	if m.HasTable("products") {
		_ = m.DropTable("products")
		log.Println("dropped legacy table: products")
	}
	for _, col := range []string{"email", "phone", "role"} {
		if m.HasColumn(&entity.User{}, col) {
			if err := m.DropColumn(&entity.User{}, col); err == nil {
				log.Printf("dropped legacy users.%s", col)
			}
		}
	}
}

func seedSuperAdmin(db *gorm.DB) {
	var count int64
	db.Model(&entity.User{}).Where("username = ?", "superadmin").Count(&count)
	if count > 0 {
		log.Println("super admin sudah ada, dilewati")
		return
	}
	pw, err := hash.Password("password")
	if err != nil {
		log.Fatalf("hash: %v", err)
	}
	u := &entity.User{Nama: "Super Administrator", Username: "superadmin", Password: pw, Type: "super_admin"}
	if err := db.Create(u).Error; err != nil {
		log.Fatalf("create super admin: %v", err)
	}
	log.Println("super admin dibuat: superadmin / password")
}

// seedBahanAjarKategori membuat 5 kategori bahan ajar awal (idempoten per nama).
// Item (unit kompetensi/modul) tidak di-seed — diisi admin lewat dashboard.
func seedBahanAjarKategori(db *gorm.DB) {
	names := []string{
		"Kompetensi Umum",
		"Kepala Produksi",
		"Manager Operasional",
		"Penjamin Mutu",
		"Administrasi Keuangan",
	}
	for i, nama := range names {
		var count int64
		db.Model(&entity.BahanAjarKategori{}).Where("nama = ?", nama).Count(&count)
		if count > 0 {
			continue
		}
		k := &entity.BahanAjarKategori{Nama: nama, Urutan: i + 1}
		if err := db.Create(k).Error; err != nil {
			log.Fatalf("create kategori bahan ajar %q: %v", nama, err)
		}
		log.Printf("kategori bahan ajar dibuat: %s", nama)
	}
}
