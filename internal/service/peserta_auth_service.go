package service

import (
	"errors"

	"knmp-backend/internal/dto"
	"knmp-backend/internal/repository"
	"knmp-backend/pkg/token"
)

// ErrPesertaVerification generik (tak membocorkan apakah satdik/NIK yang salah).
var ErrPesertaVerification = errors.New("verifikasi gagal")

type PesertaAuthService interface {
	Verify(req *dto.VerifyNIKRequest) (*dto.TokenResponse, error)
}

type pesertaAuthService struct {
	satdik  repository.SatdikRepository
	peserta repository.PesertaRepository
	tokens  *token.Manager
}

func NewPesertaAuthService(satdik repository.SatdikRepository, peserta repository.PesertaRepository, tokens *token.Manager) PesertaAuthService {
	return &pesertaAuthService{satdik: satdik, peserta: peserta, tokens: tokens}
}

func (s *pesertaAuthService) Verify(req *dto.VerifyNIKRequest) (*dto.TokenResponse, error) {
	sat, err := s.satdik.FindByKode(req.KodeSatdik)
	if err != nil {
		return nil, ErrPesertaVerification
	}
	p, err := s.peserta.FindByNIKAndSatdik(req.NIK, sat.ID)
	if err != nil {
		return nil, ErrPesertaVerification
	}
	access, err := s.tokens.GeneratePeserta(p.ID)
	if err != nil {
		return nil, err
	}
	return &dto.TokenResponse{
		AccessToken: access,
		TokenType:   "Bearer",
		ExpiresIn:   int64(s.tokens.AccessTTL().Seconds()),
	}, nil
}
