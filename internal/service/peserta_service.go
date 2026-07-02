package service

import (
	"errors"

	"knmp-backend/internal/dto"
	"knmp-backend/internal/entity"
	"knmp-backend/internal/repository"
)

// ErrForbiddenScope: peserta di luar cakupan user.
var ErrForbiddenScope = errors.New("resource outside your scope")

type PesertaService interface {
	List(role string, userID uint) ([]entity.DataPeserta, error)
	Get(role string, userID, id uint) (*entity.DataPeserta, error)
	Create(role string, userID uint, req *dto.CreatePesertaRequest) (*entity.DataPeserta, error)
	Update(role string, userID, id uint, req *dto.UpdatePesertaRequest) (*entity.DataPeserta, error)
	Delete(role string, userID, id uint) error
	// Self (peserta)
	GetSelf(pesertaID uint) (*entity.DataPeserta, error)
	UpdateSelf(pesertaID uint, req *dto.UpdateSelfRequest) (*entity.DataPeserta, error)
}

type pesertaService struct {
	repo  repository.PesertaRepository
	users repository.UserRepository
}

func NewPesertaService(repo repository.PesertaRepository, users repository.UserRepository) PesertaService {
	return &pesertaService{repo: repo, users: users}
}

// scopeSatdik mengembalikan (all, satdikID) untuk user admin berdasarkan DB.
func (s *pesertaService) scopeSatdik(role string, userID uint) (bool, uint, error) {
	if role == "super_admin" {
		return true, 0, nil
	}
	u, err := s.users.FindByID(userID)
	if err != nil {
		return false, 0, err
	}
	all, satID := ResolveScope(u.Type, u.IDSatdik)
	return all, satID, nil
}

func (s *pesertaService) List(role string, userID uint) ([]entity.DataPeserta, error) {
	all, satID, err := s.scopeSatdik(role, userID)
	if err != nil {
		return nil, err
	}
	if all {
		return s.repo.FindAllScoped(nil)
	}
	return s.repo.FindAllScoped(&satID)
}

func (s *pesertaService) authorize(role string, userID, id uint) (*entity.DataPeserta, error) {
	p, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	all, satID, err := s.scopeSatdik(role, userID)
	if err != nil {
		return nil, err
	}
	if !all && p.IDSatdik != satID {
		return nil, ErrForbiddenScope
	}
	return p, nil
}

func (s *pesertaService) Get(role string, userID, id uint) (*entity.DataPeserta, error) {
	return s.authorize(role, userID, id)
}

func (s *pesertaService) Create(role string, userID uint, req *dto.CreatePesertaRequest) (*entity.DataPeserta, error) {
	all, satID, err := s.scopeSatdik(role, userID)
	if err != nil {
		return nil, err
	}
	e := &entity.DataPeserta{}
	applyPeserta(e, req)
	if !all {
		e.IDSatdik = satID // admin dipaksa ke satdiknya
	}
	if e.IDSatdik == 0 {
		return nil, ErrForbiddenScope
	}
	if err := s.repo.Create(e); err != nil {
		return nil, err
	}
	return e, nil
}

func (s *pesertaService) Update(role string, userID, id uint, req *dto.UpdatePesertaRequest) (*entity.DataPeserta, error) {
	p, err := s.authorize(role, userID, id)
	if err != nil {
		return nil, err
	}
	all, satID, err := s.scopeSatdik(role, userID)
	if err != nil {
		return nil, err
	}
	applyPeserta(p, req)
	if !all {
		p.IDSatdik = satID // admin tak bisa memindah peserta keluar satdiknya
	}
	if err := s.repo.Update(p); err != nil {
		return nil, err
	}
	return p, nil
}

func (s *pesertaService) Delete(role string, userID, id uint) error {
	if _, err := s.authorize(role, userID, id); err != nil {
		return err
	}
	return s.repo.Delete(id)
}

func (s *pesertaService) GetSelf(pesertaID uint) (*entity.DataPeserta, error) {
	return s.repo.FindByID(pesertaID)
}

func (s *pesertaService) UpdateSelf(pesertaID uint, req *dto.UpdateSelfRequest) (*entity.DataPeserta, error) {
	p, err := s.repo.FindByID(pesertaID)
	if err != nil {
		return nil, err
	}
	applySelf(p, req) // NIK & idSatdik sengaja tidak disentuh
	if err := s.repo.Update(p); err != nil {
		return nil, err
	}
	return p, nil
}

func applyPeserta(e *entity.DataPeserta, r *dto.CreatePesertaRequest) {
	e.IDSatdik = r.IDSatdik
	e.Nama = r.Nama
	e.NIK = r.NIK
	e.NoTelepon = r.NoTelepon
	e.Email = r.Email
	e.Provinsi = r.Provinsi
	e.Kota = r.Kota
	e.Kecamatan = r.Kecamatan
	e.Alamat = r.Alamat
	e.TempatLahir = r.TempatLahir
	e.TanggalLahir = r.TanggalLahir
	e.JenisKelamin = r.JenisKelamin
	e.Pekerjaan = r.Pekerjaan
	e.GolonganDarah = r.GolonganDarah
	e.StatusMenikah = r.StatusMenikah
	e.Kewarganegaraan = r.Kewarganegaraan
	e.IbuKandung = r.IbuKandung
	e.NegaraTujuanBekerja = r.NegaraTujuanBekerja
	e.PendidikanTerakhir = r.PendidikanTerakhir
	e.Universitas = r.Universitas
	e.Jurusan = r.Jurusan
	e.Agama = r.Agama
	e.Jabatan = r.Jabatan
	e.Foto = r.Foto
	e.KTP = r.KTP
	e.KK = r.KK
	e.SuratKesehatan = r.SuratKesehatan
	e.Ijazah = r.Ijazah
	e.Status = r.Status
}

func applySelf(e *entity.DataPeserta, r *dto.UpdateSelfRequest) {
	e.Nama = r.Nama
	e.NoTelepon = r.NoTelepon
	e.Email = r.Email
	e.Provinsi = r.Provinsi
	e.Kota = r.Kota
	e.Kecamatan = r.Kecamatan
	e.Alamat = r.Alamat
	e.TempatLahir = r.TempatLahir
	e.TanggalLahir = r.TanggalLahir
	e.JenisKelamin = r.JenisKelamin
	e.Pekerjaan = r.Pekerjaan
	e.GolonganDarah = r.GolonganDarah
	e.StatusMenikah = r.StatusMenikah
	e.Kewarganegaraan = r.Kewarganegaraan
	e.IbuKandung = r.IbuKandung
	e.NegaraTujuanBekerja = r.NegaraTujuanBekerja
	e.PendidikanTerakhir = r.PendidikanTerakhir
	e.Universitas = r.Universitas
	e.Jurusan = r.Jurusan
	e.Agama = r.Agama
	e.Jabatan = r.Jabatan
	e.Foto = r.Foto
	e.KTP = r.KTP
	e.KK = r.KK
	e.SuratKesehatan = r.SuratKesehatan
	e.Ijazah = r.Ijazah
}
