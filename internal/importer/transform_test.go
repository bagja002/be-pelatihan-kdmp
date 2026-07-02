package importer

import "testing"

func TestNormalizeHeader(t *testing.T) {
	cases := map[string]string{
		"  Tanggal_Lahir ": "tanggal_lahir",
		"SATDIK_PELATIHAN": "satdik_pelatihan",
		"provinsi by nik":  "provinsibynik",
	}
	for in, want := range cases {
		if got := NormalizeHeader(in); got != want {
			t.Errorf("NormalizeHeader(%q)=%q want %q", in, got, want)
		}
	}
}

func TestParseTanggalLahir(t *testing.T) {
	cases := map[string]string{
		"18 Februari 2002": "2002-02-18",
		"2 Agustus 1999":   "1999-08-02",
		"1 Januari 1990":   "1990-01-01",
		"31 Desember 2000": "2000-12-31",
		"":                 "",
		"bukan tanggal":    "bukan tanggal",
	}
	for in, want := range cases {
		if got := ParseTanggalLahir(in); got != want {
			t.Errorf("ParseTanggalLahir(%q)=%q want %q", in, got, want)
		}
	}
}

func TestRowToPeserta(t *testing.T) {
	row := map[string]string{
		"nama":            "Budi",
		"nik":             "1234567890123456",
		"tanggal_lahir":   "2 Agustus 1999",
		"jenis_kelamin":   "Laki-laki",
		"provinsi_by_nik": "Sumatera Utara",
		"kotakab_by_nik":  "Kota Medan",
		"provinsi":        "IGNORED",
		"jabatan_knmp":    "Operator",
		"ijazah_users":    "ijazah.pdf",
	}
	p, err := RowToPeserta(row)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.Nama != "Budi" || p.NIK != "1234567890123456" {
		t.Fatalf("basic fields wrong: %+v", p)
	}
	if p.TanggalLahir != "1999-08-02" {
		t.Errorf("tanggal = %q", p.TanggalLahir)
	}
	if p.Provinsi != "Sumatera Utara" || p.Kota != "Kota Medan" {
		t.Errorf("wilayah by_nik salah: prov=%q kota=%q", p.Provinsi, p.Kota)
	}
	if p.Jabatan != "Operator" || p.Ijazah != "ijazah.pdf" {
		t.Errorf("mapped fields wrong: %+v", p)
	}
}

func TestRowToPesertaFallbackWilayah(t *testing.T) {
	row := map[string]string{
		"nik":      "1",
		"provinsi": "Aceh",
		"kota":     "Kota Banda Aceh",
	}
	p, _ := RowToPeserta(row)
	if p.Provinsi != "Aceh" || p.Kota != "Kota Banda Aceh" {
		t.Errorf("fallback wilayah salah: prov=%q kota=%q", p.Provinsi, p.Kota)
	}
}

func TestRowToPesertaNIKKosong(t *testing.T) {
	if _, err := RowToPeserta(map[string]string{"nama": "X"}); err == nil {
		t.Fatal("expected error untuk NIK kosong")
	}
}
