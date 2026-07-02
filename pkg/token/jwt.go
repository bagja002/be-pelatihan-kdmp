// Package token issues and validates JWT access/refresh tokens using HS256.
package token

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Type distinguishes access tokens from refresh tokens so one cannot be
// substituted for the other.
type Type string

const (
	Access  Type = "access"
	Refresh Type = "refresh"
)

// ErrInvalidToken is returned for any invalid, expired or wrong-type token.
var ErrInvalidToken = errors.New("token: invalid or expired token")

// Claims are the custom JWT claims carried by every token.
type Claims struct {
	UserID uint   `json:"uid"`
	Role   string `json:"role"`
	Type   Type   `json:"typ"`
	jwt.RegisteredClaims
}

// Manager creates and verifies tokens with a fixed secret and TTLs.
type Manager struct {
	secret     []byte
	accessTTL  time.Duration
	refreshTTL time.Duration
	pesertaTTL time.Duration
	issuer     string
}

// NewManager builds a token Manager.
func NewManager(secret string, accessTTL, refreshTTL, pesertaTTL time.Duration, issuer string) *Manager {
	return &Manager{
		secret:     []byte(secret),
		accessTTL:  accessTTL,
		refreshTTL: refreshTTL,
		pesertaTTL: pesertaTTL,
		issuer:     issuer,
	}
}

// GeneratePeserta issues a short-lived access token for a peserta session
// (role "peserta"), used after successful NIK verification.
func (m *Manager) GeneratePeserta(pesertaID uint) (string, error) {
	return m.generate(pesertaID, "peserta", Access, m.pesertaTTL)
}

// AccessTTL exposes the access-token lifetime.
func (m *Manager) AccessTTL() time.Duration { return m.accessTTL }

// GenerateAccess issues a signed access token.
func (m *Manager) GenerateAccess(userID uint, role string) (string, error) {
	return m.generate(userID, role, Access, m.accessTTL)
}

// GenerateRefresh issues a signed refresh token.
func (m *Manager) GenerateRefresh(userID uint, role string) (string, error) {
	return m.generate(userID, role, Refresh, m.refreshTTL)
}

func (m *Manager) generate(userID uint, role string, t Type, ttl time.Duration) (string, error) {
	now := time.Now()
	claims := Claims{
		UserID: userID,
		Role:   role,
		Type:   t,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    m.issuer,
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(m.secret)
}

// Parse validates a token's signature, expiry and type. It rejects any token
// not signed with HS256 to prevent algorithm-confusion attacks.
func (m *Manager) Parse(tokenStr string, expected Type) (*Claims, error) {
	claims := &Claims{}
	parsed, err := jwt.ParseWithClaims(
		tokenStr,
		claims,
		func(t *jwt.Token) (any, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, ErrInvalidToken
			}
			return m.secret, nil
		},
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
	)
	if err != nil || !parsed.Valid || claims.Type != expected {
		return nil, ErrInvalidToken
	}
	return claims, nil
}
