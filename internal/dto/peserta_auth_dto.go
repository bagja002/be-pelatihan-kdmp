package dto

// VerifyNIKRequest — peserta memverifikasi identitas via kode satdik + NIK.
type VerifyNIKRequest struct {
	KodeSatdik string `json:"kodeSatdik" validate:"required,max=64"`
	NIK        string `json:"nik" validate:"required,min=16,max=32"`
}
