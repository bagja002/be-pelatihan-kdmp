package dto

// RegisterPelatihRequest — field teks dari form registrasi mandiri.
// Berkas (CV & sertifikat) diambil terpisah dari multipart form oleh handler.
type RegisterPelatihRequest struct {
	NamaLengkap  string `json:"namaLengkap" validate:"required,max=255"`
	NIP          string `json:"nip" validate:"required,max=64"`
	Pendidikan   string `json:"pendidikan" validate:"omitempty,max=128"`
	Jurusan      string `json:"jurusan" validate:"omitempty,max=255"`
	Universitas  string `json:"universitas" validate:"omitempty,max=255"`
	UnitKerja    string `json:"unitKerja" validate:"omitempty,max=255"`
	Jabatan      string `json:"jabatan" validate:"omitempty,max=128"`
	Golongan     string `json:"golongan" validate:"omitempty,max=16"`
	Kriteria     string `json:"kriteria" validate:"omitempty,max=16"`
	LokasiTOT    string `json:"lokasiTot" validate:"omitempty,max=255"`
	KelasJabatan string `json:"kelasJabatan" validate:"omitempty,max=64"`
}
