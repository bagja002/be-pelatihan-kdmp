package service

import (
	"errors"
	"mime/multipart"

	"knmp-backend/internal/entity"
	"knmp-backend/internal/repository"
	"knmp-backend/internal/storage"
)

// ErrNIPExists: NIP sudah terdaftar (tidak boleh ada data ganda).
var ErrNIPExists = errors.New("nip already exists")

// ErrNoSertifikat: minimal 1 sertifikat wajib ada (aturan bisnis, ditegakkan di server).
var ErrNoSertifikat = errors.New("minimal 1 sertifikat")

type CertUpload struct {
	Nama string
	File *multipart.FileHeader
}

type RegisterPelatihInput struct {
	NamaLengkap  string
	NIP          string
	Pendidikan   string
	Jurusan      string
	Universitas  string
	UnitKerja    string
	Jabatan      string
	Golongan     string
	Kriteria     string
	LokasiTOT    string
	KelasJabatan string
	CV           *multipart.FileHeader // opsional
	Sertifikat   []CertUpload
}

// UpdateSelfInput — data swakelola edit mandiri oleh pelatih (kunci: NIP).
type UpdateSelfInput struct {
	NIP          string // kunci pencarian; tidak diubah
	NamaLengkap  string
	Pendidikan   string
	Jurusan      string
	Universitas  string
	UnitKerja    string
	Jabatan      string
	Golongan     string
	Kriteria     string
	LokasiTOT    string
	KelasJabatan string
	CV           *multipart.FileHeader // opsional; ganti CV lama bila diisi
	KeepSertifikatIDs []uint            // ID sertifikat lama yang dipertahankan
	NewSertifikat     []CertUpload      // sertifikat baru yang ditambahkan
}

type PelatihService interface {
	Register(in RegisterPelatihInput) (*entity.Pelatih, error)
	List() ([]entity.Pelatih, error)
	Get(id uint) (*entity.Pelatih, error)
	FindByNIP(nip string) (*entity.Pelatih, error)
	UpdateSelf(in UpdateSelfInput) (*entity.Pelatih, error)
	Delete(id uint) error
	Sertifikat(id uint) (*entity.SertifikatKeahlian, error)
}

type pelatihService struct {
	repo  repository.PelatihRepository
	store *storage.Store
}

func NewPelatihService(repo repository.PelatihRepository, store *storage.Store) PelatihService {
	return &pelatihService{repo: repo, store: store}
}

func (s *pelatihService) Register(in RegisterPelatihInput) (*entity.Pelatih, error) {
	// 1. Cek NIP unik lebih dulu — sebelum menyimpan berkas apa pun.
	exists, err := s.repo.ExistsNIP(in.NIP)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrNIPExists
	}

	// 2. Simpan berkas; lacak untuk cleanup bila langkah berikutnya gagal.
	var saved []string
	cleanup := func() {
		for _, rel := range saved {
			_ = s.store.Remove(rel)
		}
	}

	cvPath := ""
	if in.CV != nil {
		p, err := s.store.Save(in.CV, "pelatih")
		if err != nil {
			cleanup()
			return nil, err
		}
		saved = append(saved, p)
		cvPath = p
	}

	certs := make([]entity.SertifikatKeahlian, 0, len(in.Sertifikat))
	for _, cu := range in.Sertifikat {
		berkas := ""
		if cu.File != nil {
			p, err := s.store.Save(cu.File, "pelatih")
			if err != nil {
				cleanup()
				return nil, err
			}
			saved = append(saved, p)
			berkas = p
		}
		certs = append(certs, entity.SertifikatKeahlian{NamaSertifikat: cu.Nama, Berkas: berkas})
	}

	p := &entity.Pelatih{
		NamaLengkap:  in.NamaLengkap,
		NIP:          in.NIP,
		Pendidikan:   in.Pendidikan,
		Jurusan:      in.Jurusan,
		Universitas:  in.Universitas,
		UnitKerja:    in.UnitKerja,
		Jabatan:      in.Jabatan,
		Golongan:     in.Golongan,
		Kriteria:     in.Kriteria,
		LokasiTOT:    in.LokasiTOT,
		KelasJabatan: in.KelasJabatan,
		CV:           cvPath,
		Status:       "baru",
		Sertifikat:   certs,
	}

	// 3. Simpan ke DB. Bila gagal (mis. race pada unique NIP), bersihkan berkas.
	if err := s.repo.Create(p); err != nil {
		cleanup()
		if repository.IsDuplicateNIP(err) {
			return nil, ErrNIPExists
		}
		return nil, err
	}
	return p, nil
}

func (s *pelatihService) List() ([]entity.Pelatih, error) { return s.repo.FindAll() }

func (s *pelatihService) Get(id uint) (*entity.Pelatih, error) { return s.repo.FindByID(id) }

func (s *pelatihService) FindByNIP(nip string) (*entity.Pelatih, error) {
	return s.repo.FindByNIP(nip)
}

// UpdateSelf memperbarui data pelatih berdasarkan NIP: field teks, CV (opsional),
// dan sertifikat (pertahankan sebagian lama + tambah baru). Berlaku langsung,
// status ditandai "diperbarui".
func (s *pelatihService) UpdateSelf(in UpdateSelfInput) (*entity.Pelatih, error) {
	// 1. Temukan data lama (dan sertifikatnya) berdasarkan NIP.
	p, err := s.repo.FindByNIP(in.NIP)
	if err != nil {
		return nil, err
	}

	// 2. Tentukan sertifikat lama mana yang dipertahankan vs dihapus.
	keep := make(map[uint]bool, len(in.KeepSertifikatIDs))
	for _, id := range in.KeepSertifikatIDs {
		keep[id] = true
	}
	var deleteIDs []uint
	var deletePaths []string
	for _, c := range p.Sertifikat {
		if !keep[c.ID] {
			deleteIDs = append(deleteIDs, c.ID)
			if c.Berkas != "" {
				deletePaths = append(deletePaths, c.Berkas)
			}
		}
	}
	keptCount := len(p.Sertifikat) - len(deleteIDs)

	// 3. Aturan bisnis: minimal 1 sertifikat harus tersisa (lama + baru).
	newValid := 0
	for _, cu := range in.NewSertifikat {
		if cu.Nama != "" && cu.File != nil {
			newValid++
		}
	}
	if keptCount+newValid < 1 {
		return nil, ErrNoSertifikat
	}

	// 4. Simpan berkas baru; lacak untuk cleanup bila DB gagal.
	var saved []string
	cleanup := func() {
		for _, rel := range saved {
			_ = s.store.Remove(rel)
		}
	}

	newCVPath := ""
	if in.CV != nil {
		q, err := s.store.Save(in.CV, "pelatih")
		if err != nil {
			cleanup()
			return nil, err
		}
		saved = append(saved, q)
		newCVPath = q
	}

	newSertifikat := make([]entity.SertifikatKeahlian, 0, newValid)
	for _, cu := range in.NewSertifikat {
		if cu.Nama == "" || cu.File == nil {
			continue
		}
		q, err := s.store.Save(cu.File, "pelatih")
		if err != nil {
			cleanup()
			return nil, err
		}
		saved = append(saved, q)
		newSertifikat = append(newSertifikat, entity.SertifikatKeahlian{NamaSertifikat: cu.Nama, Berkas: q})
	}

	// 5. Terapkan field teks (NIP tetap). Ganti CV hanya bila ada berkas baru.
	oldCV := p.CV
	p.NamaLengkap = in.NamaLengkap
	p.Pendidikan = in.Pendidikan
	p.Jurusan = in.Jurusan
	p.Universitas = in.Universitas
	p.UnitKerja = in.UnitKerja
	p.Jabatan = in.Jabatan
	p.Golongan = in.Golongan
	p.Kriteria = in.Kriteria
	p.LokasiTOT = in.LokasiTOT
	p.KelasJabatan = in.KelasJabatan
	p.Status = "diperbarui"
	if newCVPath != "" {
		p.CV = newCVPath
	}

	// 6. Simpan ke DB dalam satu transaksi. Gagal → bersihkan berkas baru.
	if err := s.repo.UpdateSelf(p, deleteIDs, newSertifikat); err != nil {
		cleanup()
		return nil, err
	}

	// 7. Sukses → bersihkan berkas lama yang tak lagi dipakai (best-effort).
	if newCVPath != "" && oldCV != "" {
		_ = s.store.Remove(oldCV)
	}
	for _, rel := range deletePaths {
		_ = s.store.Remove(rel)
	}

	// 8. Kembalikan data terbaru (termasuk sertifikat gabungan).
	return s.repo.FindByNIP(in.NIP)
}

func (s *pelatihService) Delete(id uint) error {
	p, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	if err := s.repo.Delete(id); err != nil {
		return err
	}
	// Best-effort cleanup berkas di disk (abaikan error bila berkas sudah tiada).
	_ = s.store.Remove(p.CV)
	for _, c := range p.Sertifikat {
		_ = s.store.Remove(c.Berkas)
	}
	return nil
}

func (s *pelatihService) Sertifikat(id uint) (*entity.SertifikatKeahlian, error) {
	return s.repo.FindSertifikat(id)
}
