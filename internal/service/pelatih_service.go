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

type CertUpload struct {
	Nama string
	File *multipart.FileHeader
}

type RegisterPelatihInput struct {
	NamaLengkap string
	NIP         string
	Pendidikan  string
	Jurusan     string
	Universitas string
	UnitKerja   string
	Jabatan     string
	Golongan    string
	CV          *multipart.FileHeader // opsional
	Sertifikat  []CertUpload
}

type PelatihService interface {
	Register(in RegisterPelatihInput) (*entity.Pelatih, error)
	List() ([]entity.Pelatih, error)
	Get(id uint) (*entity.Pelatih, error)
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
		NamaLengkap: in.NamaLengkap,
		NIP:         in.NIP,
		Pendidikan:  in.Pendidikan,
		Jurusan:     in.Jurusan,
		Universitas: in.Universitas,
		UnitKerja:   in.UnitKerja,
		Jabatan:     in.Jabatan,
		Golongan:    in.Golongan,
		CV:          cvPath,
		Status:      "baru",
		Sertifikat:  certs,
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
