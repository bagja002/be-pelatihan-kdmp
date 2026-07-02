package service

import (
	"regexp"
	"strconv"
	"strings"

	"knmp-backend/internal/entity"
	"knmp-backend/internal/importer"
	"knmp-backend/internal/repository"
)

type RowError struct {
	LineNo int    `json:"lineNo"`
	Alasan string `json:"alasan"`
}

type ImportSummary struct {
	Total            int        `json:"total"`
	Inserted         int        `json:"inserted"`
	Duplikat         int        `json:"duplikat"`
	DilewatiTanpaNIK int        `json:"dilewatiTanpaNik"`
	SatdikDibuat     []string   `json:"satdikDibuat"`
	Errors           []RowError `json:"errors"`
}

type PesertaImportService interface {
	Import(rows []importer.Row) (*ImportSummary, error)
}

type pesertaImportService struct {
	peserta repository.PesertaRepository
	satdik  repository.SatdikRepository
}

func NewPesertaImportService(peserta repository.PesertaRepository, satdik repository.SatdikRepository) PesertaImportService {
	return &pesertaImportService{peserta: peserta, satdik: satdik}
}

var nonAlnum = regexp.MustCompile(`[^A-Z0-9]+`)

// SlugSatdik: "Menart 1 Marinir" -> "SATDIK-MENART-1-MARINIR".
func SlugSatdik(nama string) string {
	s := strings.ToUpper(strings.TrimSpace(nama))
	s = nonAlnum.ReplaceAllString(s, "-")
	s = strings.Trim(s, "-")
	if s == "" {
		s = "TANPA-NAMA"
	}
	if !strings.HasPrefix(s, "SATDIK-") {
		s = "SATDIK-" + s
	}
	return s
}

func (s *pesertaImportService) resolveSatdik(nama string, cache map[string]uint, created *[]string) uint {
	key := strings.ToLower(strings.TrimSpace(nama))
	if key == "" {
		return 0
	}
	if id, ok := cache[key]; ok {
		return id
	}
	if existing, err := s.satdik.FindByNama(strings.TrimSpace(nama)); err == nil {
		cache[key] = existing.ID
		return existing.ID
	}
	sat := &entity.Satdik{Nama: strings.TrimSpace(nama), Kode: s.uniqueKode(nama)}
	if err := s.satdik.Create(sat); err != nil {
		return 0
	}
	cache[key] = sat.ID
	*created = append(*created, sat.Nama)
	return sat.ID
}

func (s *pesertaImportService) uniqueKode(nama string) string {
	base := SlugSatdik(nama)
	kode := base
	for i := 2; ; i++ {
		if _, err := s.satdik.FindByKode(kode); err != nil {
			return kode // tidak ditemukan → bebas dipakai
		}
		kode = base + "-" + strconv.Itoa(i)
	}
}

func (s *pesertaImportService) Import(rows []importer.Row) (*ImportSummary, error) {
	sum := &ImportSummary{Total: len(rows), SatdikDibuat: []string{}, Errors: []RowError{}}

	existing, err := s.peserta.AllNIK()
	if err != nil {
		return nil, err
	}
	satCache := map[string]uint{}
	seen := map[string]bool{}
	batch := make([]entity.DataPeserta, 0, len(rows))

	for _, r := range rows {
		if r.SkipReason != "" {
			sum.DilewatiTanpaNIK++
			continue
		}
		p := r.Peserta
		idSat := s.resolveSatdik(r.SatdikNama, satCache, &sum.SatdikDibuat)
		if idSat == 0 {
			sum.Errors = append(sum.Errors, RowError{LineNo: r.LineNo, Alasan: "satdik tidak dapat diresolusi"})
			continue
		}
		p.IDSatdik = idSat
		if existing[p.NIK] || seen[p.NIK] {
			p.Duplikat = true
			sum.Duplikat++
		}
		seen[p.NIK] = true
		batch = append(batch, *p)
	}

	if err := s.peserta.CreateBatch(batch); err != nil {
		return nil, err
	}
	sum.Inserted = len(batch)
	return sum, nil
}
