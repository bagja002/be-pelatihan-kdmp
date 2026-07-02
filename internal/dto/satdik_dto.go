package dto

// CreateSatdikRequest membuat satdik baru.
type CreateSatdikRequest struct {
	Kode      string `json:"kode" validate:"required,max=64"`
	Nama      string `json:"nama" validate:"required,max=255"`
	Lokasi    string `json:"lokasi" validate:"omitempty,max=255"`
	Provinsi  string `json:"provinsi" validate:"omitempty,max=128"`
	PicSatdik string `json:"picSatdik" validate:"omitempty,max=255"`
	NoPic     string `json:"noPic" validate:"omitempty,max=64"`
}

// UpdateSatdikRequest memperbarui satdik (semua field opsional).
type UpdateSatdikRequest struct {
	Nama      string `json:"nama" validate:"omitempty,max=255"`
	Lokasi    string `json:"lokasi" validate:"omitempty,max=255"`
	Provinsi  string `json:"provinsi" validate:"omitempty,max=128"`
	PicSatdik string `json:"picSatdik" validate:"omitempty,max=255"`
	NoPic     string `json:"noPic" validate:"omitempty,max=64"`
}
