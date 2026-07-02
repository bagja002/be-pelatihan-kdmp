package service

import (
	"knmp-backend/internal/dto"
	"knmp-backend/internal/entity"
	"knmp-backend/internal/repository"
)

type SatdikService interface {
	Create(req *dto.CreateSatdikRequest) (*entity.Satdik, error)
	GetAll() ([]entity.Satdik, error)
	GetByID(id uint) (*entity.Satdik, error)
	Update(id uint, req *dto.UpdateSatdikRequest) (*entity.Satdik, error)
	Delete(id uint) error
}

type satdikService struct{ repo repository.SatdikRepository }

func NewSatdikService(repo repository.SatdikRepository) SatdikService {
	return &satdikService{repo: repo}
}

func (s *satdikService) Create(req *dto.CreateSatdikRequest) (*entity.Satdik, error) {
	e := &entity.Satdik{
		Kode: req.Kode, Nama: req.Nama, Lokasi: req.Lokasi,
		Provinsi: req.Provinsi, PicSatdik: req.PicSatdik, NoPic: req.NoPic,
	}
	if err := s.repo.Create(e); err != nil {
		return nil, err
	}
	return e, nil
}

func (s *satdikService) GetAll() ([]entity.Satdik, error)        { return s.repo.FindAll() }
func (s *satdikService) GetByID(id uint) (*entity.Satdik, error) { return s.repo.FindByID(id) }

func (s *satdikService) Update(id uint, req *dto.UpdateSatdikRequest) (*entity.Satdik, error) {
	e, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if req.Nama != "" {
		e.Nama = req.Nama
	}
	e.Lokasi = req.Lokasi
	e.Provinsi = req.Provinsi
	e.PicSatdik = req.PicSatdik
	e.NoPic = req.NoPic
	if err := s.repo.Update(e); err != nil {
		return nil, err
	}
	return e, nil
}

func (s *satdikService) Delete(id uint) error { return s.repo.Delete(id) }
