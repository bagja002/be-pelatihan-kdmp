package service

import (
	"errors"

	"knmp-backend/internal/dto"
	"knmp-backend/internal/entity"
	"knmp-backend/internal/repository"
	"knmp-backend/pkg/hash"
	"knmp-backend/pkg/token"
)

var (
	// ErrEmailTaken is returned when registering an already-used email.
	ErrEmailTaken = errors.New("email already registered")
	// ErrInvalidCredentials is a deliberately generic auth failure so that
	// clients cannot distinguish "unknown email" from "wrong password".
	ErrInvalidCredentials = errors.New("invalid credentials")
)

// AuthService handles registration, login and token refresh.
type AuthService interface {
	Register(req *dto.RegisterRequest) (*entity.User, error)
	Login(req *dto.LoginRequest) (*dto.TokenResponse, error)
	Refresh(refreshToken string) (*dto.TokenResponse, error)
	Profile(userID uint) (*entity.User, error)
}

type authService struct {
	users  repository.UserRepository
	tokens *token.Manager
}

// NewAuthService wires the auth service to its dependencies.
func NewAuthService(users repository.UserRepository, tokens *token.Manager) AuthService {
	return &authService{users: users, tokens: tokens}
}

func (s *authService) Register(req *dto.RegisterRequest) (*entity.User, error) {
	if existing, _ := s.users.FindByEmail(req.Email); existing != nil {
		return nil, ErrEmailTaken
	}

	hashed, err := hash.Password(req.Password)
	if err != nil {
		return nil, err
	}

	u := &entity.User{
		Email:    req.Email,
		Password: hashed,
		Role:     "user",
		Phone:    req.Phone,
	}
	if err := s.users.Create(u); err != nil {
		return nil, err
	}
	return u, nil
}

func (s *authService) Login(req *dto.LoginRequest) (*dto.TokenResponse, error) {
	u, err := s.users.FindByEmail(req.Email)
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
	access, err := s.tokens.GenerateAccess(u.ID, u.Role)
	if err != nil {
		return nil, err
	}
	refresh, err := s.tokens.GenerateRefresh(u.ID, u.Role)
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
