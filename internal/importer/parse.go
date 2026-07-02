package importer

import (
	"fmt"
	"io"

	"knmp-backend/internal/entity"

	"github.com/xuri/excelize/v2"
)

// Row adalah hasil parse satu baris data.
type Row struct {
	Peserta    *entity.DataPeserta // nil bila SkipReason != ""
	SatdikNama string
	LineNo     int    // nomor baris di sheet (1-based, termasuk header)
	SkipReason string // "" bila valid
}

// findSheet mengembalikan nama sheet yang ternormalisasi == "gabungdata".
func findSheet(f *excelize.File) (string, error) {
	for _, name := range f.GetSheetList() {
		if NormalizeHeader(name) == "gabungdata" {
			return name, nil
		}
	}
	return "", fmt.Errorf("sheet GabungData tidak ditemukan")
}

// Parse membaca berkas xlsx dan mengembalikan baris data (termasuk yang dilewati).
func Parse(r io.Reader) ([]Row, error) {
	f, err := excelize.OpenReader(r)
	if err != nil {
		return nil, fmt.Errorf("gagal membuka xlsx: %w", err)
	}
	defer f.Close()

	sheet, err := findSheet(f)
	if err != nil {
		return nil, err
	}
	rows, err := f.GetRows(sheet)
	if err != nil {
		return nil, fmt.Errorf("gagal membaca baris: %w", err)
	}
	if len(rows) < 2 {
		return nil, fmt.Errorf("sheet %q tidak berisi data", sheet)
	}

	// Header baris pertama → kunci ternormalisasi per indeks kolom.
	headers := make([]string, len(rows[0]))
	for i, h := range rows[0] {
		headers[i] = NormalizeHeader(h)
	}

	out := make([]Row, 0, len(rows)-1)
	for i := 1; i < len(rows); i++ {
		raw := rows[i]
		empty := true
		for _, c := range raw {
			if c != "" {
				empty = false
				break
			}
		}
		if empty {
			continue
		}
		m := make(map[string]string, len(headers))
		for j, key := range headers {
			if key == "" {
				continue
			}
			if j < len(raw) {
				m[key] = raw[j]
			}
		}
		lineNo := i + 1
		satdik := m["satdik_pelatihan"]
		p, err := RowToPeserta(m)
		if err != nil {
			out = append(out, Row{LineNo: lineNo, SatdikNama: satdik, SkipReason: err.Error()})
			continue
		}
		out = append(out, Row{Peserta: p, SatdikNama: satdik, LineNo: lineNo})
	}
	return out, nil
}
