package dto

// RegisterPelatihRequest — field teks dari form registrasi mandiri.
// Berkas (CV & sertifikat) diambil terpisah dari multipart form oleh handler.
type RegisterPelatihRequest struct {
	NamaLengkap  string `json:"namaLengkap" validate:"required,max=255"`
	NIP          string `json:"nip" validate:"required,max=64"`
	NoTelepon    string `json:"noTelepon" validate:"required,max=64"`
	Pendidikan   string `json:"pendidikan" validate:"omitempty,max=128"`
	Jurusan      string `json:"jurusan" validate:"omitempty,max=255"`
	Universitas  string `json:"universitas" validate:"omitempty,max=255"`
	UnitKerja    string `json:"unitKerja" validate:"omitempty,max=255"`
	Jabatan      string `json:"jabatan" validate:"omitempty,max=128"`
	Golongan     string `json:"golongan" validate:"omitempty,max=16"`
	Kriteria     string `json:"kriteria" validate:"omitempty,max=16"`
	LokasiTOT    string `json:"lokasiTot" validate:"required,oneof=Bogor Bandung Surabaya Makassar"`
	KelasJabatan string `json:"kelasJabatan" validate:"omitempty,max=64"`
}

// AdminUpdatePelatihRequest — field teks yang boleh diubah admin lewat dashboard.
// NIP adalah kunci dan tidak ikut diubah. Berkas (CV & sertifikat) dikelola
// terpisah lewat alur edit mandiri, bukan endpoint ini.
type AdminUpdatePelatihRequest struct {
	NamaLengkap  string `json:"namaLengkap" validate:"required,max=255"`
	NoTelepon    string `json:"noTelepon" validate:"required,max=64"`
	Pendidikan   string `json:"pendidikan" validate:"omitempty,max=128"`
	Jurusan      string `json:"jurusan" validate:"omitempty,max=255"`
	Universitas  string `json:"universitas" validate:"omitempty,max=255"`
	UnitKerja    string `json:"unitKerja" validate:"omitempty,max=255"`
	Jabatan      string `json:"jabatan" validate:"omitempty,max=128"`
	Golongan     string `json:"golongan" validate:"omitempty,max=16"`
	Kriteria     string `json:"kriteria" validate:"omitempty,max=16"`
	LokasiTOT    string `json:"lokasiTot" validate:"required,oneof=Bogor Bandung Surabaya Makassar"`
	KelasJabatan string `json:"kelasJabatan" validate:"omitempty,max=64"`
}
