package dto

// CreatePesertaRequest — admin membuat peserta (semua field profil+dokumen).
type CreatePesertaRequest struct {
	IDSatdik            uint   `json:"idSatdik" validate:"required"`
	Nama                string `json:"nama" validate:"omitempty,max=255"`
	NIK                 string `json:"nik" validate:"required,max=32"`
	NoTelepon           string `json:"noTelepon" validate:"omitempty,max=64"`
	Email               string `json:"email" validate:"omitempty,max=255"`
	Provinsi            string `json:"provinsi" validate:"omitempty,max=128"`
	Kota                string `json:"kota" validate:"omitempty,max=128"`
	Kecamatan           string `json:"kecamatan" validate:"omitempty,max=128"`
	Alamat              string `json:"alamat" validate:"omitempty,max=512"`
	TempatLahir         string `json:"tempatLahir" validate:"omitempty,max=128"`
	TanggalLahir        string `json:"tanggalLahir" validate:"omitempty,max=32"`
	JenisKelamin        string `json:"jenisKelamin" validate:"omitempty,max=32"`
	Pekerjaan           string `json:"pekerjaan" validate:"omitempty,max=128"`
	GolonganDarah       string `json:"golonganDarah" validate:"omitempty,max=8"`
	StatusMenikah       string `json:"statusMenikah" validate:"omitempty,max=32"`
	Kewarganegaraan     string `json:"kewarganegaraan" validate:"omitempty,max=32"`
	IbuKandung          string `json:"ibuKandung" validate:"omitempty,max=255"`
	NegaraTujuanBekerja string `json:"negaraTujuanBekerja" validate:"omitempty,max=128"`
	PendidikanTerakhir  string `json:"pendidikanTerakhir" validate:"omitempty,max=32"`
	Universitas         string `json:"universitas" validate:"omitempty,max=255"`
	Jurusan             string `json:"jurusan" validate:"omitempty,max=255"`
	Agama               string `json:"agama" validate:"omitempty,max=32"`
	Jabatan             string `json:"jabatan" validate:"omitempty,max=128"`
	Foto                string `json:"foto" validate:"omitempty,max=512"`
	KTP                 string `json:"ktp" validate:"omitempty,max=512"`
	KK                  string `json:"kk" validate:"omitempty,max=512"`
	SuratKesehatan      string `json:"suratKesehatan" validate:"omitempty,max=512"`
	Ijazah              string `json:"ijazah" validate:"omitempty,max=512"`
	Status              string `json:"status" validate:"omitempty,max=32"`
}

// UpdatePesertaRequest — admin update; semua field boleh (termasuk NIK & idSatdik).
type UpdatePesertaRequest = CreatePesertaRequest

// UpdateSelfRequest — peserta update dirinya; TANPA nik & idSatdik (read-only).
type UpdateSelfRequest struct {
	Nama                string `json:"nama" validate:"omitempty,max=255"`
	NoTelepon           string `json:"noTelepon" validate:"omitempty,max=64"`
	Email               string `json:"email" validate:"omitempty,max=255"`
	Provinsi            string `json:"provinsi" validate:"omitempty,max=128"`
	Kota                string `json:"kota" validate:"omitempty,max=128"`
	Kecamatan           string `json:"kecamatan" validate:"omitempty,max=128"`
	Alamat              string `json:"alamat" validate:"omitempty,max=512"`
	TempatLahir         string `json:"tempatLahir" validate:"omitempty,max=128"`
	TanggalLahir        string `json:"tanggalLahir" validate:"omitempty,max=32"`
	JenisKelamin        string `json:"jenisKelamin" validate:"omitempty,max=32"`
	Pekerjaan           string `json:"pekerjaan" validate:"omitempty,max=128"`
	GolonganDarah       string `json:"golonganDarah" validate:"omitempty,max=8"`
	StatusMenikah       string `json:"statusMenikah" validate:"omitempty,max=32"`
	Kewarganegaraan     string `json:"kewarganegaraan" validate:"omitempty,max=32"`
	IbuKandung          string `json:"ibuKandung" validate:"omitempty,max=255"`
	NegaraTujuanBekerja string `json:"negaraTujuanBekerja" validate:"omitempty,max=128"`
	PendidikanTerakhir  string `json:"pendidikanTerakhir" validate:"omitempty,max=32"`
	Universitas         string `json:"universitas" validate:"omitempty,max=255"`
	Jurusan             string `json:"jurusan" validate:"omitempty,max=255"`
	Agama               string `json:"agama" validate:"omitempty,max=32"`
	Jabatan             string `json:"jabatan" validate:"omitempty,max=128"`
	Foto                string `json:"foto" validate:"omitempty,max=512"`
	KTP                 string `json:"ktp" validate:"omitempty,max=512"`
	KK                  string `json:"kk" validate:"omitempty,max=512"`
	SuratKesehatan      string `json:"suratKesehatan" validate:"omitempty,max=512"`
	Ijazah              string `json:"ijazah" validate:"omitempty,max=512"`
}
