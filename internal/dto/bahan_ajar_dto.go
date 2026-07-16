package dto

import "knmp-backend/internal/entity"

// BahanAjarKategoriRequest — buat/ubah kategori bahan ajar (body JSON).
type BahanAjarKategoriRequest struct {
	Nama   string `json:"nama" validate:"required,max=255"`
	Urutan *int   `json:"urutan" validate:"omitempty,min=0"`
}

// BahanAjarItemResponse — item untuk daftar publik/admin.
// Path berkas tidak diekspos; hanya penanda ketersediaan.
type BahanAjarItemResponse struct {
	ID     uint   `json:"id"`
	Judul  string `json:"judul"`
	Urutan int    `json:"urutan"`
	HasPdf bool   `json:"hasPdf"`
	HasPpt bool   `json:"hasPpt"`
}

type BahanAjarKategoriResponse struct {
	ID     uint                    `json:"id"`
	Nama   string                  `json:"nama"`
	Urutan int                     `json:"urutan"`
	Items  []BahanAjarItemResponse `json:"items"`
}

func ToBahanAjarItemResponse(b entity.BahanAjar) BahanAjarItemResponse {
	return BahanAjarItemResponse{
		ID:     b.ID,
		Judul:  b.Judul,
		Urutan: b.Urutan,
		HasPdf: b.FilePdf != "",
		HasPpt: b.FilePpt != "",
	}
}

func ToBahanAjarKategoriResponse(k entity.BahanAjarKategori) BahanAjarKategoriResponse {
	items := make([]BahanAjarItemResponse, 0, len(k.Items))
	for _, b := range k.Items {
		items = append(items, ToBahanAjarItemResponse(b))
	}
	return BahanAjarKategoriResponse{ID: k.ID, Nama: k.Nama, Urutan: k.Urutan, Items: items}
}
