package handler

import "testing"

func TestSanitizeFilename(t *testing.T) {
	cases := []struct {
		name  string
		judul string
		want  string
	}{
		{"ascii biasa tidak berubah", "Unit Kompetensi 1", "Unit Kompetensi 1"},
		{"kutip newline em-dash disanitasi", "Judul \"Aneh\"\r\nDengan — Em Dash", "Judul -Aneh---Dengan - Em Dash"},
		{"kosong jadi default", "", "bahan-ajar"},
		{"hanya spasi jadi default", "   ", "bahan-ajar"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := sanitizeFilename(tc.judul)
			if got != tc.want {
				t.Errorf("sanitizeFilename(%q) = %q, ingin %q", tc.judul, got, tc.want)
			}
		})
	}
}
