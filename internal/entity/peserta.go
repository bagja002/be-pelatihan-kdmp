package entity

import (
	"time"

	"knmp-backend/internal/database"

	"gorm.io/gorm"
)

// DataPeserta menyimpan data seorang peserta pada suatu Satdik.
type DataPeserta struct {
	ID       uint `gorm:"primaryKey" json:"id"`
	IDSatdik uint `gorm:"index;not null" json:"idSatdik"`

	// Profil (dihitung dalam kelengkapan di FE)
	Nama                string `gorm:"type:varchar(255)" json:"nama"`
	NIK                 string `gorm:"column:nik;type:varchar(32);index" json:"nik"`
	NoTelepon           string `gorm:"type:varchar(64)" json:"noTelepon"`
	Email               string `gorm:"type:varchar(255)" json:"email"`
	Provinsi            string `gorm:"type:varchar(128)" json:"provinsi"`
	Kota                string `gorm:"type:varchar(128)" json:"kota"`
	Kecamatan           string `gorm:"type:varchar(128)" json:"kecamatan"` // nama kecamatan (dulu kode, karenanya sempat varchar(16))
	Alamat              string `gorm:"type:varchar(512)" json:"alamat"`
	TempatLahir         string `gorm:"type:varchar(128)" json:"tempatLahir"`
	TanggalLahir        string `gorm:"type:varchar(32)" json:"tanggalLahir"`
	JenisKelamin        string `gorm:"type:varchar(32)" json:"jenisKelamin"`
	Pekerjaan           string `gorm:"type:varchar(128)" json:"pekerjaan"`
	GolonganDarah       string `gorm:"type:varchar(8)" json:"golonganDarah"`
	StatusMenikah       string `gorm:"type:varchar(32)" json:"statusMenikah"`
	Kewarganegaraan     string `gorm:"type:varchar(32)" json:"kewarganegaraan"`
	IbuKandung          string `gorm:"type:varchar(255)" json:"ibuKandung"`
	NegaraTujuanBekerja string `gorm:"type:varchar(128)" json:"negaraTujuanBekerja"`
	PendidikanTerakhir  string `gorm:"type:varchar(32)" json:"pendidikanTerakhir"`
	Universitas         string `gorm:"type:varchar(255)" json:"universitas"`
	Jurusan             string `gorm:"type:varchar(255)" json:"jurusan"`
	Agama               string `gorm:"type:varchar(32)" json:"agama"`
	Jabatan             string `gorm:"type:varchar(128)" json:"jabatan"`

	// Dokumen (tidak dihitung kelengkapan) — nama berkas / URL
	Foto           string `gorm:"type:varchar(512)" json:"foto"`
	KTP            string `gorm:"column:ktp;type:varchar(512)" json:"ktp"`
	KK             string `gorm:"column:kk;type:varchar(512)" json:"kk"`
	SuratKesehatan string `gorm:"type:varchar(512)" json:"suratKesehatan"`
	Ijazah         string `gorm:"type:varchar(512)" json:"ijazah"`

	// Sistem
	Status   string `gorm:"type:varchar(32)" json:"status,omitempty"`
	Duplikat bool   `gorm:"default:false" json:"duplikat"`

	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func init() { database.RegisterModel(&DataPeserta{}) }
