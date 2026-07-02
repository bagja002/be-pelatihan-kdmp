package service

import (
	"errors"

	"knmp-backend/internal/dto"
	"knmp-backend/internal/entity"
	"knmp-backend/internal/repository"
	"knmp-backend/pkg/hash"
	"knmp-backend/pkg/token"
)

// ErrInvalidCredentials adalah kegagalan auth generik (tidak membocorkan
// apakah username salah atau password salah).
var ErrInvalidCredentials = errors.New("invalid credentials")

type AuthService interface {
	Login(req *dto.LoginRequest) (*dto.TokenResponse, error)
	Refresh(refreshToken string) (*dto.TokenResponse, error)
	Profile(userID uint) (*entity.User, error)
}

type authService struct {
	users  repository.UserRepository
	tokens *token.Manager
}

func NewAuthService(users repository.UserRepository, tokens *token.Manager) AuthService {
	return &authService{users: users, tokens: tokens}
}

func (s *authService) Login(req *dto.LoginRequest) (*dto.TokenResponse, error) {
	u, err := s.users.FindByUsername(req.Username)
	if err != nil {
		return nil, ErrInvalidCredentials
	}
	ok, err := hash.Verify(req.Password, u.Password)
	if err != nil || !ok {
		return nil, ErrInvalidCredentials
	}
	return s.issueTokens(u)
}

func (s *authService) Refresh(refreshToken string) (*dto.TokenResponse, error) {
	claims, err := s.tokens.Parse(refreshToken, token.Refresh)
	if err != nil {
		return nil, ErrInvalidCredentials
	}
	u, err := s.users.FindByID(claims.UserID)
	if err != nil {
		return nil, ErrInvalidCredentials
	}
	return s.issueTokens(u)
}

func (s *authService) Profile(userID uint) (*entity.User, error) {
	return s.users.FindByID(userID)
}

func (s *authService) issueTokens(u *entity.User) (*dto.TokenResponse, error) {
	access, err := s.tokens.GenerateAccess(u.ID, u.Type)
	if err != nil {
		return nil, err
	}
	refresh, err := s.tokens.GenerateRefresh(u.ID, u.Type)
	if err != nil {
		return nil, err
	}
	return &dto.TokenResponse{
		AccessToken:  access,
		RefreshToken: refresh,
		TokenType:    "Bearer",
		ExpiresIn:    int64(s.tokens.AccessTTL().Seconds()),
	}, nil
}
