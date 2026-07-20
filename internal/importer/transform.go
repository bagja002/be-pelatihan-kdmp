// Package importer mem-parse berkas Excel peserta (sheet GabungData) menjadi
// entity DataPeserta.
package importer

import (
	"strings"

	"knmp-backend/internal/entity"
)

// NIKPlaceholder dipakai untuk baris yang NIK-nya kosong di berkas — datanya
// tetap dimasukkan, hanya NIK-nya ditandai dengan nilai ini.
const NIKPlaceholder = "00000000"

// NormalizeHeader menyamakan header: trim, lowercase, buang spasi.
func NormalizeHeader(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	return strings.ReplaceAll(s, " ", "")
}

var bulanID = map[string]string{
	"januari": "01", "februari": "02", "maret": "03", "april": "04",
	"mei": "05", "juni": "06", "juli": "07", "agustus": "08",
	"september": "09", "oktober": "10", "november": "11", "desember": "12",
}

// ParseTanggalLahir mengubah "18 Februari 2002" menjadi "2002-02-18".
// Input yang tidak cocok dikembalikan apa adanya.
func ParseTanggalLahir(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return ""
	}
	parts := strings.Fields(s)
	if len(parts) != 3 {
		return s
	}
	day := parts[0]
	mm, ok := bulanID[strings.ToLower(parts[1])]
	if !ok {
		return s
	}
	year := parts[2]
	if len(day) == 1 {
		day = "0" + day
	}
	if len(day) != 2 || len(year) != 4 {
		return s
	}
	return year + "-" + mm + "-" + day
}

// get mengambil nilai kolom (kunci sudah ternormalisasi) dan trim; mengembalikan
// nilai bukan-kosong pertama dari daftar kunci.
func get(row map[string]string, keys ...string) string {
	for _, k := range keys {
		if v, ok := row[k]; ok && strings.TrimSpace(v) != "" {
			return strings.TrimSpace(v)
		}
	}
	return ""
}

// RowToPeserta memetakan satu baris menjadi DataPeserta. IDSatdik & Duplikat
// diisi oleh service. Kunci map adalah header ternormalisasi (NormalizeHeader).
func RowToPeserta(row map[string]string) (*entity.DataPeserta, error) {
	nik := get(row, "nik")
	if nik == "" {
		nik = NIKPlaceholder
	}
	p := &entity.DataPeserta{
		Nama:                get(row, "nama"),
		NIK:                 nik,
		NoTelepon:           get(row, "no_telpon", "notelpon"),
		Email:               get(row, "email"),
		Alamat:              get(row, "alamat"),
		Provinsi:            get(row, "provinsi_by_nik", "provinsibynik", "provinsi"),
		Kota:                get(row, "kotakab_by_nik", "kotakabbynik", "kota"),
		Kecamatan:           get(row, "kecamatan"),
		TempatLahir:         get(row, "tempat_lahir", "tempatlahir"),
		TanggalLahir:        ParseTanggalLahir(get(row, "tanggal_lahir", "tanggallahir")),
		JenisKelamin:        get(row, "jenis_kelamin", "jeniskelamin"),
		Pekerjaan:           get(row, "pekerjaan"),
		GolonganDarah:       get(row, "golongan_darah", "golongandarah"),
		StatusMenikah:       get(row, "status_menikah", "statusmenikah"),
		Kewarganegaraan:     get(row, "kewarganegaraan"),
		IbuKandung:          get(row, "ibu_kandung", "ibukandung"),
		NegaraTujuanBekerja: get(row, "negara_tujuan_bekerja", "negaratujuanbekerja"),
		PendidikanTerakhir:  get(row, "pendidikan_terakhir", "pendidikanterakhir"),
		Agama:               get(row, "agama"),
		Universitas:         get(row, "universitas_terakhir", "universitasterakhir"),
		Jurusan:             get(row, "jurusan"),
		Jabatan:             get(row, "jabatan_knmp", "jabatanknmp"),
		Foto:                get(row, "foto"),
		KTP:                 get(row, "ktp"),
		KK:                  get(row, "kk"),
		SuratKesehatan:      get(row, "surat_kesehatan", "suratkesehatan"),
		Ijazah:              get(row, "ijazah_users", "ijazahusers", "ijazah"),
		Status:              get(row, "status"),
	}
	return p, nil
}
